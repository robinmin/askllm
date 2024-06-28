package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	log "log/slog"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/internal/llm"
	"github.com/robinmin/askllm/internal/output"
	"github.com/robinmin/askllm/internal/prompt"
)

func ParseLevel(s string) (log.Level, error) {
	var level log.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}

func main() {
	// Define command-line flags
	engine := flag.String("e", "ollama", "LLM engine (chatgpt, gemini, ollama)")
	model := flag.String("m", "gemma2", "Model for the LLM engine")
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

	// Set the custom logger as the default
	var logFile *os.File
	lvl, err := ParseLevel(strings.ToUpper(cfg.Sys.LogLevel))
	if err == nil || *verbose {
		lvl = log.LevelDebug
	}
	opts := &log.HandlerOptions{
		Level: lvl,
	}

	if len(cfg.Sys.LogPath) > 0 {
		// Create a custom JSON logger
		filename := fmt.Sprintf("%s/log_%s.log", cfg.Sys.LogPath, time.Now().Format("20060102"))

		var err1 error
		logFile, err1 = os.Create(filename)
		if err1 != nil {
			log.Error((err1.Error())) // Handle errors appropriately
			return
		}
		logger := log.New(log.NewJSONHandler(logFile, opts))
		log.SetDefault(logger)
	} else {
		logger := log.New(log.NewTextHandler(os.Stdout, opts))
		log.SetDefault(logger)
	}
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

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
	promptText, err := prompt.GetPrompt(*promptFile, flag.Args())
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
