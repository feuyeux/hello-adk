# Hello Google Agent Development Kit - Go

This is a Go implementation of the Periodic Table Agent using Google's Agent Development Kit (ADK).

## Prerequisites

- Go 1.24+ (as recommended by the ADK Go quickstart)
  - Windows: Download and install Go 1.24+ from https://golang.org/dl/
- Ollama running locally (for the qwen2.5 model)

## Setup

1. Install dependencies (module files will be populated by `go mod tidy`):

```bash
cd hello-adk-go
go mod tidy
```


## Run

Interactive CLI:

```bash
go run .
```


```sh
User -> 氢
Agent -> 氢, Atomic Number: 1, Atomic Weight: 1.0080

User -> O
Agent -> 氧(O), Atomic Number: 8, Atomic Weight: 15.9994
```

ADK Web interface:

```bash
go run . web -port 8000 api webui -api_server_address http://localhost:8000/api
```

Then open <http://localhost:8080> and select the agent.
