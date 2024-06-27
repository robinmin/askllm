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
		out, err := glamour.Render(content, "dark")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)
		return nil
	}
	return os.WriteFile(outputFile, []byte(content), 0644)
}
