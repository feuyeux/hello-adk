package com.example;

import com.google.adk.models.BaseLlm;
import com.google.adk.models.BaseLlmConnection;
import com.google.adk.models.LlmRequest;
import com.google.adk.models.LlmResponse;
import com.google.genai.types.*;
import dev.langchain4j.agent.tool.ToolSpecification;
import dev.langchain4j.agent.tool.ToolParameters;
import dev.langchain4j.data.message.*;
import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.model.output.Response;
import io.reactivex.rxjava3.core.Flowable;

import java.util.*;

/**
 * LangChain4j Ollama 适配器 - 将 LangChain4j 的 Ollama 模型包装为 ADK 的 BaseLlm
 * 支持 Qwen2.5 等模型的 function calling
 */
public class LangChain4jOllamaLlm extends BaseLlm {
    private final ChatLanguageModel chatModel;

    public LangChain4jOllamaLlm(ChatLanguageModel chatModel, String modelName) {
        super(modelName);
        this.chatModel = chatModel;
        System.out.println("✅ LangChain4j Ollama 适配器已初始化: " + modelName);
    }

    @Override
    public Flowable<LlmResponse> generateContent(LlmRequest request, boolean streaming) {
        return Flowable.fromCallable(() -> {
            try {
                // 1. 从 request 中提取工具定义
                List<ToolSpecification> toolSpecs = extractToolSpecs(request);
                
                // 2. 转换消息
                List<ChatMessage> messages = convertMessages(request);
                
                System.out.println("📤 发送到 Ollama: " + messages.size() + " 条消息, " + toolSpecs.size() + " 个工具");
                
                // 3. 调用 Ollama
                Response<AiMessage> response;
                if (!toolSpecs.isEmpty()) {
                    response = chatModel.generate(messages, toolSpecs);
                } else {
                    response = chatModel.generate(messages);
                }
                
                System.out.println("📥 收到 Ollama 响应");
                
                // 4. 转换响应
                return convertResponse(response);
            } catch (Exception e) {
                System.err.println("❌ Ollama 调用失败: " + e.getMessage());
                e.printStackTrace();
                throw new RuntimeException("Ollama 调用失败", e);
            }
        });
    }

    @Override
    public BaseLlmConnection connect(LlmRequest request) {
        throw new UnsupportedOperationException("Ollama 不支持流式连接");
    }

    /**
     * 从 LlmRequest 中提取工具规范
     */
    private List<ToolSpecification> extractToolSpecs(LlmRequest request) {
        List<ToolSpecification> toolSpecs = new ArrayList<>();
        
        if (request.tools() == null || request.tools().isEmpty()) {
            System.out.println("🔧 没有工具");
            return toolSpecs;
        }
        
        System.out.println("🔧 发现 " + request.tools().size() + " 个工具");
        
        for (var entry : request.tools().entrySet()) {
            String toolName = entry.getKey();
            var tool = entry.getValue();
            
            System.out.println("  工具: " + toolName);
            
            // FunctionTool 实现了 declaration() 方法，返回 Optional<FunctionDeclaration>
            // 但由于 AutoValue 生成类的访问限制，我们直接用工具名和描述创建 ToolSpecification
            // 并从 @Schema 注解获取参数信息
            
            // 对于 getElementInfo 工具，我们知道它有一个 String 参数 "symbol"
            // 这是一个硬编码方案，但可以避免反射权限问题
            if (toolName.equals("getElementInfo")) {
                Map<String, Map<String, Object>> properties = new HashMap<>();
                properties.put("symbol", Map.of(
                    "type", "string",
                    "description", "The element identifier - can be symbol (H, O, Au), Chinese name (氢, 氧, 金), or English name (hydrogen, oxygen, gold)"
                ));
                
                ToolSpecification spec = ToolSpecification.builder()
                    .name("getElementInfo")
                    .description("Get detailed information about a chemical element. Use this tool for ANY question about chemical elements.")
                    .parameters(ToolParameters.builder()
                        .properties(properties)
                        .required(List.of("symbol"))
                        .build())
                    .build();
                
                toolSpecs.add(spec);
                System.out.println("  ✓ 已添加工具: getElementInfo (symbol: string)");
            } else {
                // 对于其他工具，创建基本定义
                toolSpecs.add(ToolSpecification.builder()
                    .name(toolName)
                    .description("Tool: " + toolName)
                    .build());
                System.out.println("  ✓ 已添加工具: " + toolName + " (无参数定义)");
            }
        }
        
        return toolSpecs;
    }

