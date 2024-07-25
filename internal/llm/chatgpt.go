package llm

import (
	"context"
	"fmt"

	// "os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/pkg/utils"
)

type ChatGPT struct {
	model    string
	llm      llms.Model
	context  context.Context
	chatURL  string
	modelURL string
	models   []string // List of all available models
	apiKey   string   // API token
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
		return nil, fmt.Errorf("failed to initialize ChatGPT: %v", err)
	}

	return &ChatGPT{
		model:    model,
		llm:      llm,
		context:  ctx,
		chatURL:  cfg.BaseURL + "/chat/completions",
		modelURL: cfg.BaseURL + "/models",
		apiKey:   cfg.APIKey,
	}, nil
}

func (c *ChatGPT) Query(prompt string) (string, error) {
	result, err := llms.GenerateFromSinglePrompt(
		c.context, c.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(c.model),
	)
	if err != nil {
		return "", fmt.Errorf("ChatGPT query failed: %v", err)
	}
	return result, nil
}

type ObjectType string

const (
	OTModel           ObjectType = "model"
	OTModelPermission ObjectType = "model_permission"
	OTList            ObjectType = "list"
	OTEdit            ObjectType = "edit"
	OTTextCompletion  ObjectType = "text_completion"
	OTEEmbedding      ObjectType = "embedding"
	OTFile            ObjectType = "file"
	OTFineTune        ObjectType = "fine-tune"
	OTFineTuneEvent   ObjectType = "fine-tune-event"
)

type ChatGPTModel struct {
	ID         string            `json:"id"`
	Object     ObjectType        `json:"object"`
	Created    int64             `json:"created"`
	OwnedBy    string            `json:"owned_by"`
	Permission []ModelPermission `json:"permission"`
	Root       string            `json:"root"`
	Parent     string            `json:"parent"`
}

type ModelPermission struct {
	ID                 string     `json:"id"`
	Object             ObjectType `json:"object"`
	Created            int64      `json:"created"`
	AllowCreateEngine  bool       `json:"allow_create_engine"`
	AllowSampling      bool       `json:"allow_sampling"`
	AllowLogProbs      bool       `json:"allow_logprobs"`
	AllowSearchIndices bool       `json:"allow_search_indices"`
	AllowView          bool       `json:"allow_view"`
	AllowFineTuning    bool       `json:"allow_fine_tuning"`
	Organization       string     `json:"organization"`
	Group              string     `json:"group"`
	IsBlocking         bool       `json:"is_blocking"`
}

type ChatGPTModelsListResponse struct {
	Data   []ChatGPTModel `json:"data"`
	Object ObjectType
}

func (c *ChatGPT) ListAllModelsCore() ([]ChatGPTModel, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + c.apiKey,
		"Content-Type":  "application/json",
	}

	response, err := utils.APIGet[ChatGPTModelsListResponse](c.modelURL, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching ChatGPT models: %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("no ChatGPT models found")
	}

	// Filter for only Gemini models
	// var models []ChatGPTModel
	// for _, model := range response.Data {
	// 	models = append(models, model.ID)
	// }

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no ChatGPT models found in the response")
	}

	return response.Data, nil
}

func (c *ChatGPT) ListAllModels() ([]string, error) {
	if len(c.models) > 0 {
		return c.models, nil
	}

	// Fetch all models
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
