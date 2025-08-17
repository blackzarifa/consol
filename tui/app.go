// Package tui runs the terminal interface using bubbletea
package tui

import (
	"fmt"

	"github.com/blackzarifa/consol/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor     int
	normalized string
	lineEnding string
	conflicts  []parser.Conflict
}

func InitialModel(
	normalized, lineEnding string,
	conflicts []parser.Conflict,
) model {
	return model{
		normalized: normalized,
		lineEnding: lineEnding,
		conflicts:  conflicts,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "=== CONSOL CONFLICT RESOLVER ===\n\n"
	s += "File content:\n---\n"
	s += m.normalized
	s += "\n---\n\n"

	s += fmt.Sprintf("Found %d conflicts:\n", len(m.conflicts))
	for i, c := range m.conflicts {
		s += fmt.Sprintf("  %d. Lines %d-%d\n", i+1, c.StartLine, c.EndLine)
	}

	s += fmt.Sprintf("\nOriginal line ending: %q\n", m.lineEnding)
	s += "Press 'q' to quit\n"
	return s
}