    /**
     * 转换 ADK 消息为 LangChain4j ChatMessage
     */
    private List<ChatMessage> convertMessages(LlmRequest request) {
        List<ChatMessage> messages = new ArrayList<>();
        
        if (request.contents() == null) {
            return messages;
        }
        
        for (com.google.genai.types.Content content : request.contents()) {
            Optional<List<Part>> partsOpt = content.parts();
            if (partsOpt.isEmpty()) continue;
            
            boolean hasContent = false;
            
            for (Part part : partsOpt.get()) {
                // 处理文本部分
                if (part.text().isPresent()) {
                    String text = part.text().get();
                    String role = content.role().orElse("user");
                    switch (role) {
                        case "user" -> messages.add(new UserMessage(text));
                        case "model" -> messages.add(new AiMessage(text));
                        case "system" -> messages.add(new SystemMessage(text));
                    }
                    hasContent = true;
                }
                
                // 处理工具调用结果（function response）
                if (part.functionResponse().isPresent()) {
                    var funcResponse = part.functionResponse().get();
                    String funcName = funcResponse.name().orElse("unknown_tool");
                    
                    // 获取工具调用结果内容
                    String resultText = "{}";
                    // 尝试获取FunctionResponse的内容
                    try {
                        // 使用反射获取可能的方法
                        var method = funcResponse.getClass().getMethod("content");
                        var funcContent = method.invoke(funcResponse);
                        if (funcContent instanceof Optional) {
                            Optional<?> optContent = (Optional<?>) funcContent;
                            if (optContent.isPresent()) {
                                resultText = optContent.get().toString();
                            }
                        } else {
                            resultText = funcContent.toString();
                        }
                    } catch (Exception e) {
                        // 如果反射失败，将FunctionResponse直接转换为字符串
                        resultText = funcResponse.toString();
                    }
                    
                    // 将工具调用结果转换为 LangChain4j 的 ToolExecutionResultMessage
                    // 注意：ToolExecutionResultMessage 需要工具名、执行ID和结果
                    messages.add(new ToolExecutionResultMessage(funcName, UUID.randomUUID().toString(), resultText));
                    hasContent = true;
                    System.out.println("🔧 工具调用结果: " + funcName + " -> " + resultText);
                }
            }
            
            if (!hasContent) continue;
        }
        
        return messages;
    }

    /**
     * 转换 LangChain4j 响应为 ADK LlmResponse
     */
    private LlmResponse convertResponse(Response<AiMessage> response) {
        AiMessage aiMessage = response.content();
        List<Part> responseParts = new ArrayList<>();
        
        // 检查是否有 function call
        if (aiMessage.hasToolExecutionRequests()) {
            System.out.println("🔧 检测到 function call!");
            
            for (dev.langchain4j.agent.tool.ToolExecutionRequest toolReq : aiMessage.toolExecutionRequests()) {
                String funcName = toolReq.name();
                String argsJson = toolReq.arguments();
                
                System.out.println("  → 函数: " + funcName);
                System.out.println("  → 参数: " + argsJson);
                
                Map<String, Object> args = parseArgs(argsJson);
                Part funcPart = Part.fromFunctionCall(funcName, args);
                responseParts.add(funcPart);
            }
        } else if (aiMessage.text() != null && !aiMessage.text().isEmpty()) {
            responseParts.add(Part.fromText(aiMessage.text()));
        }
        
        com.google.genai.types.Content responseContent = com.google.genai.types.Content.builder()
            .role("model")
            .parts(responseParts)
            .build();
        
        return LlmResponse.builder()
            .content(responseContent)
            .turnComplete(true)
            .partial(false)
            .build();
    }

    /**
     * 解析 JSON 参数字符串
     */
    private Map<String, Object> parseArgs(String jsonStr) {
        Map<String, Object> result = new HashMap<>();
        if (jsonStr == null || jsonStr.isEmpty()) return result;
        
        try {
            jsonStr = jsonStr.trim();
            if (jsonStr.startsWith("{")) jsonStr = jsonStr.substring(1);
            if (jsonStr.endsWith("}")) jsonStr = jsonStr.substring(0, jsonStr.length() - 1);
            jsonStr = jsonStr.trim();
            
            if (jsonStr.isEmpty()) return result;
            
            String[] pairs = jsonStr.split(",");
            for (String pair : pairs) {
                String[] kv = pair.split(":", 2);
                if (kv.length == 2) {
                    String key = kv[0].trim().replaceAll("\"", "");
                    String value = kv[1].trim().replaceAll("\"", "");
                    result.put(key, value);
                }
            }
        } catch (Exception e) {
            System.err.println("⚠️ 参数解析失败: " + e.getMessage());
        }
        
        return result;
    }
}
