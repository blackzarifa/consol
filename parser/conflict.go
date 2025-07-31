package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type Conflict struct {
	StartLine int
	EndLine   int
	Ours      string
	Theirs    string
}

var (
	conflictStart     = regexp.MustCompile(`(?m)^<{7}( .*)?`)
	conflictSeparator = regexp.MustCompile(`(?m)^={7}`)
	conflictEnd       = regexp.MustCompile(`(?m)^>{7}( .*)?`)
)

func getLineNumber(content string, pos int) int {
	line := 1
	for i := 0; i < pos && i < len(content); i++ {
		if content[i] == '\n' {
			line++
			continue
		}

		if content[i] == '\r' {
			if content[i+1] == '\n' {
				i++
			}
			line++
		}
	}
	return line
}

// Normalized all line endings in content to \n -
// Returns the normalized string and the original line ending
func normalizeLineEndings(content string) (normalized, lineEnding string) {
	lineEnding = "\n"
	if strings.Contains(content, "\r\n") {
		lineEnding = "\r\n"
	} else if strings.Contains(content, "\r") {
		lineEnding = "\r"
	}
	normalized = strings.ReplaceAll(content, lineEnding, "\n")

	return
}

func HasConflict(content string) bool {
	if !conflictStart.MatchString(content) {
		return false
	}
	return true
}

// Parses an entire file string to return an array of Conflict
func ParseFile(content string) ([]Conflict, string) {
	var conflicts []Conflict

	normalized, lineEnding := normalizeLineEndings(content)

	startIndexes := conflictStart.FindAllStringIndex(normalized, -1)
	separatorIndexes := conflictSeparator.FindAllStringIndex(normalized, -1)
	endIndexes := conflictEnd.FindAllStringIndex(normalized, -1)

	for i, start := range startIndexes {
		if i >= len(separatorIndexes) || i >= len(endIndexes) {
			break
		}

		separator := separatorIndexes[i]
		end := endIndexes[i]

		startLine := getLineNumber(normalized, start[0])
		endLine := getLineNumber(normalized, end[0])
		ours := normalized[start[1]+1 : separator[0]-1]
		theirs := normalized[separator[1]+1 : end[0]-1]

		conflict := Conflict{StartLine: startLine, EndLine: endLine, Ours: ours, Theirs: theirs}
		conflicts = append(conflicts, conflict)
		fmt.Println(conflict.Ours)
		fmt.Println("----------------")
		fmt.Println(conflict.Theirs)
		fmt.Println("================")
	}

	return conflicts, ""
}
