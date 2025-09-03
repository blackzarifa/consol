package tui

import (
	"strings"

	"github.com/blackzarifa/consol/parser"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateViewportContent() {
	oursBranch := lipgloss.NewStyle().
		Background(lipgloss.Color("28")).
		Foreground(lipgloss.Color("15")).Bold(true).
		Width(m.viewport.Width)
	theirsBranch := lipgloss.NewStyle().
		Background(lipgloss.Color("19")).
		Foreground(lipgloss.Color("15")).Bold(true).
		Width(m.viewport.Width)
	oursStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("22")).
		Foreground(lipgloss.Color("15")).
		Width(m.viewport.Width)
	theirsStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("18")).
		Foreground(lipgloss.Color("15")).
		Width(m.viewport.Width)

	var lines []string
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

		styledLine := m.styleLineSegment(
			line, lineType, state, oursBranch, theirsBranch, oursStyle, theirsStyle,
		)

		if i == m.cursor {
			cursorStyle := lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Dark: "255", Light: "0"}).
				Foreground(lipgloss.AdaptiveColor{Dark: "0", Light: "255"}).
				Blink(true)

			if len(line) == 0 {
				line = " "
			}

			firstChar := cursorStyle.Render(string(line[0]))
			restStyled := m.styleLineSegment(
				line[1:], lineType, state, oursBranch, theirsBranch, oursStyle, theirsStyle,
			)
			styledLine = firstChar + restStyled
		}

		lines = append(lines, styledLine)
	}

	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(max(0, m.cursor-m.viewport.Height/2))
}

func (m *model) styleLineSegment(
	line, lineType, state string,
	oursBranch, theirsBranch, oursStyle, theirsStyle lipgloss.Style,
) string {
	switch lineType {
	case "conflictStart":
		return oursBranch.Render(line)
	case "conflictSeparator":
		return line
	case "conflictEnd":
		return theirsBranch.Render(line)
	default:
		switch state {
		case "ours":
			return oursStyle.Render(line)
		case "theirs":
			return theirsStyle.Render(line)
		default:
			return line
		}
	}
}
