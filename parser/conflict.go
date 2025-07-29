package parser

import (
	"fmt"
	"regexp"
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
		// if content[i] == '\r' {
		// 	if content[i+1] == '\n' {
		// 		line += 2
		// 		continue
		// 	}
		// 	line++
		// 	continue
		// }
	}
	return line
}

func HasConflict(content string) bool {
	if !conflictStart.MatchString(content) {
		return false
	}
	return true
}

func ParseFile(content string) ([]Conflict, string) {
	var conflicts []Conflict

	startIndexes := conflictStart.FindAllStringIndex(content, -1)
	separatorIndexes := conflictSeparator.FindAllStringIndex(content, -1)
	endIndexes := conflictEnd.FindAllStringIndex(content, -1)

	for i, start := range startIndexes {
		if i >= len(separatorIndexes) || i >= len(endIndexes) {
			break
		}

		separator := separatorIndexes[i]
		end := endIndexes[i]

		startLine := getLineNumber(content, start[0])
		separatorLine := getLineNumber(content, separator[0])
		endLine := getLineNumber(content, end[0])

		fmt.Println(startLine, separatorLine, endLine)
	}

	return conflicts, ""
}
