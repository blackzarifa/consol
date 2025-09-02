package tui

import (
	"strings"

	"github.com/blackzarifa/consol/parser"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateViewportContent() {
	oursBranch := lipgloss.NewStyle().
		Background(lipgloss.Color("28")).
		Foreground(lipgloss.Color("15")).Bold(true)
	theirsBranch := lipgloss.NewStyle().
		Background(lipgloss.Color("19")).
		Foreground(lipgloss.Color("15")).Bold(true)
	oursStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("22")).
		Foreground(lipgloss.Color("15"))
	theirsStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("18")).
		Foreground(lipgloss.Color("15"))

	var lines []string
	state := "normal" // normal, ours, theirs

	for i, line := range m.normalized {
		styledLine := line

		if parser.ConflictStart.MatchString(line) {
			state = "ours"
			styledLine = oursBranch.Render(line)
		} else if parser.ConflictSeparator.MatchString(line) {
			state = "theirs"
		} else if parser.ConflictEnd.MatchString(line) {
			state = "normal"
			styledLine = theirsBranch.Render(line)
		} else {
			switch state {
			case "ours":
				styledLine = oursStyle.Render(line)
			case "theirs":
				styledLine = theirsStyle.Render(line)
			}
		}

		if i == m.cursor {
			cursorStyle := lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Dark: "238", Light: "252"}).
				Foreground(lipgloss.Color("15"))
			styledLine = cursorStyle.Render(line)
		}

		lines = append(lines, styledLine)
	}

	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(max(0, m.cursor-m.viewport.Height/2))
}
