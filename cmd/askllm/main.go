package main

import (
	"flag"
	"fmt"
	"net/url"
	// "os"
	"regexp"
	"strings"
	"time"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/internal/llm"
	"github.com/robinmin/askllm/internal/output"
	"github.com/robinmin/askllm/internal/prompt"
	"github.com/robinmin/askllm/pkg/utils/log"
)

func isQueryString(s string) bool {
	// Check for presence of equals sign (=) and separator (&)
	// return len(s) > 0 && strings.ContainsRune(s, '=') && strings.ContainsRune(s, '&')
	// Regex pattern for key-value pairs with optional '&' separators
	pattern := `^([a-zA-Z0-9_*]+)(?:\s*=\s*|)=([^&]*)?(&([a-zA-Z0-9_*]+)(?:\s*=\s*|)=([^&]*)*)*$`
	return regexp.MustCompile(pattern).MatchString(s)
}

// ParseQueryString parses a query string and returns a map[string]string or an error.
func ParseQueryString(queryString string) (map[string]any, error) {
	if !isQueryString(queryString) {
		return nil, fmt.Errorf("invalid query string format")
	}

	query, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("error parsing query string: %w", err)
	}

	queryParams := make(map[string]any)
	for key, values := range query {
		queryParams[key] = values[0] // Assuming single value for each key
	}

	return queryParams, nil
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
	var promptText string
	payload := strings.Join(flag.Args(), " ")
	if len(*promptFile) > 0 {
		fileName := strings.ToLower(*promptFile)
		if strings.HasSuffix(fileName, ".yaml") || strings.HasSuffix(fileName, ".yml") {
			// load prompt from prompt template YAML file
			pt, err := prompt.NewPromptTemplate(*promptFile)
			if err != nil {
				log.Error("Failed to create instance of PromptTemplate: " + err.Error())
				return
			}

			var vars map[string]any
			if isQueryString(payload) {
				vars, err = ParseQueryString(payload)
				if err != nil {
					log.Error("Failed to unmarshal variables from query string: " + err.Error())
					return
				}
			}

			promptText, err = pt.GetPrompt(vars)
			if err != nil {
				log.Error("Error getting prompt: " + err.Error())
				return
			}
		} else {
			// load prompt from external file (compatible with old version)
			promptText, err = prompt.GetPrompt(*promptFile, payload)
			if err != nil {
				log.Error("Error getting prompt: " + err.Error())
				return
			}
		}
	} else {
		// load prompt from command line directly
		promptText = payload
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
