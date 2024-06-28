package llm

import (
	"context"
	"fmt"

	"github.com/robinmin/askllm/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	model   string
	llm     *ollama.LLM
	context context.Context
}

func NewOllama(model string, cfg config.LLMEngineConfig) (*Ollama, error) {
	ctx := context.Background()
	var err error
	var llm *ollama.LLM
	if cfg.BaseURL != "" {
		llm, err = ollama.New(ollama.WithModel(model), ollama.WithServerURL(cfg.BaseURL))
	} else {
		llm, err = ollama.New(ollama.WithModel(model))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Ollama: %w", err)
	}

	return &Ollama{
		model:   model,
		llm:     llm,
		context: ctx,
	}, nil
}

func (o *Ollama) Query(prompt string) (string, error) {
	result, err := llms.GenerateFromSinglePrompt(
		o.context, o.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(o.model),
	)
	if err != nil {
		return "", fmt.Errorf("Ollama query failed: %w", err)
	}
	return result, nil
}
