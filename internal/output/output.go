package output

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"

	"github.com/robinmin/askllm/pkg/utils/log"
)

func HandleOutput(outputFile, content string) error {
	if outputFile == "" || outputFile == "stdout" {
		// show markdown in console
		r, _ := glamour.NewTermRenderer(
			// detect background color and pick either the default dark or light theme
			glamour.WithAutoStyle(),
			// wrap output at specific width (default is 80)
			glamour.WithWordWrap(120),
		)
		out, err := r.Render(content)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		fmt.Println(out)
		return nil
	}
	return os.WriteFile(outputFile, []byte(content), 0644)
}

func OutputMarkdown(content string) error {
	// show markdown in console
	r, _ := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width (default is 80)
		glamour.WithWordWrap(120),
	)
	out, err := r.Render(content)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	fmt.Println(out)
	return nil
}
