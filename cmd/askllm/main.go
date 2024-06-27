package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/internal/llm"
	"github.com/robinmin/askllm/internal/output"
	"github.com/robinmin/askllm/internal/prompt"
)

func main() {
	// Define command-line flags
	engine := flag.String("e", "ollama", "LLM engine (chatgpt, gemini, ollama)")
	model := flag.String("m", "gemma2", "Model for the LLM engine")
	configFile := flag.String("c", "~/.askllm/config.yaml", "Configuration file")
	promptFile := flag.String("p", "", "Prompt file")
	outputFile := flag.String("o", "", "Output file")

	// Parse command-line flags
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize LLM engine
	llmEngine, err := llm.NewEngine(*engine, *model, cfg)
	if err != nil {
		log.Fatalf("Error initializing LLM engine: %v", err)
	}

	// Get prompt
	promptText, err := prompt.GetPrompt(*promptFile, flag.Args())
	if err != nil {
		log.Fatalf("Error getting prompt: %v", err)
	}

	// Query LLM
	response, err := llmEngine.Query(promptText)
	if err != nil {
		log.Fatalf("Error querying LLM: %v", err)
	}

	// Handle output
	if err := output.HandleOutput(*outputFile, response); err != nil {
		log.Fatalf("Error handling output: %v", err)
	}

	fmt.Println("LLM query completed successfully.")
}
