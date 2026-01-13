# Hello Google Agent Development Kit - Go

This is a Go implementation of the Periodic Table Agent using Google's Agent Development Kit (ADK).

## Prerequisites

- Go 1.24+ (as recommended by the ADK Go quickstart)

## Setup

1. Install dependencies (module files will be populated by `go mod tidy`):

```bash
cd hello-adk-go
go mod tidy
```

1. Set your API key:

```bash
echo 'export GOOGLE_API_KEY="YOUR_API_KEY"' > .env
# Windows PowerShell
Set-Content .env 'set GOOGLE_API_KEY="YOUR_API_KEY"'
```

## Run

Interactive CLI:

```bash
# remember to load keys: source .env OR run the PowerShell command above
go run agent.go
```

ADK Web interface:

```bash
# remember to load keys
go run agent.go web api webui
```

Then open <http://localhost:8080> and select the agent.

## Agent Features

- Look up chemical elements by symbol (e.g., "H", "O", "Fe")
- Look up chemical elements by Chinese name (e.g., "氢", "氧", "铁")
- Returns: Chinese name, English name, atomic number, atomic weight

## Notes

- The tool is implemented with `functiontool.New` and exposes `get_element_info`.
- Model: `gemini-2.5-flash` (change in `agent.go` if needed).

## Project Structure

```
hello-adk-go/
├── agent.go             # Agent definition and launcher
├── periodic_table.go    # Periodic table data and lookup
├── go.mod               # Module definition
├── .env                 # API key placeholder
└── .gitignore
```
