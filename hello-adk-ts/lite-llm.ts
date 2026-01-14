import {
  BaseLlm,
  BaseLlmConnection,
  LlmRequest,
  LlmResponse,
} from "@google/adk";
import { Content, Blob as GenAIBlob } from "@google/genai";

export class LiteLlm extends BaseLlm {
  constructor(modelName: string) {
    super({
      model: modelName,
    });
  }

  async *generateContentAsync(
    llmRequest: LlmRequest,
    stream?: boolean,
  ): AsyncGenerator<LlmResponse, void> {
    // Convert ADK request to Ollama API format
    const messages = llmRequest.contents
      .map((content) => {
        // Handle different part types
        let messageContent: any = null;
        let toolCalls: any[] = [];
        let toolResult: any = null;

        let isToolResponse = false;
        let toolCallId = "";
        let toolResponseContent = "";

        for (const part of content.parts || []) {
          if ("text" in part && part.text) {
            messageContent = part.text;
          } else if ("functionCall" in part && part.functionCall) {
            // Handle function call parts
            // ADK doesn't provide a callId, so we generate one
            const callId = `call_${Date.now()}_${Math.floor(Math.random() * 1000)}`;
            toolCalls.push({
              function: {
                name: part.functionCall.name,
                arguments: part.functionCall.args, // Directly pass the object, not stringified
              },
              id: callId,
              type: "function",
            });
          } else if ("functionResponse" in part && part.functionResponse) {
            // Handle function response parts - format for Ollama
            // Ollama expects tool results to be in this format:
            // {
            //   "role": "tool",
            //   "content": "result",
            //   "tool_call_id": "call_id"
            // }
            isToolResponse = true;
            toolCallId = `call_${Date.now()}_${Math.floor(Math.random() * 1000)}`; // Generate a consistent call ID

            const functionResponse = part.functionResponse;
            const response = functionResponse.response;

            if (response && "status" in response) {
              if (response.status === "success" && "report" in response) {
                toolResponseContent = String(response.report);
              } else if (
                response.status === "error" &&
                "error_message" in response
              ) {
                toolResponseContent = String(response.error_message);
              }
            }
          }
        }

        // Create message based on content type
        const message: any = {};

        if (isToolResponse) {
          // Format tool response according to Ollama API expectations
          message.role = "tool";
          message.content = toolResponseContent;
          message.tool_call_id = toolCallId;
        } else {
          // Regular message format
          message.role = content.role === "model" ? "assistant" : content.role;

          if (messageContent) {
            message.content = messageContent;
          }

          if (toolCalls.length > 0) {
            message.tool_calls = toolCalls;
          }
        }

        return message;
      })
      .filter(
        (msg) =>
          (msg.content !== null &&
            msg.content !== undefined &&
            msg.content !== "") ||
          (msg.tool_calls && msg.tool_calls.length > 0) ||
          (msg.role === "tool" && msg.tool_call_id),
      );

    // Convert tools to Ollama format
    const tools = Object.values(llmRequest.toolsDict)
      .map((tool) => {
        // Get function declaration from the tool
        const declaration = tool._getDeclaration();
        if (!declaration || !declaration.parameters) {
          return null;
        }

        // Convert GenAI parameter types to standard JSON Schema types
        const convertParamType = (type: string) => {
          switch (type) {
            case "OBJECT":
              return "object";
            case "STRING":
              return "string";
            case "NUMBER":
              return "number";
            case "BOOLEAN":
              return "boolean";
            default:
              return type.toLowerCase();
          }
        };

        // Convert parameter properties
        const properties: any = {};
        for (const [key, param] of Object.entries(
          declaration.parameters.properties || {},
        )) {
          properties[key] = {
            type: convertParamType(param.type || "STRING"),
            description: param.description,
          };
        }

        // Create Ollama-compatible function schema
        const functionSchema = {
          name: declaration.name,
          description: declaration.description,
          parameters: {
            type: "object",
            properties: properties,
            required: declaration.parameters.required || [],
          },
        };

        return {
          type: "function",
          function: functionSchema,
        };
      })
      .filter((tool) => tool !== null);

    try {
      // Call Ollama API
      const response = await fetch("http://localhost:11434/api/chat", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          model: this.model.replace("ollama_chat/", ""),
          messages,
          tools: tools,
          stream: stream || false,
          temperature: 0.7,
        }),
      });

      // Get response text first for debugging
      const responseText = await response.text();
      console.log("Ollama API Response:", responseText);

      if (!response.ok) {
        throw new Error(
          `Ollama API error: ${response.status} ${response.statusText}\nResponse: ${responseText}`,
        );
      }

      let data;
      try {
        data = JSON.parse(responseText);
      } catch (jsonError) {
        console.error(
          "Failed to parse Ollama API response as JSON:",
          jsonError,
        );
        throw new Error(`Ollama API returned invalid JSON: ${responseText}`);
      }

      // Convert Ollama response to ADK format
      let result: LlmResponse;

      // Check if the response contains tool calls
      if (
        data.message &&
        data.message.tool_calls &&
        data.message.tool_calls.length > 0
      ) {
        // Create a function call part for ADK
        const functionCall = data.message.tool_calls[0].function;

        // Parse arguments if they're a string
        let args;
        try {
          args =
            typeof functionCall.arguments === "string"
              ? JSON.parse(functionCall.arguments)
              : functionCall.arguments;
        } catch (parseError) {
          console.error("Failed to parse function arguments:", parseError);
          args = functionCall.arguments;
        }

        const functionCallPart = {
          functionCall: {
            name: functionCall.name,
            args: args,
          },
        };

        result = {
          content: {
            role: "assistant",
            parts: [functionCallPart as any],
          },
          turnComplete: true,
          partial: false,
        };
      } else if (data.message) {
        // Regular text response
        result = {
          content: {
            role: "assistant",
            parts: [{ text: data.message.content || "" }],
          },
          turnComplete: true,
          partial: false,
        };
      } else {
        // Handle unexpected response format
        result = {
          content: {
            role: "assistant",
            parts: [
              {
                text: "Sorry, I encountered an error processing your request.",
              },
            ],
          },
          turnComplete: true,
          partial: false,
        };
      }

      yield result;
    } catch (error) {
      console.error("Error in generateContentAsync:", error);
      // Return a friendly error message instead of throwing
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      yield {
        content: {
          role: "assistant",
          parts: [{ text: `Sorry, I encountered an error: ${errorMessage}` }],
        },
        turnComplete: true,
        partial: false,
      };
    }
  }

  async connect(llmRequest: LlmRequest): Promise<BaseLlmConnection> {
    // Create a simple BaseLlmConnection implementation
    return {
      async sendHistory(history: Content[]): Promise<void> {
        // Not implemented for this simple case
      },

      async sendContent(content: Content): Promise<void> {
        // Not implemented for this simple case
      },

      async sendRealtime(blob: GenAIBlob): Promise<void> {
        // Not implemented for this simple case
      },

      async *receive(): AsyncGenerator<LlmResponse, void, void> {
        // Not implemented for this simple case
      },

      async close(): Promise<void> {
        // Not implemented for this simple case
      },
    };
  }
}

// Register the LiteLlm class for use with Ollama models
BaseLlm.supportedModels.push(/ollama_chat\/.*/);
