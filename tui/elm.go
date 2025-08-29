package tui

import (
	"fmt"

	"github.com/blackzarifa/consol/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	conflicts       []parser.Conflict
	normalized      []string
	lineEnding      string
	contentSize     int
	currentConflict int
	cursor          int
	height          int
	offset          int
	lastKeyG        bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.contentSize = m.height - 5 // 5 is the size of the header + footnote

	case tea.KeyMsg:
		scrolloff := 10

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--

				if m.offset > 0 && m.cursor < m.offset+scrolloff {
					m.offset--
				}
			}
		case "j", "down":
			if m.cursor < len(m.normalized)-1 {
				m.cursor++

				lastVisibleLine := m.offset + m.contentSize
				linesVisibleBelow := lastVisibleLine - m.cursor

				if linesVisibleBelow < scrolloff && lastVisibleLine < len(m.normalized)-1 {
					m.offset++
				}
			}
		case "g", "home":
			if msg.String() == "g" && !m.lastKeyG {
				m.lastKeyG = true
				break
			}
			m.lastKeyG = false
			m.cursor = 0
			m.offset = 0
		case "G", "end":
			length := len(m.normalized) - 1
			m.cursor = length
			m.offset = length - m.contentSize
		case "n":
			if m.currentConflict >= len(m.conflicts)-1 {
				break
			}
			m.currentConflict++
			m.cursor = m.conflicts[m.currentConflict].StartLine
			m.offset = m.calculateOffset(m.cursor)
		case "p":
			if m.currentConflict <= 0 {
				break
			}
			m.currentConflict--
			m.cursor = m.conflicts[m.currentConflict].StartLine
			m.offset = m.calculateOffset(m.cursor)
		}

	}

	return m, nil
}

func (m model) View() string {
	s := "=== CONSOL CONFLICT RESOLVER ===\n\n"

	for i, line := range m.normalized {
		if i > m.contentSize+m.offset {
			break
		} else if i < m.offset {
			continue
		} else if i == m.cursor {
			s += fmt.Sprintf(">>> %s <<< Length:%d CS:%d c:%d off:%d\n", line, len(m.normalized), m.contentSize, m.cursor, m.offset)
			continue
		}

		s += fmt.Sprintf("%d  -  ", i)
		s += line + "\n"
	}

	if s[len(s)-1:] == "\n" {
		s = s[:len(s)-1]
	}

	s += "\n\nPress 'q' to quit"
	if m.lastKeyG {
		s += "  |  'g' - Press again to go to the beginning"
	}
	return s
}

func (m model) calculateOffset(cursor int) int {
	centered := cursor - m.contentSize/2
	maxOffset := len(m.normalized) - m.contentSize - 1
	return max(0, min(centered, maxOffset))
}
