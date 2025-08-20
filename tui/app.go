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
	conflicts  []parser.Conflict
	normalized string
	lineEnding string
	cursor     int
	height     int
	offset     int
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
	tea.SetWindowTitle("Consol - Conflict reSolver")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j", "down":
			m.cursor++
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "=== CONSOL CONFLICT RESOLVER ===\n\n"

	lines := strings.Split(m.normalized, "\n")
	for i := range lines {
		if i >= m.height-5 {
			break
		}

		if i == m.cursor {
			s += fmt.Sprintf(">>> %s <<<\n", lines[i])
			continue
		}

		s += lines[i] + "\n"
	}

	s += "\nPress 'q' to quit\n"
	return s
}

func initialModel(
	normalized, lineEnding string,
	conflicts []parser.Conflict,
) model {
	return model{
		conflicts:  conflicts,
		normalized: normalized,
		lineEnding: lineEnding,
		cursor:     0,
		offset:     0,
	}
}
