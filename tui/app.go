// Package tui runs the terminal interface using bubbletea
package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/blackzarifa/consol/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	conflicts   []parser.Conflict
	normalized  []string
	lineEnding  string
	contentSize int
	cursor      int
	height      int
	offset      int
}

func RunProgram(normalized, lineEnding string, conflicts []parser.Conflict) {
	normalizedArr := strings.Split(normalized, "\n")
	if len(normalizedArr) > 0 && normalizedArr[len(normalizedArr)-1] == "" {
		normalizedArr = normalizedArr[:len(normalizedArr)-1]
	}

	p := tea.NewProgram(
		initialModel(normalizedArr, lineEnding, conflicts),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

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
				if m.cursor < m.offset {
					m.offset--
				}
			}
		case "j", "down":
			if m.cursor < len(m.normalized)-1 {
				m.cursor++

				if len(m.normalized) <= m.contentSize {
					break
				}

				lastVisibleLine := m.offset + m.contentSize
				linesVisibleBelow := lastVisibleLine - m.cursor

				if linesVisibleBelow < scrolloff && lastVisibleLine < len(m.normalized)-1 {
					m.offset++
				}
			}
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

func initialModel(
	normalizedArr []string, lineEnding string,
	conflicts []parser.Conflict,
) model {
	return model{
		conflicts:  conflicts,
		normalized: normalizedArr,
		lineEnding: lineEnding,
		cursor:     0,
		offset:     0,
	}
}
