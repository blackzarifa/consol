package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blackzarifa/consol/parser"
	"github.com/blackzarifa/consol/tui"
)

var version = "dev"

func main() {
	if len(os.Args) >= 2 {
		arg := os.Args[1]

		switch arg {
		case "--version", "-v":
			fmt.Printf("consol %s\n", version)
			return
		case "--help", "-h":
			showUsage()
			return
		default:
			handleFile(arg)
			return
		}
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

func showUsage() {
	fmt.Println("Usage: consol [file]")
	fmt.Println("\nInteractive Git merge conflict resolver")
	fmt.Println("\nCommands:")
	fmt.Println("  consol              - Auto-discover and select conflict files")
	fmt.Println("  consol <file>       - Resolve conflicts in specific file")
	fmt.Println("  consol --version    - Show version")
	fmt.Println("  consol --help       - Show this help")
}

func handleFile(filename string) bool {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: Cannot read file '%s': %v\n", filename, err)
		return false
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
