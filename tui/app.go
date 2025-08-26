// Package tui runs the terminal interface using bubbletea
package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/blackzarifa/consol/parser"
	tea "github.com/charmbracelet/bubbletea"
)

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
