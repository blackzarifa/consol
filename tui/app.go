// Package tui runs the terminal interface using bubbletea
package tui

import (
	"log"

	"github.com/blackzarifa/consol/parser"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func RunConflictResolver(
	normalizedArr []string,
	lineEnding, filename string,
	conflicts []parser.Conflict,
) bool {
	p := tea.NewProgram(
		newConflictResolver(normalizedArr, lineEnding, filename, conflicts),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}

	if m := finalModel.(model); m.backToSelector {
		return true
	}
	return false
}

func newConflictResolver(
	normalizedArr []string,
	lineEnding, filename string,
	conflicts []parser.Conflict,
) model {
	initialCursor := 0
	currentConflict := 0

	if len(conflicts) > 0 {
		initialCursor = conflicts[0].StartLine - 1
	}

	m := model{
		resolvedLines:   make([]bool, len(normalizedArr)),
		conflicts:       conflicts,
		normalized:      normalizedArr,
		lineEnding:      lineEnding,
		currentConflict: currentConflict,
		cursor:          initialCursor,
		viewport:        viewport.New(50, 25),
		filename:        filename,
		help:            help.New(),
	}

	m.updateViewportContent()
	return m
}
