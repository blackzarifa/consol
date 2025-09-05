package tui

import (
	"slices"
	"strings"

	"github.com/blackzarifa/consol/parser"
)

func RemoveIndex[T any](slice []T, index int) []T {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}

func (m *model) resolveConflict(resolution string) {
	if len(m.conflicts) == 0 {
		return
	}

	cc := m.conflicts[m.currentConflict]
	resLines := strings.Split(resolution, "\n")

	for len(m.resolvedLines) < len(m.normalized) {
		m.resolvedLines = append(m.resolvedLines, false)
	}

	newResolvedLines := make([]bool, len(resLines))
	for i := range newResolvedLines {
		newResolvedLines[i] = true
	}

	m.normalized = slices.Replace(m.normalized, cc.StartLine-1, cc.EndLine, resLines...)
	m.resolvedLines = slices.Replace(
		m.resolvedLines,
		cc.StartLine-1,
		cc.EndLine,
		newResolvedLines...,
	)

	m.conflicts = RemoveIndex(m.conflicts, m.currentConflict)

	originalNumLines := cc.EndLine - cc.StartLine + 1
	lineDiff := len(resLines) - originalNumLines

	m.updateConflictLines(cc, lineDiff)
	m.adjustCurrentConflict()
}

// updateConflictLines updates the start and end lines in each other conflict after resolved,
// according to the given line diff
func (m *model) updateConflictLines(resolved parser.Conflict, lineDiff int) {
	for i := range m.conflicts {
		if m.conflicts[i].StartLine <= resolved.StartLine {
			continue
		}
		m.conflicts[i].StartLine += lineDiff
		m.conflicts[i].EndLine += lineDiff
	}
}

// adjustCurrentConflict fixes the value of the m.currentConflict if it's higher than the amount of
// conflicts
func (m *model) adjustCurrentConflict() {
	if len(m.conflicts) == 0 {
		m.currentConflict = 0
		return
	}

	if m.currentConflict >= len(m.conflicts) {
		m.currentConflict = len(m.conflicts) - 1
	}

	m.cursor = m.conflicts[m.currentConflict].StartLine - 1
	m.updateViewportContent()
}

func (m *model) cursorToConflict() {
	for i, c := range m.conflicts {
		if m.cursor >= c.StartLine-1 && m.cursor <= c.EndLine-1 {
			m.currentConflict = i
		}
	}
}

func (m *model) isLineResolved(visualLineNum int) bool {
	normalizedIndex := m.visualLineToNormalizedIndex(visualLineNum)
	if normalizedIndex < 0 || normalizedIndex >= len(m.resolvedLines) {
		return false
	}
	return m.resolvedLines[normalizedIndex]
}

// visualLineToNormalizedIndex converts a visual line number to normalized slice index
func (m *model) visualLineToNormalizedIndex(visualLineNum int) int {
	currentVisualLine := 1
	state := "normal" // normal, ours, theirs

	for i, line := range m.normalized {
		var lineType string
		if parser.ConflictStart.MatchString(line) {
			lineType = "conflictStart"
			state = "ours"
		} else if parser.ConflictSeparator.MatchString(line) {
			lineType = "conflictSeparator"
			state = "theirs"
		} else if parser.ConflictEnd.MatchString(line) {
			lineType = "conflictEnd"
			state = "normal"
		}

		showLineNumber := (state == "normal" || state == "ours") && lineType == ""

		if showLineNumber {
			if currentVisualLine == visualLineNum {
				return i
			}
			currentVisualLine++
		}
	}

	return -1
}
