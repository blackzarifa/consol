package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
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
	regex := regexp.MustCompile(`(?m)^<{7} \S*`)

	if !regex.MatchString(content) {
		fmt.Println("No conflicts found.")
		return
	}
}
