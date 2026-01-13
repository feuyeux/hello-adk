package main

import (
	"bufio"
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

// LiteLlm implements the model.LLM interface for Ollama models
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

// GenerateContent generates content using the Ollama model
func (l *LiteLlm) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		// Convert ADK request to Ollama request format
		ollamaReq, err := l.convertToOllamaRequest(req)
		if err != nil {
			yield(nil, err)
			return
		}

		// Send request to Ollama
		ollamaResp, err := l.sendRequest(ctx, ollamaReq)
		if err != nil {
			yield(nil, err)
			return
		}

		// Convert Ollama response to ADK response
		adkResp, err := l.convertToADKResponse(ollamaResp)
		if err != nil {
			yield(nil, err)
			return
		}

		// Yield the single response
		yield(adkResp, nil)
	}
}

// OllamaCompletionRequest represents the request format for Ollama completion API
type OllamaCompletionRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaCompletionResponse represents the response format from Ollama completion API
type OllamaCompletionResponse struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Response string `json:"response"`
	Done bool `json:"done"`
}

// convertToOllamaRequest converts ADK LLMRequest to Ollama completion request format
func (l *LiteLlm) convertToOllamaRequest(req *model.LLMRequest) (*OllamaCompletionRequest, error) {
	// Build prompt from all messages
	var promptBuilder strings.Builder

	// Add system prompt - ensure it's only added once
	systemPromptAdded := false

	// Convert messages
	for _, content := range req.Contents {
		// Extract text content from parts
		var textContent strings.Builder
		for _, part := range content.Parts {
			// Handle text parts
			if part.Text != "" {
				textContent.WriteString(part.Text)
				textContent.WriteString("\n")
			}
		}

		// Remove trailing newlines and whitespace
		cleanContent := strings.TrimSpace(textContent.String())
		if cleanContent == "" {
			continue // Skip empty messages
		}

		// Check if this is a system message
		if content.Role == "system" && !systemPromptAdded {
			promptBuilder.WriteString(cleanContent)
			promptBuilder.WriteString("\n\n")
			systemPromptAdded = true
		} else if content.Role == "user" {
			// For user messages, add them directly
			promptBuilder.WriteString("用户问题: ")
			promptBuilder.WriteString(cleanContent)
			promptBuilder.WriteString("\n")
			promptBuilder.WriteString("请使用get_element_info工具来查询相关元素信息。\n")
		}
	}

	// If no system prompt was added, add default one
	if !systemPromptAdded {
		promptBuilder.WriteString("You are a helpful agent that provides information about chemical elements.\n")
		promptBuilder.WriteString("请使用get_element_info工具来查询相关元素信息。\n\n")
	}

	return &OllamaCompletionRequest{
		Model:  l.modelName,
		Prompt: promptBuilder.String(),
	}, nil
}

// sendRequest sends a request to the Ollama completion API
func (l *LiteLlm) sendRequest(ctx context.Context, req *OllamaCompletionRequest) (*OllamaCompletionResponse, error) {
	url := fmt.Sprintf("%s/api/generate", l.baseURL)

	// Convert request to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Ollama request: %v", err)
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

	// Read response line by line and accumulate the response
	var fullResponse strings.Builder
	var lastResponse OllamaCompletionResponse

	// Use scanner to read line by line
	scanner := bufio.NewScanner(httpResp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse the line
		var partialResp OllamaCompletionResponse
		if err := json.Unmarshal([]byte(line), &partialResp); err != nil {
			return nil, fmt.Errorf("failed to parse Ollama response line: %v\nLine: %s", err, line)
		}

		// Accumulate the response text
		if partialResp.Response != "" {
			fullResponse.WriteString(partialResp.Response)
		}

		// Save the last response (which contains the done flag)
		lastResponse = partialResp
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read Ollama response: %v", err)
	}

	// Set the accumulated response text
	lastResponse.Response = fullResponse.String()

	return &lastResponse, nil
}

// convertToADKResponse converts Ollama completion response to ADK response format
func (l *LiteLlm) convertToADKResponse(resp *OllamaCompletionResponse) (*model.LLMResponse, error) {
	// Create ADK content part
	part := &genai.Part{
		Text: resp.Response,
	}

	// Create ADK content
	content := &genai.Content{
		Role:  "assistant",
		Parts: []*genai.Part{part},
	}

	// Create ADK response
	return &model.LLMResponse{
		Content:      content,
		TurnComplete: true,
		Partial:      false,
	}, nil
}
