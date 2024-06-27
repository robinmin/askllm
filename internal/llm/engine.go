package llm

import (
	"fmt"
	"strings"

	"github.com/robinmin/askllm/internal/config"
)

type Engine interface {
	Query(prompt string) (string, error)
}

func NewEngine(engineType, model string, cfg *config.Config) (Engine, error) {
	switch strings.ToLower(engineType) {
	case "chatgpt":
		return NewChatGPT(model, cfg.LLMEngines["chatgpt"])
	case "gemini":
		return NewGemini(model, cfg.LLMEngines["gemini"])
	case "ollama":
		return NewOllama(model, cfg.LLMEngines["ollama"])
	default:
		return nil, fmt.Errorf("unsupported LLM engine: %s", engineType)
	}
}
