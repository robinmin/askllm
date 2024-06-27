package llm

import (
	"context"
	"fmt"

	"github.com/robinmin/askllm/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

type Gemini struct {
	model   string
	llm     llms.Model
	context context.Context
}

func NewGemini(model string, cfg config.LLMEngineConfig) (*Gemini, error) {
	ctx := context.Background()
	llm, err := googleai.New(ctx, googleai.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini: %w", err)
	}

	return &Gemini{
		model:   model,
		llm:     llm,
		context: ctx,
	}, nil
}

func (g *Gemini) Query(prompt string) (string, error) {
	// result, err := g.llm.Call(g.context, prompt, llms.WithModel(g.model))
	result, err := llms.GenerateFromSinglePrompt(
		g.context, g.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(g.model),
	)
	if err != nil {
		return "", fmt.Errorf("Gemini query failed: %w", err)
	}
	return result, nil
}
