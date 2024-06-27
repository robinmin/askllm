package prompt

import (
	"os"
	"strings"
)

func GetPrompt(promptFile string, args []string) (string, error) {
	if promptFile != "" {
		return readPromptFile(promptFile)
	}
	return strings.Join(args, " "), nil
}

func readPromptFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
