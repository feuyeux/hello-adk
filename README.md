# Hello Google Agent Develoopment Kit

1. [Installing ADK](https://google.github.io/adk-docs/get-started/installation/)
2. [Quickstart](https://google.github.io/adk-docs/get-started/quickstart/)

```sh
.venv\Scripts\Activate.ps1
source .venv/Scripts/activate

pip install google-adk
pip install litellm
```

```sh
ollama show qwen2.5
  Model
    architecture        qwen2     
    parameters          7.6B      
    context length      32768     
    embedding length    3584      
    quantization        Q4_K_M    

  Capabilities
    completion    
    tools         

  System
    You are Qwen, created by Alibaba Cloud. You are a helpful assistant.    

  License
    Apache License               
    Version 2.0, January 2004    
```

Run the following command to launch the dev UI.

```sh
adk web
```

Run the following command, to chat with your Google Search agent.

```sh
adk run hello_ollama
```

```sh
给我氯元素的信息

```

`adk api_server` enables you to create a local FastAPI server in a single command, enabling you to test local cURL requests before you deploy your agent.

```sh
adk api_server
```
