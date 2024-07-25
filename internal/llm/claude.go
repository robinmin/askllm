package llm

import (
	"context"
	"fmt"
	// "strings"

	// "github.com/PuerkitoBio/goquery"
	"github.com/robinmin/askllm/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	// "github.com/robinmin/askllm/pkg/utils"
)

type Claude struct {
	model   string
	llm     llms.Model
	context context.Context
	chatURL string
	// modelURL string
	models []string // List of all available models
}

func NewClaude(model string, cfg config.LLMEngineConfig) (*Claude, error) {
	var llm llms.Model
	var err error

	ctx := context.Background()
	if model == "" {
		model = cfg.Model
	}
	if cfg.BaseURL != "" {
		llm, err = anthropic.New(anthropic.WithToken(cfg.APIKey), anthropic.WithModel(model), anthropic.WithBaseURL(cfg.BaseURL))
	} else {
		llm, err = anthropic.New(anthropic.WithToken(cfg.APIKey), anthropic.WithModel(model))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Claude: %v", err)
	}

	return &Claude{
		model:   model,
		llm:     llm,
		context: ctx,
		chatURL: cfg.BaseURL + "/chat/completions",
		// modelURL: "https://docs.anthropic.com/en/docs/about-claude/models#model-names",
	}, nil
}

func (c *Claude) Query(prompt string) (string, error) {
	result, err := llms.GenerateFromSinglePrompt(
		c.context, c.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(c.model),
	)
	if err != nil {
		return "", fmt.Errorf("Claude query failed: %v", err)
	}
	return result, nil
}

type ClaudeModel struct {
	ID              string
	Description     string
	ContentWindow   int
	MaxOutputTokens int
}

func (c *Claude) ListAllModelsCore() ([]ClaudeModel, error) {
	// headers := map[string]string{
	// 	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	// }

	// responseBody, err := utils.APIRequestCore("GET", c.modelURL, nil, headers)
	// if err != nil {
	// 	return nil, fmt.Errorf("error fetching Anthropic documentation page: %v", err)
	// }

	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(responseBody)))
	// if err != nil {
	// 	return nil, fmt.Errorf("error parsing HTML: %v", err)
	// }

	// var models []ClaudeModel

	// doc.Find("h2:contains('Model names')").NextUntil("h2").Find("li").Each(func(i int, s *goquery.Selection) {
	// 	text := strings.TrimSpace(s.Text())
	// 	parts := strings.SplitN(text, ":", 2)
	// 	if len(parts) == 2 {
	// 		model := ClaudeModel{
	// 			ID:          strings.TrimSpace(parts[0]),
	// 			Description: strings.TrimSpace(parts[1]),
	// 		}
	// 		models = append(models, model)
	// 	}
	// })

	// if len(models) == 0 {
	// 	return nil, fmt.Errorf("no models found on the page")
	// }

	models := []ClaudeModel{
		{
			ID:              "claude-3-5-sonnet-20240620",
			Description:     "Claude 3.5 Sonnet",
			ContentWindow:   200000,
			MaxOutputTokens: 8192,
		},
		// {ID: "claude-3-5-opus", Description: "Claude 3.5 Opus(Later this year)"},   // Pending official release date
		// {ID: "claude-3-5-haiku", Description: "Claude 3.5 Haiku(Later this year)"}, // Pending official release date
		{
			ID:              "claude-3-opus-20240229",
			Description:     "Claude 3 Opus",
			ContentWindow:   200000,
			MaxOutputTokens: 4096,
		}, {
			ID:              "claude-3-sonnet-20240229",
			Description:     "Claude 3 Sonnet",
			ContentWindow:   200000,
			MaxOutputTokens: 4096,
		}, {
			ID:              "claude-3-haiku-20240307",
			Description:     "Claude 3 Haiku",
			ContentWindow:   200000,
			MaxOutputTokens: 4096,
		}, {
			ID:              "claude-2.1",
			Description:     "Claude 2.1(Legacy Model)",
			ContentWindow:   200000,
			MaxOutputTokens: 4096,
		}, {
			ID:              "claude-2.0",
			Description:     "Claude 2(Legacy Model)",
			ContentWindow:   100000,
			MaxOutputTokens: 4096,
		}, {
			ID:              "claude-instant-1.2",
			Description:     "Claude Instant 1.2(Legacy Model)",
			ContentWindow:   100000,
			MaxOutputTokens: 4096,
		},
	}
	return models, nil
}

func (c *Claude) ListAllModels() ([]string, error) {
	if len(c.models) > 0 {
		return c.models, nil
	}

	models, err := c.ListAllModelsCore()
	if err != nil {
		return nil, err
	}
	c.models = []string{}
	for _, model := range models {
		c.models = append(c.models, model.ID)
	}
	return c.models, nil
}
