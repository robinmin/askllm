package main

import (
	"flag"
	"fmt"
	"strings"

	// "os"

	"time"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/internal/llm"
	"github.com/robinmin/askllm/internal/output"
	"github.com/robinmin/askllm/internal/prompt"
	"github.com/robinmin/askllm/pkg/utils/log"
)

func main() {
	// Define command-line flags
	engine := flag.String("e", "ollama", "LLM engine (chatgpt, gemini, ollama)")
	model := flag.String("m", "", "Model for the LLM engine")
	configFile := flag.String("c", "~/.askllm/config.yaml", "Configuration file")
	promptFile := flag.String("p", "", "Prompt file")
	outputFile := flag.String("o", "", "Output file")
	verbose := flag.Bool("v", false, "verbose output")

	// Parse command-line flags
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Error("Error loading configuration: " + err.Error())
		return
	}

	logFile := log.InitLogger(cfg.Sys.LogPath, "askllm", cfg.Sys.LogLevel, *verbose)
	defer log.CloseLogger(logFile)

	if *verbose {
		log.Debug(*configFile)
	}

	log.Info("Starting askllm...(engine: " + *engine + ", model: " + *model + " @ " + config.VERSION + ")")

	// Initialize LLM engine
	llmEngine, err := llm.NewEngine(*engine, *model, cfg)
	if err != nil {
		log.Error("Error initializing LLM engine: " + err.Error())
		return
	}

	// Get prompt
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Error getting current working directory:", err)
	// 	return
	// }
	// log.Infof("currentDir = %v", currentDir)

	// load prompt from external file (compatible with old version)
	payload := strings.Join(flag.Args(), " ")
	promptText, err := prompt.GeneratePrompt(*promptFile, payload)
	if err != nil {
		log.Error("Error getting prompt: " + err.Error())
		return
	}

	// Query LLM
	startTime := time.Now()
	response, err := llmEngine.Query(promptText)
	if err != nil {
		log.Error("Error querying LLM: " + err.Error())
		return
	}
	elapsedTime := time.Since(startTime)

	// Handle output
	if err := output.HandleOutput(*outputFile, response); err != nil {
		log.Error("Error handling output: " + err.Error())
		return
	}

	log.Info(fmt.Sprintf("============== DONE ==============(%s)", elapsedTime))
}
