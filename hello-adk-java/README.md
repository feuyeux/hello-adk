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

Then access the web interface at <http://localhost:8000>
> Note: The web interface is currently not available in the Java version due to missing dependencies. Please use the CLI runner.

## Example Usage

```sh
You > Tell me about hydrogen
```