package tui

import (
	"fmt"
	"strings"
)

func (m *model) updateViewportContent() {
	var lines []string
	for i, line := range m.normalized {
		if i != m.cursor {
			lines = append(lines, line)
			continue
		}

		lines = append(lines, fmt.Sprintf(
			">>> %s <<< Current:%d lenCon:%d",
			line, m.currentConflict, len(m.conflicts),
		))
	}

	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(max(0, m.cursor-m.viewport.Height/2))
}
