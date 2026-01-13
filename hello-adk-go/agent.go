package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type elementRequest struct {
	Symbol string `json:"symbol"`
}

func getElementInfo(ctx tool.Context, args elementRequest) (map[string]string, error) {
	query := strings.TrimSpace(args.Symbol)
	if query == "" {
		return map[string]string{
			"status":        "error",
			"error_message": "Element symbol or Chinese name is required.",
		}, nil
	}

	element, ok := lookupElement(query)
	if !ok {
		return map[string]string{
			"status":        "error",
			"error_message": fmt.Sprintf("Element symbol or Chinese name '%s' not found.", query),
		}, nil
	}

	report := fmt.Sprintf("%s(%s), Atomic Number: %d, Atomic Weight: %.4f", element.ChineseName, element.Name, element.AtomicNumber, element.AtomicWeight)
	return map[string]string{
		"status": "success",
		"report": report,
	}, nil
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	elementTool, err := functiontool.New(functiontool.Config{
		Name:        "get_element_info",
		Description: "Get information about a chemical element by symbol or Chinese name.",
	}, getElementInfo)
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	periodicAgent, err := llmagent.New(llmagent.Config{
		Name:        "periodic_table_agent",
		Model:       model,
		Description: "Agent to answer questions about chemical element information.",
		Instruction: "You are a helpful agent that provides information about chemical elements. Use the 'get_element_info' tool to look up element information by symbol or Chinese name.",
		Tools:       []tool.Tool{elementTool},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	cfg := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(periodicAgent),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, cfg, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
