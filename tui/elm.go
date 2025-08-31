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
			m.cursor = len(m.normalized) - 1
			m.offset = m.calculateOffset(m.cursor)
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
		case "o", "t":
			if len(m.conflicts) == 0 {
				break
			}
			cc := m.conflicts[m.currentConflict]

			if msg.String() == "o" {
				m.resolveConflict(cc.Ours)
			}
			m.resolveConflict(cc.Theirs)
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
			s += fmt.Sprintf(
				">>> %s <<< Current:%d lenCon:%d\n",
				line, m.currentConflict, len(m.conflicts),
			)
			continue
		}
		s += line + "\n"
	}

	if s[len(s)-1:] == "\n" {
		s = s[:len(s)-1]
	}

	s += "\n\nPress 'q' to quit  |  'np' to navigate conflicts"
	if m.lastKeyG {
		s += "  |  'g' - Press again to go to the beginning"
	}
	return s
}
