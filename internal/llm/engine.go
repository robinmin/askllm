package llm

import (
	"fmt"
	"strings"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/pkg/utils/log"
)

type Engine interface {
	Query(prompt string) (string, error)
}

func NewEngine(engineType, model string, cfg *config.Config) (Engine, error) {
	// use provided engine type first
	tmpEngine := strings.TrimSpace(strings.ToLower(engineType))

	if tmpEngine == "" {
		tmpEngine = "ollama" // if still no engine type is provided, use ollama
	}

	enginCfg, ok := cfg.LLMEngines[tmpEngine]
	if !ok {
		tmpEngine = "ollama" // if still no engine type is provided, use ollama
		enginCfg = cfg.LLMEngines[tmpEngine]
	}

	// use provided model
	tmpModel := strings.TrimSpace(strings.ToLower(model))
	if tmpModel == "" {
		tmpModel = enginCfg.Model // if no model is not provided, use the default model
	}

	log.Infof("Using LLM engine: %s, model: %s", tmpEngine, tmpModel)

	switch tmpEngine {
	case "chatgpt":
		return NewChatGPT(tmpModel, enginCfg)
	case "gemini":
		return NewGemini(tmpModel, enginCfg)
	case "ollama":
		return NewOllama(tmpModel, enginCfg)
	case "claude":
		return NewClaude(tmpModel, enginCfg)
	case "groq":
		return NewGroq(tmpModel, enginCfg)
	default:
		return nil, fmt.Errorf("unsupported LLM engine: %s", tmpEngine)
	}
}
