package llm

import (
	// "bytes"
	"context"
	// "encoding/json"
	"fmt"
	// "io"
	// "net/http"

	"github.com/robinmin/askllm/internal/config"
	"github.com/robinmin/askllm/pkg/utils"
)

type Groq struct {
	model         string
	context       context.Context
	apiKey        string
	orgnizationId string
	chatURL       string
	modelURL      string
	models        []string // List of all available models
}

type chatCompletionRequest struct {
	Messages []message `json:"messages"`
	Model    string    `json:"model"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewGroq(model string, cfg config.LLMEngineConfig) (*Groq, error) {
	ctx := context.Background()
	if model == "" {
		model = cfg.Model
	}

	return &Groq{
		model:         model,
		context:       ctx,
		apiKey:        cfg.APIKey,
		orgnizationId: cfg.OrgnizationId,
		chatURL:       cfg.BaseURL + "/chat/completions",
		modelURL:      cfg.BaseURL + "/models",
	}, nil
}

func (g *Groq) Query(prompt string) (string, error) {
	reqBody := chatCompletionRequest{
		Messages: []message{
			{Role: "user", Content: prompt},
		},
		Model: g.model,
	}

	headers := map[string]string{
		"Authorization": "Bearer " + g.apiKey,
		"Content-Type":  "application/json",
	}

	chatResp, err := utils.APIPost[chatCompletionRequest, chatCompletionResponse](g.chatURL, reqBody, headers)
	if err != nil {
		return "", fmt.Errorf("error fetching models: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}

type GroqModel struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	Created       int64  `json:"created"`
	OwnedBy       string `json:"owned_by"`
	Active        bool   `json:"active"`
	ContextWindow int    `json:"context_window"`
}

type GroqModelListResponse struct {
	Object string      `json:"object"`
	Data   []GroqModel `json:"data"`
}

func (g *Groq) ListAllModelsCore() ([]GroqModel, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + g.apiKey,
		"Content-Type":  "application/json",
	}

	response, err := utils.APIGet[GroqModelListResponse](g.modelURL, headers)
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %v", err)
	}

	if response == nil {
		return nil, fmt.Errorf("received nil response")
	}

	return response.Data, nil
}

func (g *Groq) ListAllModels() ([]string, error) {
	if len(g.models) > 0 {
		return g.models, nil
	}
	models, err := g.ListAllModelsCore()
	if err != nil {
		return nil, err
	}

	g.models = []string{}
	for _, model := range models {
		g.models = append(g.models, model.ID)
	}
	return g.models, nil
}
