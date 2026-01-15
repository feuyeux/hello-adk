# ADK 中 Context、Session 和 Session State 

| 概念 | 生命周期 | 存储方式 | 作用范围 | 主要用途 |
|------|----------|----------|----------|----------|
| Context | 单次请求 | 临时 | 请求级 | 传递请求信息、控制执行 |
| Session | 多次请求 | 持久化 | 用户级 | 管理用户会话、存储对话历史 |
| Session State | 会话期间 | 可持久化 | 会话级 | 记录执行状态、工具调用结果 |

## 1. Context (上下文)

### 1.1 定义与作用
- **定义**：Context 是函数调用的执行环境，主要用于传递**请求级别**的信息
- **作用**：
  - 控制函数执行的生命周期、超时和取消
  - 传递请求元数据（如请求 ID、用户信息等）
  - 跨函数调用传递上下文信息

### 1.2 特点
- 临时的，仅在特定请求期间存在
- 不持久化存储，请求结束后即销毁
- 用于单次请求的上下文传递
- 可包含请求级别的配置和元数据

### 1.3 示例代码 (Go 版本)
```go
func (l *LiteLlm) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
    // 使用 ctx 控制请求超时和取消
    // 从 ctx 中获取请求元数据
}
```

### 1.4 示例代码 (Python 版本)
```python
from google.adk.tools.tool_context import ToolContext
import logging

def get_element_info(symbol: str, tool_context: ToolContext = None) -> dict:
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
```

### 1.5 官方文档
- [Context Overview](https://google.github.io/adk-docs/context/)
- [Context Caching](https://google.github.io/adk-docs/context/caching/)
- [Context Compaction](https://google.github.io/adk-docs/context/compaction/)

## 2. Session (会话)

### 2.1 定义与作用
- **定义**：Session 是用户与 Agent 之间的完整交互会话，代表一次完整的对话过程
- **作用**：
  - 管理用户身份和会话标识
  - 存储完整的对话历史记录
  - 维护长期的用户-Agent 交互状态
  - 支持会话重放和回溯

### 2.2 特点
- 持久化存储，跨越多个请求
- 与特定用户关联，包含用户身份信息
- 包含完整的对话历史记录
- 支持会话管理（创建、更新、关闭等）

### 2.3 示例代码 (Java 版本)
```java
Session session = runner
        .sessionService()
        .createSession(runner.appName(), "user1234")
        .blockingGet();
```

### 2.4 示例代码 (Python 版本)
```python
from google.adk.tools.tool_context import ToolContext
import logging

def get_element_info(symbol: str, tool_context: ToolContext = None) -> dict:
    logging.info(f"查询元素信息: {symbol}")
    
    # 获取和使用context信息
    if tool_context:
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
```

### 2.5 官方文档
- [Sessions Overview](https://google.github.io/adk-docs/sessions/)
- [Session Management](https://google.github.io/adk-docs/sessions/session/)
- [Session Rewind](https://google.github.io/adk-docs/sessions/rewind/)
- [Session Memory](https://google.github.io/adk-docs/sessions/memory/)

## 3. Session State (会话状态)

### 3.1 定义与作用
- **定义**：Session State 是 Session 内部的状态数据，记录 Agent 执行过程中的中间状态和结果
- **作用**：
  - 记录 Agent 执行过程中的工具调用状态
  - 存储对话过程中的中间结果
  - 维护 Agent 的执行上下文
  - 支持会话状态的持久化和恢复

### 3.2 特点
- 与特定 Session 关联，是 Session 的一部分
- 通常是临时的，但可以随 Session 一起持久化
- 包含执行状态、工具调用记录、对话历史等
- 支持状态的更新和查询

### 3.3 示例代码 (TypeScript 版本)
在实际应用中，Session State 通常通过 ADK 框架内部管理，以下是概念性示例：
```typescript
interface SessionState {
  currentToolCall: ToolCall | null;
  conversationHistory: Message[];
  executionPath: string[];
  agentConfig: AgentConfig;
}
```

### 3.4 示例代码 (Python 版本)
```python
from google.adk.tools.tool_context import ToolContext
from datetime import datetime
import logging

def get_element_info(symbol: str, tool_context: ToolContext = None) -> dict:
    logging.info(f"查询元素信息: {symbol}")
    
    if tool_context:
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
```

### 3.5 实际应用示例 - 元素查询历史
以下是一个实际的Python实现，展示了如何使用Session State来管理元素查询历史：

```python
from google.adk.tools.tool_context import ToolContext
from datetime import datetime
import logging

def get_element_info(symbol: str, tool_context: ToolContext = None) -> dict:
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
    
    # 元素查询逻辑...
```

### 3.6 官方文档
- [Session State](https://google.github.io/adk-docs/sessions/state/)

## 4. Artifacts (工件)

### 4.1 定义与作用
- **定义**：Artifacts是与会话关联的命名、版本化二进制数据
- **作用**：
  - 持久化存储工具调用的结果
  - 保存会话期间生成的文件
  - 支持数据的版本控制和检索
  - 便于在web界面中查看和下载

### 4.2 特点
- 与特定Session关联
- 持久化存储，可以在会话结束后继续访问
- 支持多种MIME类型（JSON、文本、图像等）
- 在web界面的Artifacts tab中显示

### 4.3 示例代码 (Python 版本)
```python
from google.adk.tools.tool_context import ToolContext
from google.genai import types
import json

def get_element_info(symbol: str, tool_context: ToolContext = None) -> dict:
    # 查询元素信息
    element_info = periodic_map().get(symbol, {})
    
    if tool_context and element_info:
        # 将元素信息保存为Artifact
        element_data = {
            "symbol": symbol,
            "name": element_info.get("name", ""),
            "chinese_name": element_info.get("chinese_name", ""),
            "atomic_number": element_info.get("atomic_number", 0),
            "atomic_weight": element_info.get("atomic_weight", 0)
        }
        
        # 保存为JSON格式
        json_content = json.dumps(element_data, indent=2, ensure_ascii=False)
        json_bytes = json_content.encode('utf-8')
        
        # 创建Artifact
        artifact = types.Part(
            inline_data=types.Blob(
                mime_type="application/json",
                data=json_bytes
            )
        )
        
        try:
            import asyncio
            loop = asyncio.get_event_loop()
            if loop.is_running():
                # 如果已经在事件循环中，使用create_task
                loop.create_task(tool_context.save_artifact(f"element_info_{symbol}.json", artifact))
            else:
                # 如果不在事件循环中，直接运行
                asyncio.run(tool_context.save_artifact(f"element_info_{symbol}.json", artifact))
            logging.info(f"已保存元素信息artifact: element_info_{symbol}.json")
        except Exception as e:
            logging.error(f"保存artifact失败: {e}")
    
    return {"status": "success", "report": f"查询到元素 {symbol} 的信息"}

### 4.4 官方文档
- [Artifacts Overview](https://google.github.io/adk-docs/artifacts/)
