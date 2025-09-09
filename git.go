package main

import (
	"os/exec"
	"strings"
)

func findConflictFiles() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	conflictFiles := []string{}
	trimmedOutput := strings.TrimSpace(string(output))
	for line := range strings.SplitSeq(trimmedOutput, "\n") {
		if len(line) >= 3 && line[:2] == "UU" {
			filename := strings.TrimSpace(line[3:])
			conflictFiles = append(conflictFiles, filename)
		}
	}

	return conflictFiles, nil
}