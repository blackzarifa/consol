// Package tui runs the terminal interface using bubbletea
package tui

import (
	"fmt"
	"os"

	"github.com/blackzarifa/consol/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	conflicts    []parser.Conflict
	currentIndex int
	cursor       int
	normalized   string
	lineEnding   string
}

func RunProgram(normalized, lineEnding string, conflicts []parser.Conflict) {
	p := tea.NewProgram(
		initialModel(normalized, lineEnding, conflicts),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
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

	s += m.normalized

	s += "Press 'q' to quit\n"
	return s
}

func initialModel(
	normalized, lineEnding string,
	conflicts []parser.Conflict,
) model {
	return model{
		normalized: normalized,
		lineEnding: lineEnding,
		conflicts:  conflicts,
	}
}
