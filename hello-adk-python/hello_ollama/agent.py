import logging
from datetime import datetime
from google.adk.models.lite_llm import LiteLlm
from google.adk.agents import Agent
from google.adk.tools.tool_context import ToolContext

from hello_ollama.periodic_table import periodic_map


def get_element_info(symbol: str, tool_context: ToolContext = None) -> dict:
    """Returns information about a chemical element by its Chinese name."""
    logging.info(f"查询元素信息: {symbol}")
    
    # 获取和使用context信息
    if tool_context:
        # 打印context信息
        logging.info(f"Context信息: {tool_context}")
        
        # 打印session信息
        session_id = tool_context.session.id if hasattr(tool_context, 'session') and hasattr(tool_context.session, 'id') else 'unknown'
        logging.info(f"Session信息: ID={session_id}")
        
        # 打印session state信息
        if hasattr(tool_context.session, 'state'):
            session_state = tool_context.session.state
            logging.info(f"Session State信息: {session_state}")
        
        # 在会话状态中存储查询历史
        if hasattr(tool_context.session, 'state'):
            if isinstance(tool_context.session.state, dict):
                if 'query_history' not in tool_context.session.state:
                    tool_context.session.state['query_history'] = []
                tool_context.session.state['query_history'].append(symbol)
                logging.info(f"查询历史: {tool_context.session.state['query_history']}")
            else:
                if not hasattr(tool_context.session.state, 'query_history'):
                    tool_context.session.state.query_history = []
                tool_context.session.state.query_history.append(symbol)
                logging.info(f"查询历史: {tool_context.session.state.query_history}")
        
        # 在delta-aware state中添加数据，会在web界面的State tab中显示
        tool_context.state['last_query'] = symbol
        tool_context.state['last_query_time'] = datetime.now().isoformat()
        
        # 保存查询结果作为artifact
        element_info = periodic_map().get(symbol, {})
        if element_info:
            import json
            from google.genai import types
            
            # 创建一个包含元素信息的文本artifact
            element_data = {
                "symbol": symbol,
                "name": element_info.get("name", ""),
                "chinese_name": element_info.get("chinese_name", ""),
                "atomic_number": element_info.get("atomic_number", 0),
                "atomic_weight": element_info.get("atomic_weight", 0)
            }
            
            # 保存为JSON格式的artifact
            json_content = json.dumps(element_data, indent=2, ensure_ascii=False)
            
            # 使用二进制数据格式存储，符合Artifacts的最佳实践
            json_bytes = json_content.encode('utf-8')
            artifact = types.Part(
                inline_data=types.Blob(
                    mime_type="application/json",  # 正确设置MIME类型
                    data=json_bytes
                )
            )
            
            try:
                import asyncio
                loop = asyncio.get_event_loop()
                if loop.is_running():
                    # 如果已经在事件循环中，直接调用
                    loop.create_task(tool_context.save_artifact(f"element_info_{symbol}.json", artifact))
                else:
                    # 如果不在事件循环中，运行一个新的事件循环
                    asyncio.run(tool_context.save_artifact(f"element_info_{symbol}.json", artifact))
                logging.info(f"已保存元素信息artifact: element_info_{symbol}.json")
            except Exception as e:
                logging.error(f"保存artifact失败: {e}")
    
    if symbol in periodic_map():
        element = periodic_map()[symbol]
    else:
        logging.warning(f"元素符号或中文名 '{symbol}' 未找到。")
        return {
            "status": "error",
            "error_message": f"元素符号或中文名 '{symbol}' 未找到。",
        }

    chinese_name = element["chinese_name"]
    element_name = element["name"]
    atomic_number = element["atomic_number"]
    atomic_weight_value = element["atomic_weight"]
    report = f"{chinese_name}（{element_name}），原子序数：{atomic_number}，原子量：{atomic_weight_value}"
    logging.info(f"查询成功: {report}")
    return {"status": "success", "report": report}


# ollama
adk_model = LiteLlm(model="ollama_chat/qwen2.5")
# google
# adk_model="gemini-2.0-flash",
root_agent = Agent(
    model=adk_model,
    name="ollama_agent",
    description=("Agent to answer questions about chemical element information."),
    instruction=("You are a helpful agent."),
    tools=[get_element_info],
)

if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s"
    )
