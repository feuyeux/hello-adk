package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"strings"

	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

// LiteLlm implements the model.LLM interface for Ollama models using /api/chat
type LiteLlm struct {
	modelName string
	baseURL   string
}

// NewLiteLlm creates a new LiteLlm instance
func NewLiteLlm(modelName string) (model.LLM, error) {
	return &LiteLlm{
		modelName: modelName,
		baseURL:   "http://localhost:11434", // Default Ollama URL
	}, nil
}

// Name returns the model name
func (l *LiteLlm) Name() string {
	return l.modelName
}

// GenerateContent generates content using the Ollama /api/chat endpoint
func (l *LiteLlm) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		// Convert ADK request to Ollama chat request format
		chatReq, err := l.convertToChatRequest(req)
		if err != nil {
			yield(nil, err)
			return
		}

		// Send request to Ollama /api/chat
		chatResp, err := l.sendChatRequest(ctx, chatReq)
		if err != nil {
			yield(nil, err)
			return
		}

		// Convert Ollama response to ADK response
		adkResp, err := l.convertToADKResponse(chatResp)
		if err != nil {
			yield(nil, err)
			return
		}

		// Yield the single response
		yield(adkResp, nil)
	}
}

// Ollama Chat API Request Structures

// ChatRequest represents the request format for Ollama /api/chat endpoint
type ChatRequest struct {
	Model    string   `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
	Stream   bool     `json:"stream"`
	Temperature float64 `json:"temperature,omitempty"`
}

// Message represents a chat message in Ollama /api/chat format
type Message struct {
	Role      string      `json:"role"`
	Content   string      `json:"content,omitempty"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool call in the message
type ToolCall struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Function ToolCallDetail `json:"function"`
}

// ToolCallDetail represents the function call details
type ToolCallDetail struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"` // JSON object from Ollama
}

// Tool represents a tool definition in Ollama format
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction represents the function definition
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolParameter represents a parameter in the function schema
type ToolParameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ToolParameterDefinition represents the parameters object
type ToolParameterDefinition struct {
	Type       string                    `json:"type"`
	Properties map[string]ToolParameter  `json:"properties"`
	Required   []string                 `json:"required"`
}

// Ollama Chat API Response Structures

// ChatResponse represents the response format from Ollama /api/chat endpoint
type ChatResponse struct {
	Model     string     `json:"model"`
	CreatedAt string     `json:"created_at"`
	Message   ChatMessage `json:"message"`
	Done      bool       `json:"done"`
}

