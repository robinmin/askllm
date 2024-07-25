package llm

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	// "time"

	"github.com/PuerkitoBio/goquery"
	"github.com/robinmin/askllm/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/robinmin/askllm/pkg/utils"
)

type Ollama struct {
	model    string
	llm      *ollama.LLM
	context  context.Context
	chatURL  string
	modelURL string
	models   []string // List of all available models
}

func NewOllama(model string, cfg config.LLMEngineConfig) (*Ollama, error) {
	var err error
	var llm *ollama.LLM

	ctx := context.Background()
	if model == "" {
		model = cfg.Model
	}
	if cfg.BaseURL != "" {
		llm, err = ollama.New(ollama.WithModel(model), ollama.WithServerURL(cfg.BaseURL))
	} else {
		llm, err = ollama.New(ollama.WithModel(model))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Ollama: %v", err)
	}

	return &Ollama{
		model:    model,
		llm:      llm,
		context:  ctx,
		chatURL:  cfg.BaseURL + "/chat/completions",
		modelURL: cfg.ExtraURL,
	}, nil
}

func (o *Ollama) Query(prompt string) (string, error) {
	result, err := llms.GenerateFromSinglePrompt(
		o.context, o.llm, prompt,
		llms.WithTemperature(0.2),
		llms.WithModel(o.model),
	)
	if err != nil {
		return "", fmt.Errorf("Ollama query failed: %v", err)
	}
	return result, nil
}

// OllamaModel represents information about an Ollama model.
type OllamaModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pulls       string `json:"pulls"`
	Tags        string `json:"tags"`
	Updated     string `json:"updated"`
}

func (o *Ollama) ListAllModelsCore() ([]OllamaModel, error) {
	// Fetch the webpage content
	// resp, err := http.Get(o.modelURL)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch webpage: %v", err)
	// }
	// defer func() {
	// 	_ = resp.Body.Close()
	// }()

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }
	respBody, err := utils.APIRequestCore(http.MethodGet, o.modelURL, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch webpage: %v", err)
	}
	resp := strings.NewReader(string(respBody))

	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var models []OllamaModel

	// Iterate through each model card
	doc.Find(".grid > li:nth-child(n)").Each(func(i int, s *goquery.Selection) {
		// Extract model information
		name := strings.TrimSpace(s.Find("h2").Text())
		description := strings.TrimSpace(s.Find("p[class='max-w-md break-words']").Text())

		// Use regex to extract pulls and tags from the footer text
		footerText := strings.TrimSpace(s.Find("p[class='my-2 flex space-x-5 text-[13px] font-medium text-neutral-500']").Text())
		footerText = strings.ReplaceAll(footerText, "\u00a0", " ")
		footerText = strings.ReplaceAll(footerText, "\n", " ")
		re := regexp.MustCompile(`(\S+) +Pulls +(\S+) +Tags +Updated +(\S+)`)
		matches := re.FindStringSubmatch(footerText)

		var pulls, tags, updated string
		if len(matches) == 4 {
			pulls = matches[1]
			tags = matches[2]
			updated = matches[3]
		} else {
			pulls = ""
			tags = ""
			updated = ""
		}

		// Create an OllamaModel and append it to the slice
		if len(name) > 0 {
			models = append(models, OllamaModel{
				Name:        name,
				Description: description,
				Pulls:       pulls,
				Tags:        tags,
				Updated:     updated,
			})
		}
	})

	return models, nil
}

func (o *Ollama) ListAllModels() ([]string, error) {
	if len(o.models) > 0 {
		return o.models, nil
	}

	models, err := o.ListAllModelsCore()
	if err != nil {
		return nil, err
	}
	o.models = []string{}
	for _, model := range models {
		o.models = append(o.models, model.Name)
	}
	return o.models, nil
}
