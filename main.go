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
		handleFile(os.Args[1])
		return
	}

	for {
		filename, err := selectFile()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if filename == "" {
			return
		}

		backToSelector := handleFile(filename)
		if !backToSelector {
			return
		}
	}
}

func handleFile(filename string) bool {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	content := string(file)

	if !parser.HasConflict(content) {
		fmt.Println("No conflicts found in file.")
		return true
	}

	conflicts, normalized, lineEnding := parser.ParseFile(content)

	normalizedArr := strings.Split(normalized, "\n")
	lastLine := len(normalizedArr) - 1
	if len(normalizedArr) > 0 && normalizedArr[lastLine] == "" {
		normalizedArr = normalizedArr[:lastLine]
	}

	return tui.RunConflictResolver(normalizedArr, lineEnding, filename, conflicts)
}

func selectFile() (string, error) {
	files, err := findConflictFiles()
	if err != nil {
		return "", fmt.Errorf("finding conflict files: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No git conflict files found in this repository.")
		return "", nil
	}

	conflictCounts := make([]int, len(files))
	for i, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			conflictCounts[i] = 0
			continue
		}

		if parser.HasConflict(string(content)) {
			conflicts, _, _ := parser.ParseFile(string(content))
			conflictCounts[i] = len(conflicts)
		} else {
			conflictCounts[i] = 0
		}
	}

	return tui.RunFileSelector(files, conflictCounts)
}
