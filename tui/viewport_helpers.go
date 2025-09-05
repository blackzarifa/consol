package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	lineNumWidth = 5
	spacing      = 2
)

func (m *model) updateViewportContent() {
	var lines []string
	lineNumber := 1
	state := "normal"

	for i, line := range m.normalized {
		lineType, newState, showLineNumber := processConflictLine(line, state)
		state = newState

		if lineType == "conflictStart" {
			lines = append(lines, renderConflictMessage(m.viewport.Width))
		}
		isResolved := i < len(m.resolvedLines) && m.resolvedLines[i]
		lineNumStr := m.renderLineNumber(lineNumber, showLineNumber, isResolved)
		if showLineNumber {
			lineNumber++
		}

		displayLine := m.renderLine(line, lineNumStr, lineType, newState, isResolved, i == m.cursor)
		lines = append(lines, displayLine)
	}

	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)
	m.viewport.SetYOffset(max(0, m.cursor-m.viewport.Height/2))
}

func renderConflictMessage(screenWidth int) string {
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Dark: "246", Light: "240"}).
		Width(screenWidth)

	acceptMessage := messageStyle.Render("Accept incoming change (o) | Ignore (t)")
	messagePrefix := strings.Repeat(" ", lineNumWidth+spacing)
	return messagePrefix + acceptMessage
}

func (m *model) renderLineNumber(lineNumber int, showLineNumber bool, isResolved bool) string {
	if !showLineNumber {
		return strings.Repeat(" ", lineNumWidth)
	}

	lineNumStr := fmt.Sprintf("%*d", lineNumWidth, lineNumber)

	if !isResolved {
		return lineNumStr
	}

	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("136")).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("136"))
	return lineNumStyle.Render(lineNumStr)
}

func (m *model) renderLine(
	line, lineNumStr, lineType, state string,
	isResolved, isCursor bool,
) string {
	lineSpacing := strings.Repeat(" ", spacing)
	if isResolved {
		lineSpacing = lineSpacing[0 : len(lineSpacing)-1]
	}

	if !isCursor {
		styledLine := m.styleLineSegment(line, lineType, state)
		return lineNumStr + lineSpacing + styledLine
	}

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
