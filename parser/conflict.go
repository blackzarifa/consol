package parser

import (
	"regexp"
)

type Conflict struct {
	Startline int
	Endline   int
	Ours      string
	Theirs    string
}

var (
	conflictStart     = regexp.MustCompile(`(?m)^<{7} \S*`)
	conflictSeparator = regexp.MustCompile(`(?m)^={7}\S$`)
	conflictEnd       = regexp.MustCompile(`(?m)^>{7} \S*`)
)

func HasConflict(content string) bool {
	if !conflictStart.MatchString(content) {
		return false
	}
	return true
}

func ParseFile(content string) ([]Conflict, string) {
	var conflicts = []Conflict

	return conflicts, ""
}
