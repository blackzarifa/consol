package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blackzarifa/consol/parser"
	"github.com/blackzarifa/consol/tui"
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

	tui.RunProgram(normalized, lineEnding, conflicts)
}
