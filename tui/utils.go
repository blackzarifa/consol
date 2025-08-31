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

	m.normalized = slices.Replace(
		m.normalized, cc.StartLine-1, cc.EndLine, resLines...,
	)
	m.conflicts = RemoveIndex(m.conflicts, m.currentConflict)

	originalNumLines := cc.EndLine - cc.StartLine + 1
	lineDiff := len(resLines) - originalNumLines

	m.updateConflictLines(cc, lineDiff)
	m.adjustCurrentConflict()
}

func (m *model) updateConflictLines(
	resolved parser.Conflict,
	lineDiff int,
) {
	for i := range m.conflicts {
		if m.conflicts[i].StartLine <= resolved.StartLine {
			continue
		}
		m.conflicts[i].StartLine += lineDiff
		m.conflicts[i].EndLine += lineDiff
	}
}

func (m *model) adjustCurrentConflict() {
	if len(m.conflicts) == 0 {
		m.currentConflict = 0
		return
	}

	if m.currentConflict >= len(m.conflicts) {
		m.currentConflict = len(m.conflicts) - 1
	}

	m.cursor = m.conflicts[m.currentConflict].StartLine - 1
	m.offset = m.calculateOffset(m.cursor)
}

func (m model) calculateOffset(cursor int) int {
	centered := cursor - m.contentSize/2
	maxOffset := len(m.normalized) - m.contentSize - 1
	return max(0, min(centered, maxOffset))
}
