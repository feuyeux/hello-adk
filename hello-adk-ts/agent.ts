import { FunctionTool, LlmAgent } from '@google/adk';
import { z } from 'zod';
import { getElement } from './periodic-table';

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
      const errorMsg = `Element symbol or Chinese name '${symbol}' not found.`;
      console.warn(errorMsg);
      return {
        status: 'error',
        error_message: errorMsg,
      };
    }

    const report = `${element.chineseName}(${element.name}), Atomic Number: ${element.atomicNumber}, Atomic Weight: ${element.atomicWeight}`;
    console.log(`查询成功: ${report}`);
    
    return {
      status: 'success',
      report: report,
    };
  },
});

export const rootAgent = new LlmAgent({
  name: 'periodic_table_agent',
  model: 'gemini-2.5-flash',
  description: 'Agent to answer questions about chemical element information.',
  instruction: `You are a helpful agent that provides information about chemical elements.
Use the 'get_element_info' tool to look up element information by symbol or Chinese name.`,
  tools: [getElementInfoTool],
});
