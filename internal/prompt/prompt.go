package prompt

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	// "fmt"
	"os"
	"text/template"

	"github.com/robinmin/askllm/pkg/utils"
	"github.com/robinmin/askllm/pkg/utils/log"
)

// PromptTemplate: This struct represents the overall configuration of the prompt template
type PromptTemplate struct {
	Id          string `yaml:"id"`          // Unique identifier for the template
	Name        string `yaml:"name"`        // Name of the personality analyzer template
	Description string `yaml:"description"` // Description of the template's functionality
	Author      string `yaml:"author"`      // Name of the template's author
	Variables   []struct {
		Name       string `yaml:"name"`       // Name of the variable
		Vtype      string `yaml:"vtype"`      // Variable type
		Otype      string `yaml:"otype"`      // Output type of the variable
		Default    string `yaml:"default"`    // Default value for the variable
		Validation string `yaml:"validation"` // Regular expression for validation
	} `yaml:"variables"` // List of variables used by the template
	Template string `yaml:"template"` //  The template string to be used for analysis
}

func GetPrompt(promptFile string, input string) (string, error) {
	log.Info(input)
	if promptFile != "" {
		return readPromptFile(promptFile)
	}
	return input, nil
}

func readPromptFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewPromptTemplate(promptFile string) (*PromptTemplate, error) {
	result, err := utils.LoadConfig[PromptTemplate](promptFile)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// extract default values as a hash map
func (pt *PromptTemplate) getDefaultVars() (map[string]any, error) {
	defaults := make(map[string]any)
	for _, variable := range pt.Variables {
		defaults[variable.Name] = variable.Default
	}
	return defaults, nil
}

// render the template
func (pt *PromptTemplate) GetPrompt(vars map[string]any) (string, error) {
	// get default values
	defaults, err := pt.getDefaultVars()
	if err != nil {
		return "", err
	}

	// merge with inputs
	for key, val := range vars {
		// log.Info(fmt.Sprintf("key = %s : value = %s", key, val))
		defaults[key] = val
	}

	// replace value for all vtype=file/url with content if any
	for _, v := range pt.Variables {
		if strings.ToLower(v.Vtype) == "file" {
			filePath, ok := defaults[v.Name].(string)
			if ok && isValidFilePath(filePath) {
				// Replace the variable with the file content
				fileContent, err := os.ReadFile(filePath)
				if err == nil {
					defaults[v.Name] = string(fileContent)
				}
				// do nothing if the file is not exists
			}
		} else if strings.ToLower(v.Vtype) == "url" {
			value, ok := defaults[v.Name].(string)
			if ok && len(value) > 0 {
				// Load web content from the URL
				resp, err := http.Get(value)
				if err != nil {
					log.Errorf("Failed to fetch URL %s: %v", value, err)
					continue
				}
				// defer resp.Body.Close()
				defer func() {
					if err := resp.Body.Close(); err != nil {
						log.Errorf("Failed to close http resp.Body on %s: %v", value, err)
					}
				}()

				// Read the response body
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					// do nothing if fialed to load web content
					log.Errorf("Failed to read response body from URL %s: %v", value, err)
					continue
				} else {
					// Replace the variable with the web content
					defaults[v.Name] = string(body)
				}
			}
		}
	}

	// render the prompt template
	tmpl, err := template.New("").Parse(pt.Template)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, defaults)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func isValidFilePath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false // Path does not exist
		}
		// Handle other potential errors (e.g., permissions)
		return false
	}
	return true // Path exists (may not be a regular file)
}
