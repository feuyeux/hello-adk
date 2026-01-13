# Hello Google Agent Development Kit

1. [Installing ADK](https://google.github.io/adk-docs/get-started/installation/)
2. [Quickstart](https://google.github.io/adk-docs/get-started/quickstart/)

## Install

Create the virtual environment (all OS):

```sh
python -m venv .venv
```

### Mac / Linux

```sh
source .venv/bin/activate
python -m pip install --upgrade pip
python -m pip install google-adk litellm

# Verify (optional)
python -m pip show google-adk litellm | grep -E '^Name|^Version'
```

### Windows PowerShell

```powershell
# If script execution is blocked, run this once per session:
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

. .\.venv\Scripts\Activate.ps1
python -m pip install --upgrade pip
python -m pip install google-adk litellm

# Verify (optional)
python -m pip show google-adk litellm | grep -E '^Name|^Version'
```

### Windows CMD

```bat
.venv\Scripts\activate.bat
python -m pip install --upgrade pip
python -m pip install google-adk litellm

REM Verify (optional)
python -m pip show google-adk litellm
```

### Windows Git Bash

```sh
source .venv/Scripts/activate
python -m pip install --upgrade pip
python -m pip install google-adk litellm

# Verify (optional)
python -m pip show google-adk litellm | grep -E '^Name|^Version'
```

```sh
ollama show qwen2.5
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
[user]: 给我氯元素的信息
```

`adk api_server` enables you to create a local FastAPI server in a single command, enabling you to test local cURL requests before you deploy your agent.

```sh
adk api_server
```
