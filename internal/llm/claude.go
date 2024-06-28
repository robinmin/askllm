package llm

import (
	"context"
	"fmt"

	"github.com/robinmin/askllm/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
)

type Claude struct {
	model   string
	llm     llms.Model
	context context.Context
}

func NewClaude(model string, cfg config.LLMEngineConfig) (*Claude, error) {
	ctx := context.Background()
	var llm llms.Model
	var err error
	if cfg.BaseURL != "" {
		llm, err = anthropic.New(anthropic.WithToken(cfg.APIKey), anthropic.WithModel(model), anthropic.WithBaseURL(cfg.BaseURL))
	} else {
		llm, err = anthropic.New(anthropic.WithToken(cfg.APIKey), anthropic.WithModel(model))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Claude: %w", err)
	}

	return &Claude{
		model:   model,
		llm:     llm,
		context: ctx,
	}, nil
}

func (c *Claude) Query(prompt string) (string, error) {
	result, err := llms.GenerateFromSinglePrompt(
		c.context, c.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(c.model),
	)
	if err != nil {
		return "", fmt.Errorf("Claude query failed: %w", err)
	}
	return result, nil
}
