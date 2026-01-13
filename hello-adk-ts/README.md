# Hello Google Agent Development Kit - TypeScript

This is a TypeScript implementation of the Periodic Table Agent using Google's Agent Development Kit (ADK).

## Installation

### Prerequisites

- Node.js 20.12.7 or later
- npm 9.2.0 or later

### Setup

1. Install dependencies:

```bash
npm install
```

1. Create the `.env` file and set your API key:

```bash
echo 'GEMINI_API_KEY="YOUR_API_KEY"' > .env
```

1. Compile TypeScript (optional, @google/adk-devtools can handle this):

```bash
npx tsc
```

## Running the Agent

### Command-line Interface

Run the agent with the interactive CLI:

```bash
npx @google/adk-devtools run agent.ts
```

### Web Interface

Run the agent with the ADK web interface:

```bash
npx @google/adk-devtools web
```

Then access the web interface at <http://localhost:8000>

## Agent Features

The Periodic Table Agent can:

- Look up chemical elements by symbol (e.g., "H", "O", "Fe")
- Look up chemical elements by Chinese name (e.g., "氢", "氧", "铁")
- Return element information including:
  - Chinese name
  - English name
  - Atomic number
  - Atomic weight

## Example Usage

```
You > Tell me about hydrogen
Agent > Hydrogen (氢), Atomic Number: 1, Atomic Weight: 1.008

You > What is the atomic weight of 金 (gold)?
Agent > 金(Gold), Atomic Number: 79, Atomic Weight: 196.97

You > quit
```

## Project Structure

```
hello-adk-ts/
├── agent.ts              # Main agent definition
├── periodic-table.ts     # Periodic table data and element lookup
├── package.json          # npm configuration
├── tsconfig.json         # TypeScript configuration
└── .env                  # Environment variables (API keys)
```

## References

- [ADK TypeScript Quickstart](https://google.github.io/adk-docs/get-started/typescript/)
- [ADK Documentation](https://google.github.io/adk-docs/)
