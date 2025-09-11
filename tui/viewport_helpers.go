package tui

import (
	"fmt"
	"strings"
)

const (
	spacing = 2
)

func (m *model) updateViewportContent() {
	var lines []string
	lineNumber := 1
	state := "normal"

	lineNumWidth := len(fmt.Sprintf("%d", len(m.normalized)))

	for i, line := range m.normalized {
		lineType, newState, showLineNumber := processConflictLine(line, state)
		state = newState

		realLineNumber := i + 1
		isCurrentConflictLine := realLineNumber == m.conflicts[m.currentConflict].StartLine
		if lineType == "conflictStart" && isCurrentConflictLine {
			lines = append(lines, m.renderConflictMessage(m.viewport.Width))
		}

		isResolved := i < len(m.resolvedLines) && m.resolvedLines[i]
		lineNumStr := m.renderLineNumber(lineNumber, showLineNumber, isResolved, lineNumWidth)
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

func (m *model) renderConflictMessage(screenWidth int) string {
	messageStyle := ConflictMessageStyle.
		Width(screenWidth)

	lineNumWidth := len(fmt.Sprintf("%d", len(m.normalized)))
	messagePrefix := strings.Repeat(" ", lineNumWidth+spacing)
	acceptMessage := "Accept change (o) | Ignore (t)"
	conflictNum := fmt.Sprintf(" | Conflict %d/%d", m.currentConflict+1, len(m.conflicts))

	return messageStyle.Render(messagePrefix + acceptMessage + conflictNum)
}

func (m *model) renderLineNumber(
	lineNumber int,
	showLineNumber bool,
	isResolved bool,
	lineNumWidth int,
) string {
	if !showLineNumber {
		return strings.Repeat(" ", lineNumWidth)
	}

	lineNumStr := fmt.Sprintf("%*d", lineNumWidth, lineNumber)

	if !isResolved {
		return lineNumStr
	}

	return ResolvedLineNumStyle.Render(lineNumStr)
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

	if len(line) == 0 {
		line = " "
	}

	firstChar := CursorStyle.Render(string(line[0]))
	restStyled := m.styleLineSegment(line[1:], lineType, state)
	return lineNumStr + lineSpacing + firstChar + restStyled
}

func (m *model) styleLineSegment(line, lineType, state string) string {
	switch lineType {
	case "conflictStart":
		return OursBranchStyle.Width(m.viewport.Width).Render(line)
	case "conflictSeparator":
		return line
	case "conflictEnd":
		return TheirsBranchStyle.Width(m.viewport.Width).Render(line)
	default:
		switch state {
		case "ours":
			return OursContentStyle.Width(m.viewport.Width).Render(line)
		case "theirs":
			return TheirsContentStyle.Width(m.viewport.Width).Render(line)
		default:
			return line
		}
	}
}
