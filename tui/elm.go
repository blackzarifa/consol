package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/blackzarifa/consol/parser"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	resolvedLines   []bool
	conflicts       []parser.Conflict
	normalized      []string
	viewport        viewport.Model
	lineEnding      string
	currentConflict int
	cursor          int
	lastKeyG        bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.updateViewportContent()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
				m.cursorToConflict()
				m.updateViewportContent()
			}
		case "j", "down":
			if m.cursor < len(m.normalized)-1 {
				m.cursor++
				m.cursorToConflict()
				m.updateViewportContent()
			}
		case "g", "home":
			if msg.String() == "g" && !m.lastKeyG {
				m.lastKeyG = true
				break
			}
			m.lastKeyG = false
			m.cursor = 0
			m.updateViewportContent()
		case "G", "end":
			m.cursor = len(m.normalized) - 1
			m.updateViewportContent()
		case "n", "p":
			if msg.String() == "n" {
				if m.currentConflict >= len(m.conflicts)-1 {
					break
				}
				m.currentConflict++
			} else {
				if m.currentConflict <= 0 {
					break
				}
				m.currentConflict--
			}
			m.cursor = m.conflicts[m.currentConflict].StartLine - 1
			m.updateViewportContent()
		case "o", "t":
			if len(m.conflicts) == 0 {
				break
			}
			cc := m.conflicts[m.currentConflict]

			if msg.String() == "o" {
				m.resolveConflict(cc.Ours)
				break
			}
			m.resolveConflict(cc.Theirs)
		case "w":
			toSave := strings.Join(m.normalized, m.lineEnding) + m.lineEnding
			os.WriteFile(os.Args[1], []byte(toSave), 0o664)
		}

	}

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	return "=== CONSOL CONFLICT RESOLVER ===\n"
}

func (m model) footerView() string {
	footer := "\n'q' to quit  |  'jknp' to navigate"
	if m.lastKeyG {
		footer += "  |  'g' - Press again to go to the beginning"
	}
	return footer
}
