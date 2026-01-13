# Hello Google Agent Development Kit - Java

This is a Java implementation of the Periodic Table Agent using Google's Agent Development Kit (ADK).

**Note**: The Java ADK (v0.5.0) currently supports Gemini models natively. For Ollama support, you would need to:

1. Use the Python ADK which has LiteLLM integration
2. Wait for official Ollama/LangChain4j integration in future ADK versions
3. Set up an OpenAI-compatible proxy for Ollama (e.g., using litellm proxy)

## Installation

1. Install dependencies:

```bash
mvn clean install
```

## Running the Agent

### Command-line Interface

Run the agent with the interactive CLI:

**Windows:**

```cmd
run.bat
```

**Linux/Mac or direct Maven command:**

```bash
mvn compile exec:java -Dexec.mainClass="com.example.AgentCliRunner"
```

Optional environment variables:

- `MODEL_NAME`: Model to use (default: `gemini-2.0-flash-exp`)
- `GOOGLE_AI_API_KEY`: Your Google AI API key (required)

### Web Interface

Run the agent with the ADK web interface:

```bash
mvn compile exec:java -Dexec.mainClass="com.google.adk.web.AdkWebServer" -Dexec.args="--adk.agents.source-dir=target --server.port=8000"
```

Then access the web interface at <http://localhost:8000>

## Example Usage

```sh
You > Tell me about hydrogen
Agent > Hydrogen (氢), Atomic Number: 1, Atomic Weight: 1.0080

You > What is the atomic weight of 金 (gold)?
Agent > 金(Gold), Atomic Number: 79, Atomic Weight: 196.9700

You > quit
```

## Project Structure

```
hello-adk-java/
├── src/main/java/com/example/
│   ├── PeriodicTableAgent.java   # Main agent definition with Ollama integration
│   ├── PeriodicTable.java        # Periodic table data and element lookup
│   └── AgentCliRunner.java       # CLI runner for the agent
└── pom.xml                       # Maven configuration
```

## Architecture

The agent uses:

- **Google Gemini**: LLM model via Google AI API
- **Google ADK**: Agent framework with tool calling support  
- **LangChain4j**: (Added but not yet integrated - future Ollama support)

**Current Limitation**: Java ADK 0.5.0 only supports Gemini models. The Python ADK has broader model support through LiteLLM.

## References

- [ADK Java Quickstart](https://google.github.io/adk-docs/get-started/java/)
- [ADK Documentation](https://google.github.io/adk-docs/)
- [LangChain4j](https://docs.langchain4j.dev/) - For future Ollama integration
