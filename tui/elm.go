package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	tea.SetWindowTitle("Consol - Conflict reSolver")
	return nil
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
		case "G", "end":
			length := len(m.normalized) - 1
			m.cursor = length
			m.offset = length - m.contentSize
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
	return s
}
