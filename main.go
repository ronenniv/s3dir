package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ronenniv/s3dir/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		log.Println("log enabled")
		defer f.Close()
	}
	p := tea.NewProgram(model.InitialModel("rniv3030-bucket"), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
