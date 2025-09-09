package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blackzarifa/consol/parser"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct{}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("k"), key.WithHelp("↑/k", "up")),
		key.NewBinding(key.WithKeys("j"), key.WithHelp("↓/j", "down")),
		key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "save")),
		key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "more")),
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(key.WithKeys("k"), key.WithHelp("↑/k", "up")),
			key.NewBinding(key.WithKeys("j"), key.WithHelp("↓/j", "down")),
			key.NewBinding(key.WithKeys("g"), key.WithHelp("g/home", "go to start")),
			key.NewBinding(key.WithKeys("G"), key.WithHelp("G/end", "go to end")),
		},
		{
			key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "previous conflict")),
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next conflict")),
			key.NewBinding(key.WithKeys("w"), key.WithHelp("w/ctrl+s", "save file")),
			key.NewBinding(key.WithKeys("q"), key.WithHelp("q/ctrl+c", "quit")),
		},
		{
			key.NewBinding(key.WithKeys("b"), key.WithHelp("b/esc", "back to file list")),
			key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "close help")),
		},
	}
}

type model struct {
	resolvedLines   []bool
	conflicts       []parser.Conflict
	normalized      []string
	help            help.Model
	viewport        viewport.Model
	filename        string
	lineEnding      string
	statusMessage   string
	currentConflict int
	cursor          int
	lastKeyG        bool
	backToSelector  bool
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
		m.statusMessage = ""
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "b", "esc":
			m.backToSelector = true
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
				m.updateViewportContent()
				break
			}
			m.resolveConflict(cc.Theirs)
			m.updateViewportContent()
		case "w", "ctrl+s":
			toSave := strings.Join(m.normalized, m.lineEnding) + m.lineEnding
			err := os.WriteFile(m.filename, []byte(toSave), 0o664)
			if err != nil {
				m.statusMessage = fmt.Sprintf("Error saving file: %v", err)
			} else {
				m.statusMessage = "File Saved"
			}

			cmd = tea.Tick(
				2*time.Second,
				func(t time.Time) tea.Msg { return tea.KeyMsg{} },
			)
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		}

	}

	return m, cmd
}

func (m model) Init() tea.Cmd {
	tea.SetWindowTitle("Consol - Conflict reSolver")
	return nil
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	title := titleStyle.Render("Consol Conflict reSolver")

	centerStyle := lipgloss.NewStyle().
		Width(m.viewport.Width).
		Align(lipgloss.Center)

	return centerStyle.Render(title)
}

func (m model) footerView() string {
	if m.statusMessage != "" {
		return "\n" + m.statusMessage
	}

	return "\n" + m.help.View(keyMap{}) + "\n"
}
