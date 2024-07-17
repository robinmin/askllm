package prompt_test

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	testee "github.com/robinmin/askllm/internal/prompt"
	"github.com/robinmin/askllm/pkg/utils"
)

func TestNewPromptTemplate(t *testing.T) {
	t.Run("ValidPromptFile", func(t *testing.T) {
		// Create a temporary file
		data := []byte(generateSamplePrompt())
		tmpFile, err := utils.WriteTempFile("prompt_web_content_extractor", "yaml", data)
		assert.NoError(t, err)

		// Defer cleanup to ensure it happens even if the function returns early
		defer func() {
			err := utils.CleanupTempFile(tmpFile)
			assert.NoError(t, err)
		}()

		// Assuming utils.LoadConfig is mocked to return a valid PromptTemplate
		pt, err := testee.NewPromptTemplate(tmpFile)
		assert.NoError(t, err)
		assert.NotNil(t, pt)
	})

	t.Run("InvalidPromptFile", func(t *testing.T) {
		// Assuming utils.LoadConfig is mocked to return an error
		pt, err := testee.NewPromptTemplate("invalid_prompt.yaml")
		assert.Error(t, err)
		assert.Nil(t, pt)
	})
}

func TestGeneratePrompt(t *testing.T) {
	t.Run("ValidYAMLFile", func(t *testing.T) {
		// Create a temporary file
		data := []byte(generateSamplePrompt())
		tmpFile, err := utils.WriteTempFile("prompt_web_content_extractor", "yaml", data)
		assert.NoError(t, err)

		// Defer cleanup to ensure it happens even if the function returns early
		defer func() {
			err := utils.CleanupTempFile(tmpFile)
			assert.NoError(t, err)
		}()

		// Assuming NewPromptTemplate and GetPrompt are mocked to return valid results
		prompt, err := testee.GeneratePrompt(tmpFile, "content=abc&url_content=123")
		assert.NoError(t, err)
		assert.NotEmpty(t, prompt)
		// fmt.Println(prompt)
	})

	t.Run("InvalidYAMLFile", func(t *testing.T) {
		// Assuming NewPromptTemplate returns an error
		prompt, err := testee.GeneratePrompt("invalid_prompt.yaml", "key1=value1&key2=value2")
		assert.Error(t, err)
		assert.Empty(t, prompt)
	})

	t.Run("ValidPlainText", func(t *testing.T) {
		content := "This is a plain text prompt"
		prompt, err := testee.GeneratePrompt("", content)
		assert.NoError(t, err)
		assert.Equal(t, content, prompt)
	})
}

func generateSamplePrompt() string {
	return `
id: prompt_web_content_extractor
name: "Prompt to fetch the web content and summarize it"
description: "A tiny tool to fetch the web content and summarize it."
author: "Robin Min"
variables:
  - name: "content"
    vtype: "string"
    otype: "text"
    default: ""
    validation: ""
  - name: "file_content"
    vtype: "file"
    otype: "text"
    default: ""
    validation: ""
  - name: "url_content"
    vtype: "url"
    otype: "text"
    default: ""
    validation: ""
template: |
  Here comes the major content of a web page. I need you provide a concise summary in original language and then translate it into Chinese.
  You also need to list out the major points mentioned by the web page. Both the summary and key points will be listed in oroginal
  langiage and Chinese. Here comes the web page content:
  
  {{if .content }}{{ .content }}{{end}}{{if .file_content }}{{ .file_content }}{{end}}{{if .url_content }}{{ .url_content }}{{end}}
  
  Show me the summary and key points in both original language and Chinese.
  `
}