// ChatMessage represents the message in the response
type ChatMessage struct {
	Role      string      `json:"role"`
	Content   string      `json:"content,omitempty"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
}

// convertToChatRequest converts ADK LLMRequest to Ollama /api/chat request format
func (l *LiteLlm) convertToChatRequest(req *model.LLMRequest) (*ChatRequest, error) {
	messages, err := l.convertMessages(req.Contents)
	if err != nil {
		return nil, err
	}

	tools := l.convertTools(req.Tools)

	return &ChatRequest{
		Model:       l.modelName,
		Messages:    messages,
		Tools:       tools,
		Stream:      false,
		Temperature: 0.7,
	}, nil
}

// convertMessages converts ADK contents to Ollama messages
func (l *LiteLlm) convertMessages(contents []*genai.Content) ([]Message, error) {
	var messages []Message

	for _, content := range contents {
		if content.Parts == nil || len(content.Parts) == 0 {
			continue
		}

		msg := Message{
			Role: content.Role,
		}

		var hasContent bool
		var toolCalls []ToolCall

		for _, part := range content.Parts {
			// Handle text parts
			if part.Text != "" {
				msg.Content = part.Text
				hasContent = true
			}

			// Handle function call parts
			if part.FunctionCall != nil {
				var argsMap map[string]interface{}
				argsJSON, err := json.Marshal(part.FunctionCall.Args)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal function call args: %v", err)
				}
				if err := json.Unmarshal(argsJSON, &argsMap); err != nil {
					return nil, fmt.Errorf("failed to unmarshal function call args to map: %v", err)
				}

				toolCall := ToolCall{
					ID:   fmt.Sprintf("call_%d", len(toolCalls)),
					Type: "function",
					Function: ToolCallDetail{
						Name:      part.FunctionCall.Name,
						Arguments: argsMap,
					},
				}
				toolCalls = append(toolCalls, toolCall)
				hasContent = true
			}

			// Handle function response parts (tool results)
			if part.FunctionResponse != nil && part.FunctionResponse.Name != "" {
				// Convert function response to a tool message
				responseContent := "{}"
				if part.FunctionResponse.Response != nil {
					responseContent = fmt.Sprintf("%v", part.FunctionResponse.Response)
				}

				// For tool results, Ollama expects a specific format
				msg = Message{
					Role:    "tool",
					Content: responseContent,
				}
				messages = append(messages, msg)
				hasContent = false // Skip adding this message again
				break
			}
		}

		// If we have tool calls, set them in the ToolCalls field
		if len(toolCalls) > 0 {
			msg.ToolCalls = toolCalls
			hasContent = true
		}

		if hasContent && msg.Role != "tool" {
			// Adjust role name: model -> assistant
			if msg.Role == "model" {
				msg.Role = "assistant"
			}
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

// convertTools converts ADK tools to Ollama format
func (l *LiteLlm) convertTools(tools interface{}) []Tool {
	var ollamaTools []Tool

	// Handle different possible tool representations
	if tools == nil {
		return ollamaTools
	}

	// Try to convert tools to a map of tool names to tool definitions
	toolsMap, ok := tools.(map[string]interface{})
	if !ok {
		// If not a map, try to handle as slice or other type
		return ollamaTools
	}

	for toolName := range toolsMap {
		var toolFunc ToolFunction

		// Currently, we know about get_element_info
		if toolName == "get_element_info" {
			toolFunc = ToolFunction{
				Name:        "get_element_info",
				Description: "Get information about a chemical element by symbol or Chinese name.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"symbol": map[string]interface{}{
							"type":        "string",
							"description": "The element symbol (e.g., H, O, Fe, Au) or Chinese name (e.g., 氢, 氧, 铁, 金)",
						},
					},
					"required": []string{"symbol"},
				},
			}
		} else {
			// Generic tool definition for unknown tools
			toolFunc = ToolFunction{
				Name:        toolName,
				Description: "Tool: " + toolName,
				Parameters: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			}
		}

		ollamaTools = append(ollamaTools, Tool{
			Type:     "function",
			Function: toolFunc,
		})
	}

	return ollamaTools
}

// sendChatRequest sends a request to Ollama /api/chat endpoint
func (l *LiteLlm) sendChatRequest(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/api/chat", l.baseURL)

	// Convert request to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Ollama chat request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama: %v", err)
	}
	defer httpResp.Body.Close()

	// Check status code
	if httpResp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("Ollama API returned error: %s, response: %s", httpResp.Status, string(respBody))
	}

	// Read and parse response
	var chatResp ChatResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse Ollama chat response: %v", err)
	}

	return &chatResp, nil
}

// convertToADKResponse converts Ollama chat response to ADK response format
func (l *LiteLlm) convertToADKResponse(resp *ChatResponse) (*model.LLMResponse, error) {
	var parts []*genai.Part

	// Check if response contains tool calls
	if len(resp.Message.ToolCalls) > 0 {
		for _, toolCall := range resp.Message.ToolCalls {
			part := &genai.Part{
				FunctionCall: &genai.FunctionCall{
					Name: toolCall.Function.Name,
					Args: toolCall.Function.Arguments,
				},
			}
			parts = append(parts, part)
		}
	} else if resp.Message.Content != "" {
		// Regular text response
		part := &genai.Part{
			Text: resp.Message.Content,
		}
		parts = append(parts, part)
	}

	// Create ADK content
	content := &genai.Content{
		Role:  "assistant",
		Parts: parts,
	}

	// Create ADK response
	return &model.LLMResponse{
		Content:      content,
		TurnComplete: true,
		Partial:      false,
	}, nil
}
