package tui

import (
	"slices"
	"strings"
)

func RemoveIndex[T any](slice []T, index int) []T {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}

func (m model) calculateOffset(cursor int) int {
	centered := cursor - m.contentSize/2
	maxOffset := len(m.normalized) - m.contentSize - 1
	return max(0, min(centered, maxOffset))
}

func (m *model) resolveConflict(resolution string) {
	if len(m.conflicts) == 0 {
		return
	}

	cIndex := m.currentConflict
	cc := m.conflicts[cIndex]
	resLines := strings.Split(resolution, "\n")

	m.normalized = slices.Replace(
		m.normalized, cc.StartLine-1, cc.EndLine, resLines...,
	)
	m.conflicts = RemoveIndex(m.conflicts, cIndex)

	originalNumLines := cc.EndLine - cc.StartLine + 1
	lineDiff := len(resLines) - originalNumLines

	for i := range m.conflicts {
		if m.conflicts[i].StartLine <= cc.StartLine {
			continue
		}
		m.conflicts[i].StartLine += lineDiff
		m.conflicts[i].EndLine += lineDiff
	}

	if len(m.conflicts) == 0 {
		m.currentConflict = 0
		return
	}

	if cIndex >= len(m.conflicts) {
		cIndex = len(m.conflicts) - 1
	}

	m.currentConflict = cIndex
	m.cursor = m.conflicts[cIndex].StartLine - 1
	m.offset = m.calculateOffset(m.cursor)
}
