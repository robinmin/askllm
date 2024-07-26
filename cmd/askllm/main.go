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
	action     *string
	engine     *string
	model      *string
	configFile *string
	promptFile *string
	outputFile *string
	verbose    *bool
)

func init() {
	action = flag.String("a", "client", "subcommand, so far support 'client', 'server', 'models'")
	engine = flag.String("e", "", "LLM engine (chatgpt, gemini, ollama, claude, groq)")
	model = flag.String("m", "", "Model for the LLM engine")
	configFile = flag.String("c", "~/.askllm/config.yaml", "Locatuon of configuration file")
	promptFile = flag.String("p", "", "Prompt file or prompt text")
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

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}
		log.Debugf("currentDir = %v", currentDir)
	}

	startTime := time.Now()
	log.Info("Starting askllm...(engine: " + *engine + ", model: " + *model + " @ " + config.VERSION + ")")
	payload := strings.Join(flag.Args(), " ")

	switch strings.ToLower(*action) {
	case "client":
		err = runClientAction(*promptFile, payload, *engine, *model, cfg)
	case "server":
		err = runServerAction(*promptFile, payload, *engine, *model, cfg)
	case "models":
		err = runModelsAction(*promptFile, payload, *engine, *model, cfg)
	default:
		log.Error("Invalid action: " + *action)
		err = fmt.Errorf("invalid action: %s", *action)
		flag.Usage()
	}
	if err != nil {
		log.Error("Error: " + err.Error())
	}

	elapsedTime := time.Since(startTime)
	log.Info(fmt.Sprintf("============== DONE ==============(%s)", elapsedTime))
}

func runClientAction(promptFile string, payload string, engine string, model string, cfg *config.Config) error {
	// load prompt from external file (compatible with old version)
	pt, promptText, err := prompt.GeneratePrompt(promptFile, payload)
	if err != nil {
		log.Error("Error getting prompt: " + err.Error())
		return err
	}

	// // Initialize LLM engine
	realEngine, realModel := pt.GetParameters(engine, model, cfg.Sys.DefaultEngine, llm.GetDefaultModel(cfg.Sys.DefaultEngine))
	llmEngine, err := llm.NewEngine(realEngine, realModel, cfg)
	if err != nil {
		log.Error("Error initializing LLM engine: " + err.Error())
		return err
	}

	// Query LLM
	response, err := llmEngine.Query(promptText)
	if err != nil {
		log.Error("Error querying LLM: " + err.Error())
		return err
	}

	// Handle output
	if err := output.HandleOutput(*outputFile, response); err != nil {
		log.Error("Error handling output: " + err.Error())
		return err
	}
	return nil
}

func runServerAction(promptFile string, payload string, engine string, model string, cfg *config.Config) error {
	log.Error("Not implemented")
	return fmt.Errorf("not implemented")
}

func runModelsAction(promptFile string, payload string, engine string, model string, cfg *config.Config) error {
	log.Info("List models for LLM engine: " + engine)

	// if engine == "" {
	// 	engine = cfg.Sys.DefaultEngine
	// }
	allModels, err := llm.GetAllModels(engine, cfg)
	if err != nil {
		log.Error("Error listing models: " + err.Error())
		return err
	}

	var content []string
	for engineName, modelList := range allModels {
		content = append(content, "#### "+engineName+":\n")
		defaultModel := llm.GetDefaultModel(engineName)
		for _, tmpModel := range modelList {
			if tmpModel == defaultModel {
				content = append(content, "- [*] "+tmpModel)
			} else {
				content = append(content, "- "+tmpModel)
			}
		}
		content = append(content, "")
	}
	if err := output.OutputMarkdown(strings.Join(content, "\n")); err != nil {
		log.Error("Error in output markdown : " + err.Error())
		return err
	}
	return nil
}
