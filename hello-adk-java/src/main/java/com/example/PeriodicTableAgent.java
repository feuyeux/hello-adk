package com.example;

import com.google.adk.agents.BaseAgent;
import com.google.adk.agents.LlmAgent;
import com.google.adk.tools.Annotations.Schema;
import com.google.adk.tools.FunctionTool;
import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.model.ollama.OllamaChatModel;

import java.time.Duration;
import java.util.Map;

public class PeriodicTableAgent {
    // Lazy initialization to avoid issues during CheckApi
    private static BaseAgent rootAgent;
    
    public static BaseAgent getRootAgent() {
        if (rootAgent == null) {
            rootAgent = initAgent();
        }
        return rootAgent;
    }
    
    public static BaseAgent ROOT_AGENT = null; // Will be initialized on first access

    private static String getEnvOrDefault(String key, String defaultValue) {
        String value = System.getenv(key);
        return (value == null || value.isEmpty()) ? defaultValue : value;
    }

    private static BaseAgent initAgent() {
        // 使用 LangChain4j 封装 Ollama
        String ollamaBaseUrl = getEnvOrDefault("OLLAMA_API_BASE", "http://localhost:11434");
        String ollamaModel = getEnvOrDefault("OLLAMA_MODEL", "qwen2.5");
        
        System.out.println("Initializing Periodic Table Agent with LangChain4j Ollama");
        System.out.println("Model: " + ollamaModel);
        System.out.println("Base URL: " + ollamaBaseUrl);
        
        // 创建 LangChain4j Ollama 模型
        ChatLanguageModel chatModel = OllamaChatModel.builder()
                .baseUrl(ollamaBaseUrl)
                .modelName(ollamaModel)
                .temperature(0.7)
                .timeout(Duration.ofSeconds(120))
                .build();
        
        // 使用适配器包装为 ADK BaseLlm
        LangChain4jOllamaLlm adkModel = new LangChain4jOllamaLlm(chatModel, ollamaModel);
        
        return LlmAgent.builder()
            .name("periodic-table-agent")
            .description("Agent to answer questions about chemical element information")
            .instruction("""
                You are a helpful agent that provides information about chemical elements from the periodic table.
                
                CRITICAL RULES:
                1. You MUST ALWAYS use the 'getElementInfo' tool for ANY question about chemical elements
                2. NEVER answer from your own knowledge about elements
                3. ALWAYS call getElementInfo first, then use its response to answer
                
                When a user asks about ANY chemical element:
                - Step 1: Call getElementInfo with the element symbol, Chinese name, or English name
                - Step 2: Wait for the tool response
                - Step 3: Return the exact information from the tool
                
                Examples of when to call getElementInfo:
                ✓ "Tell me about hydrogen" -> Call getElementInfo("H") or getElementInfo("hydrogen")
                ✓ "What is 金?" -> Call getElementInfo("金")
                ✓ "Atomic weight of gold" -> Call getElementInfo("Au") or getElementInfo("gold")
                ✓ "Properties of oxygen" -> Call getElementInfo("O") or getElementInfo("oxygen")
                
                Remember: ALWAYS use the tool, never rely on your own knowledge for element information.
                """)
            .model(adkModel)
            .tools(FunctionTool.create(PeriodicTableAgent.class, "getElementInfo"))
            .build();
    }

    @Schema(description = "Get detailed information about a chemical element. IMPORTANT: Use this tool for ANY question about chemical elements. Input can be: 1) Element symbol like 'H', 'O', 'Au', 'Fe'; 2) Chinese name like '氢', '氧', '金', '铁'; 3) English name will be converted to symbol. Returns complete element information including Chinese name, English name, atomic number, and atomic weight.")
    public static Map<String, Object> getElementInfo(
            @Schema(description = "The element identifier - can be symbol (H, O, Au), Chinese name (氢, 氧, 金), or English name (hydrogen, oxygen, gold). Examples: 'H', '氢', 'hydrogen', 'Au', '金', 'gold'") String symbol) {
        
        System.out.println("🔧 Tool called: getElementInfo with symbol=" + symbol);
        
        PeriodicTable table = PeriodicTable.getInstance();
        PeriodicTable.Element element = table.getElement(symbol);
        
        if (element == null) {
            System.out.println("❌ Element not found: " + symbol);
            return Map.of(
                "status", "error",
                "error_message", "Element symbol or Chinese name '" + symbol + "' not found."
            );
        }
        
        String report = String.format(
            "%s(%s), Atomic Number: %d, Atomic Weight: %.4f",
            element.chineseName,
            element.name,
            element.atomicNumber,
            element.atomicWeight
        );
        
        System.out.println("✅ Tool result: " + report);
        
        return Map.of(
            "status", "success",
            "report", report
        );
    }
}
