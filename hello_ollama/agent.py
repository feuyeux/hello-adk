import logging
from google.adk.models.lite_llm import LiteLlm
from google.adk.agents import Agent

from hello_ollama.periodic_table import periodic_map


def get_element_info(symbol: str) -> dict:
    """Returns information about a chemical element by its Chinese name."""
    logging.info(f"查询元素信息: {symbol}")
    if symbol in periodic_map():
        element = periodic_map()[symbol]
    else:
        logging.warning(f"元素符号或中文名 '{symbol}' 未找到。")
        return {"status": "error", "error_message": f"元素符号或中文名 '{symbol}' 未找到。"}

    chinese_name = element['chinese_name']
    element_name = element['name']
    atomic_number = element['atomic_number']
    atomic_weight_value = element['atomic_weight']
    report = (
        f"{chinese_name}（{element_name}），原子序数：{atomic_number}，原子量：{atomic_weight_value}"
    )
    logging.info(f"查询成功: {report}")
    return {"status": "success", "report": report}


root_agent = Agent(
    # ollama
    # model=LiteLlm(model="ollama_chat/qwen2.5"),
    # google
    model="gemini-2.0-flash",
    name="ollama_agent",
    description=(
        "Agent to answer questions about chemical element information."
    ),
    instruction=(
        "You are a helpful agent."
    ),
    tools=[get_element_info]
)

if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, format='%(asctime)s %(levelname)s %(message)s')
