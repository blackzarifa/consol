package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blackzarifa/consol/parser"
	"github.com/blackzarifa/consol/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: consol <file>")
		return
	}

	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	content := string(file)

	if !parser.HasConflict(content) {
		fmt.Println("No conflicts found.")
		return
	}

	conflicts, normalized, lineEnding := parser.ParseFile(content)

	p := tea.NewProgram(tui.InitialModel(normalized, lineEnding, conflicts))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
