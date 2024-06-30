package llm

import (
	"context"
	"fmt"

	// "os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/robinmin/askllm/internal/config"
)

type ChatGPT struct {
	model   string
	llm     llms.Model
	context context.Context
}

func NewChatGPT(model string, cfg config.LLMEngineConfig) (*ChatGPT, error) {
	var llm llms.Model
	var err error

	ctx := context.Background()
	if model == "" {
		model = cfg.Model
	}
	if cfg.OrgnizationId != "" {
		if cfg.BaseURL != "" {
			llm, err = openai.New(
				openai.WithToken(cfg.APIKey),
				openai.WithModel(model),
				openai.WithOrganization(cfg.OrgnizationId),
				openai.WithBaseURL(cfg.BaseURL),
			)
		} else {
			llm, err = openai.New(
				openai.WithToken(cfg.APIKey),
				openai.WithModel(model),
				openai.WithOrganization(cfg.OrgnizationId),
			)
		}
	} else {
		if cfg.BaseURL != "" {
			llm, err = openai.New(openai.WithToken(cfg.APIKey), openai.WithModel(model), openai.WithBaseURL(cfg.BaseURL))
		} else {
			llm, err = openai.New(openai.WithToken(cfg.APIKey), openai.WithModel(model))
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ChatGPT: %w", err)
	}

	return &ChatGPT{
		model:   model,
		llm:     llm,
		context: ctx,
	}, nil
}

func (c *ChatGPT) Query(prompt string) (string, error) {
	result, err := llms.GenerateFromSinglePrompt(
		c.context, c.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(c.model),
	)
	if err != nil {
		return "", fmt.Errorf("ChatGPT query failed: %w", err)
	}
	return result, nil
}
