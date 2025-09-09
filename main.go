package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blackzarifa/consol/parser"
	"github.com/blackzarifa/consol/tui"
)

func main() {
	if len(os.Args) >= 2 {
		processFile(os.Args[1])
		return
	}

	file, err := os.ReadFile(os.Args[1])
}

func processFile(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	content := string(file)

	if !parser.HasConflict(content) {
		fmt.Println("No conflicts found in file.")
		return
	}

	conflicts, normalized, lineEnding := parser.ParseFile(content)

	normalizedArr := strings.Split(normalized, "\n")
	lastLine := len(normalizedArr) - 1
	if len(normalizedArr) > 0 && normalizedArr[lastLine] == "" {
		normalizedArr = normalizedArr[:lastLine]
	}

	tui.RunProgram(normalizedArr, lineEnding, filename, conflicts)
}
}
