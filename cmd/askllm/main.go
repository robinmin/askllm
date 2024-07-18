package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/internal/llm"
	"github.com/robinmin/askllm/internal/output"
	"github.com/robinmin/askllm/internal/prompt"
	"github.com/robinmin/askllm/pkg/utils/log"
)

var (
	// Define command-line flags
	engine     *string
	model      *string
	configFile *string
	promptFile *string
	outputFile *string
	verbose    *bool
)

func init() {
	engine = flag.String("e", "", "LLM engine (chatgpt, gemini, ollama)")
	model = flag.String("m", "", "Model for the LLM engine")
	configFile = flag.String("c", "~/.askllm/config.yaml", "Locatuon of configuration file")
	promptFile = flag.String("p", "", "Prompt file")
	outputFile = flag.String("o", "", "Output file")
	verbose = flag.Bool("v", false, "verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (version %s):\n", os.Args[0], config.VERSION)
		flag.PrintDefaults()
	}
}

func main() {
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
		log.Debug("Loading config file from : " + *configFile)
	}

	// Get pwd
	if *verbose {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}
		log.Debugf("currentDir = %v", currentDir)
	}

	log.Info("Starting askllm...(engine: " + *engine + ", model: " + *model + " @ " + config.VERSION + ")")
	payload := strings.Join(flag.Args(), " ")
	startTime := time.Now()

	// load prompt from external file (compatible with old version)
	pt, promptText, err := prompt.GeneratePrompt(*promptFile, payload)
	if err != nil {
		log.Error("Error getting prompt: " + err.Error())
		return
	}

	// // Initialize LLM engine
	realEngine, realModel := pt.GetParameters(*engine, *model)
	llmEngine, err := llm.NewEngine(realEngine, realModel, cfg)
	if err != nil {
		log.Error("Error initializing LLM engine: " + err.Error())
		return
	}

	// Query LLM
	response, err := llmEngine.Query(promptText)
	if err != nil {
		log.Error("Error querying LLM: " + err.Error())
		return
	}

	// Handle output
	if err := output.HandleOutput(*outputFile, response); err != nil {
		log.Error("Error handling output: " + err.Error())
		return
	}

	elapsedTime := time.Since(startTime)
	log.Info(fmt.Sprintf("============== DONE ==============(%s)", elapsedTime))
}
