@echo off
cd /d %~dp0
mvn compile exec:java -Dexec.mainClass=com.example.AgentCliRunner
