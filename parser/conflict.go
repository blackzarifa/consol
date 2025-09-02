// Package parser has utilities for parsing files and checking for conflicts
package parser

import (
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
	ConflictStart     = regexp.MustCompile(`(?m)^<{7}( .*)?`)
	ConflictSeparator = regexp.MustCompile(`(?m)^={7}`)
	ConflictEnd       = regexp.MustCompile(`(?m)^>{7}( .*)?`)
)

func HasConflict(content string) bool {
	return ConflictStart.MatchString(content)
}

// ParseFile parses an entire file string to return conflicts, normalized
// content, and line ending
func ParseFile(content string) ([]Conflict, string, string) {
	var conflicts []Conflict

	normalized, lineEnding := normalizeLineEndings(content)

	startIndexes := ConflictStart.FindAllStringIndex(normalized, -1)
	separatorIndexes := ConflictSeparator.FindAllStringIndex(normalized, -1)
	endIndexes := ConflictEnd.FindAllStringIndex(normalized, -1)

	for i, start := range startIndexes {
		if i >= len(separatorIndexes) || i >= len(endIndexes) {
			break
		}

		separator := separatorIndexes[i]
		end := endIndexes[i]

		conflict := parseConflict(normalized, start, separator, end)
		conflicts = append(conflicts, conflict)
	}

	return conflicts, normalized, lineEnding
}

// Normalizes all line endings in content to \n then returns the normalized
// string and the original line ending
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

// Parses a Conflict from a content string and start, separator and end indexes
func parseConflict(content string, start, separator, end []int) Conflict {
	startLine := getLineNumber(content, start[0])
	endLine := getLineNumber(content, end[0])
	ours := content[start[1]+1 : separator[0]-1]
	theirs := content[separator[1]+1 : end[0]-1]

	return Conflict{
		StartLine: startLine,
		EndLine:   endLine,
		Ours:      ours,
		Theirs:    theirs,
	}
}

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
