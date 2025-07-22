package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

	if !strings.Contains(content, "<<<<<<< HEAD") {
		fmt.Println("No conflicts found!")
		return
	}
}
