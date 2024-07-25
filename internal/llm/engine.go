package llm

import (
	"fmt"
	"strings"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/pkg/utils/log"
)

type Engine interface {
	Query(prompt string) (string, error)
	ListAllModels() ([]string, error)
}

func NewEngine(engineType, model string, cfg *config.Config) (Engine, error) {
	// use provided engine type first
	tmpEngine := strings.TrimSpace(strings.ToLower(engineType))
	tmpModel := strings.TrimSpace(strings.ToLower(model))

	if tmpEngine == "" {
		if len(cfg.Sys.DefaultEngine) > 0 {
			tmpEngine = cfg.Sys.DefaultEngine
		} else {
			tmpEngine = "ollama" // if still no engine type is provided, use ollama
		}
	}
	if tmpModel == "" {
		tmpModel = GetDefaultModel(tmpEngine)
	}

	engineCfg, ok := cfg.LLMEngines[tmpEngine]
	if !ok {
		tmpEngine = "ollama" // if still no engine type is provided, use ollama
		engineCfg = cfg.LLMEngines[tmpEngine]
	}

	log.Infof("Using LLM engine: %s, model: %s", tmpEngine, tmpModel)

	switch tmpEngine {
	case "chatgpt":
		return NewChatGPT(tmpModel, engineCfg)
	case "gemini":
		return NewGemini(tmpModel, engineCfg)
	case "ollama":
		return NewOllama(tmpModel, engineCfg)
	case "claude":
		return NewClaude(tmpModel, engineCfg)
	case "groq":
		return NewGroq(tmpModel, engineCfg)
	default:
		return nil, fmt.Errorf("unsupported LLM engine: %s", tmpEngine)
	}
}

func GetDefaultModel(engine string) string {
	switch strings.TrimSpace(strings.ToLower(engine)) {
	case "chatgpt":
		return "gpt-4o-mini"
	case "gemini":
		return "gemini-1.5-pro"
	case "ollama":
		return "gemma2"
	case "claude":
		return "claude-3-sonnet-20240229"
	case "groq":
		return "gemma2-9b-it"
	default:
		return ""
	}
}

func GetAllModels(engineType string, cfg *config.Config) (map[string][]string, error) {
	result := map[string][]string{}
	tmpEngine := strings.TrimSpace(strings.ToLower(engineType))

	if _, ok := cfg.LLMEngines[tmpEngine]; ok {
		// if found, return models for specified engine
		if llmObj, err := NewEngine(engineType, "", cfg); err == nil {
			if tmp, err := llmObj.ListAllModels(); err == nil {
				result[tmpEngine] = tmp
			}
		} else {
			return nil, err
		}
	} else {
		// if not found, return models for all engines
		for tmpEngine := range cfg.LLMEngines {
			if llmObj, err := NewEngine(tmpEngine, "", cfg); err == nil {
				if tmp, err := llmObj.ListAllModels(); err == nil {
					result[tmpEngine] = tmp
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}
	return result, nil
}
