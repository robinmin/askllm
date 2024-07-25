package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/pkg/utils"
)

type Gemini struct {
	model    string
	llm      llms.Model
	context  context.Context
	chatURL  string
	modelURL string
	models   []string // List of all available models
	apiKey   string   // API token
}

func NewGemini(model string, cfg config.LLMEngineConfig) (*Gemini, error) {
	ctx := context.Background()
	if model == "" {
		model = cfg.Model
	}
	llm, err := googleai.New(ctx, googleai.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini: %v", err)
	}

	return &Gemini{
		model:    model,
		llm:      llm,
		context:  ctx,
		chatURL:  cfg.BaseURL + "/chat/completions",
		modelURL: cfg.ExtraURL + "/models",
		apiKey:   cfg.ExtraKey,
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
		return "", fmt.Errorf("Gemini query failed: %v", err)
	}
	return result, nil
}

type GeminiModel struct {
	Name                       string   `json:"name"`
	Version                    string   `json:"version"`
	DisplayName                string   `json:"displayName"`
	Description                string   `json:"description"`
	InputTokenLimit            int      `json:"inputTokenLimit"`
	OutputTokenLimit           int      `json:"outputTokenLimit"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	Temperature                float64  `json:"temperature"`
	TopP                       float64  `json:"topP"`
	TopK                       int      `json:"topK"`
}

type GeminiModelListResponse struct {
	Models []GeminiModel `json:"models"`
}

func (g *Gemini) ListAllModelsCore() ([]GeminiModel, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Add API key as a query parameter
	url := fmt.Sprintf("%s?key=%s", g.modelURL, g.apiKey)

	response, err := utils.APIGet[GeminiModelListResponse](url, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching Gemini models: %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("no Gemini models found")
	}

	if len(response.Models) == 0 {
		return nil, fmt.Errorf("no Gemini models found in the response")
	}

	return response.Models, nil
}

func (g *Gemini) ListAllModels() ([]string, error) {
	if len(g.models) > 0 {
		return g.models, nil
	}

	models, err := g.ListAllModelsCore()
	if err != nil {
		return nil, err
	}
	// var names []string
	g.models = []string{}
	for _, model := range models {
		// erase the leads substring "models/" if it exists
		if strings.HasPrefix(model.Name, "models/") {
			g.models = append(g.models, model.Name[len("models/"):])
		} else {
			g.models = append(g.models, model.Name)
		}
	}
	return g.models, nil
}
