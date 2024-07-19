package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/robinmin/askllm/internal/config"
)

type Groq struct {
	model         string
	context       context.Context
	apiKey        string
	orgnizationId string
	baseURL       string
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
		baseURL:       cfg.BaseURL + "/chat/completions",
	}, nil
}

func (g *Groq) Query(prompt string) (string, error) {
	reqBody := chatCompletionRequest{
		Messages: []message{
			{Role: "user", Content: prompt},
		},
		Model: g.model,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(g.context, "POST", g.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var chatResp chatCompletionResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
