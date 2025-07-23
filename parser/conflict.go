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
	startConflict = regexp.MustCompile(`(?m)^<{7} \S*`)
	endConflict   = regexp.MustCompile(`(?m)^>{7} \S*`)
)

func HasConflict(content string) bool {
	if !startConflict.MatchString(content) {
		return false
	}
	return true
}

func ParseFile(content string) ([]Conflict, string) {
	// TODO
	return []Conflict{}, ""
}
