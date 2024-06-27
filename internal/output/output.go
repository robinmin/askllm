package output

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/glamour"
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
			log.Fatal(err)
		}
		fmt.Println(out)
		return nil
	}
	return os.WriteFile(outputFile, []byte(content), 0644)
}
