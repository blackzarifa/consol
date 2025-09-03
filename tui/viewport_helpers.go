package tui

import (
	"fmt"
	"strings"

	"github.com/blackzarifa/consol/parser"
	"github.com/charmbracelet/lipgloss"
)

const (
	lineNumWidth = 5
	spacing      = 2
)

func (m *model) updateViewportContent() {
	var lines []string
	state := "normal" // normal, ours, theirs
	lineNumber := 1

	for i, line := range m.normalized {
		var lineType string
		if parser.ConflictStart.MatchString(line) {
			lineType = "conflictStart"
			state = "ours"
			lines = append(lines, m.renderConflictMessage())
		} else if parser.ConflictSeparator.MatchString(line) {
			lineType = "conflictSeparator"
			state = "theirs"
		} else if parser.ConflictEnd.MatchString(line) {
			lineType = "conflictEnd"
			state = "normal"
		}

		showLineNumber := (state == "normal" || state == "ours") && lineType == ""
		lineNumStr, isResolved := m.renderLineNumber(lineNumber, line, showLineNumber)
		if showLineNumber {
			lineNumber++
		}

		displayLine := m.renderLine(line, lineNumStr, isResolved, i == m.cursor, lineType, state)
		lines = append(lines, displayLine)
	}

	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(max(0, m.cursor-m.viewport.Height/2))
}

func (m *model) renderConflictMessage() string {
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Dark: "246", Light: "240"}).
		Width(m.viewport.Width)

	acceptMessage := messageStyle.Render("Accept incoming change (o) | Ignore (t)")
	messagePrefix := strings.Repeat(" ", lineNumWidth+spacing)
	return messagePrefix + acceptMessage
}

func (m *model) renderLineNumber(lineNumber int, line string, showLineNumber bool) (string, bool) {
	if showLineNumber {
		lineNumStr := fmt.Sprintf("%*d", lineNumWidth, lineNumber)

		if m.resolvedLines != nil && m.resolvedLines[strings.TrimSpace(line)] {
			lineNumStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("136")).
				BorderRight(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("136"))
			return lineNumStyle.Render(lineNumStr), true
		}
		return lineNumStr, false
	}
	return strings.Repeat(" ", lineNumWidth), false
}

func (m *model) renderLine(
	line, lineNumStr string,
	isResolved, isCursor bool,
	lineType, state string,
) string {
	lineSpacing := strings.Repeat(" ", spacing)
	if isResolved {
		lineSpacing = lineSpacing[0 : len(lineSpacing)-1]
	}

	if isCursor {
		cursorStyle := lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Dark: "255", Light: "0"}).
			Foreground(lipgloss.AdaptiveColor{Dark: "0", Light: "255"}).
			Blink(true)

		if len(line) == 0 {
			line = " "
		}

		firstChar := cursorStyle.Render(string(line[0]))
		restStyled := m.styleLineSegment(line[1:], lineType, state)
		return lineNumStr + lineSpacing + firstChar + restStyled
	}

	styledLine := m.styleLineSegment(line, lineType, state)
	return lineNumStr + lineSpacing + styledLine
}

func (m *model) styleLineSegment(line, lineType, state string) string {
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
