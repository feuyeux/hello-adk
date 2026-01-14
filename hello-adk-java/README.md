# Hello Google Agent Development Kit - Java

## Installation

1. Install dependencies:

```bash
mvn clean install
```

## Running the Agent

```bash
mvn compile exec:java
```

Or using the full command:

```bash
mvn compile exec:java -Dexec.mainClass="com.example.AgentCliRunner"
```

## Running the Web Interface

```bash
mvn compile exec:java -Dexec.mainClass=com.google.adk.web.AdkWebServer -Dexec.args="--adk.agents.source-dir=src/main/java/com/example --server.port=8080"
```

Then access the web interface at <http://localhost:8080>

## Running the CLI Interface

```bash
mvn compile exec:java
```

Or using the full command:

```bash
mvn compile exec:java -Dexec.mainClass="com.example.AgentCliRunner"
```

## Example Usage

```sh
You > Tell me about hydrogen
```