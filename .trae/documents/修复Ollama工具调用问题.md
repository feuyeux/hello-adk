### 问题分析

在当前的实现中，`lite-llm.ts` 没有正确处理工具调用请求。虽然 `agent.ts` 中定义了 `getElementInfoTool` 工具并将其添加到了 `LlmAgent` 中，但在调用 Ollama API 时，工具信息没有被传递给模型，导致模型不知道如何使用这些工具。

### 解决方案

需要修改 `lite-llm.ts` 文件，实现以下功能：

1. **将工具信息转换为Ollama API格式**：在 `generateContentAsync` 方法中，将 `llmRequest.toolsDict` 转换为 Ollama API 可以理解的工具格式
2. **将工具信息添加到请求中**：将转换后的工具信息包含在发送给 Ollama API 的请求中
3. **处理Ollama返回的工具调用响应**：当 Ollama 返回工具调用请求时，正确解析并返回给 ADK 框架

### 修改步骤

1. 修改 `generateContentAsync` 方法，添加工具转换逻辑
2. 将工具信息添加到 Ollama API 请求中
3. 确保正确处理 Ollama 返回的各种响应类型

### 预期结果

修改后，当执行 `npx @google/adk-devtools run agent.ts` 并询问元素信息时，模型将能够调用 `get_element_info` 工具来获取元素数据，而不是直接尝试回答。