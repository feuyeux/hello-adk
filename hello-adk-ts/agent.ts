import { FunctionTool, LlmAgent } from '@google/adk';
import { z } from 'zod';
import { getElement } from './periodic-table';
import { LiteLlm } from './lite-llm';

/* Periodic Table Element Info Tool */
const getElementInfoTool = new FunctionTool({
  name: 'get_element_info',
  description: 'Get information about a chemical element by symbol or Chinese name.',
  parameters: z.object({
    symbol: z
      .string()
      .describe('Element symbol (e.g., "H", "O") or Chinese name (e.g., "氢", "氧")'),
  }),
  execute: ({ symbol }: { symbol: string }) => {
    console.log(`查询元素信息: ${symbol}`);
    
    const element = getElement(symbol);
    
    if (!element) {
      const errorMsg = `元素符号或中文名称 '${symbol}' 未找到。`;
      console.warn(errorMsg);
      return {
        status: 'error',
        error_message: errorMsg,
      };
    }

    const report = `${element.chineseName}（${element.name}），原子序数：${element.atomicNumber}，原子量：${element.atomicWeight}`;
    console.log(`查询成功: ${report}`);
    
    return {
      status: 'success',
      report: report,
    };
  },
});

// Use Ollama Qwen2.5 model via LiteLlm
const adkModel = new LiteLlm('ollama_chat/qwen2.5');

export const rootAgent = new LlmAgent({
  name: 'ollama_agent',
  model: adkModel,
  description: 'Agent to answer questions about chemical element information.',
  instruction: `You are a helpful agent that provides information about chemical elements.
Use the 'get_element_info' tool to look up element information by symbol or Chinese name.
- Call this tool ONLY ONCE per element.
- Do NOT call the tool multiple times for the same element.
- Use the information returned by the tool to answer the user's question.
- If you already have the information about an element from a previous tool call, do NOT call the tool again.
- Only call the tool if you need additional information that you don't already have.
Always use the tool to get accurate information.`,
  tools: [getElementInfoTool],
});
