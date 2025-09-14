package tui

import (
	"fmt"

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
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "k", "up":
			m = m.moveCursorUp()
		case "j", "down":
			m = m.moveCursorDown()
		case "g", "home":
			m = m.handleGoToStart(msg)
		case "G", "end":
			m = m.goToEnd()
		case "n", "p":
			m = m.navigateConflict(msg.String())
		case "o", "t":
			m = m.resolveCurrentConflict(msg.String())
		case "w", "ctrl+s":
			m, cmd = m.saveFile()
		}
	}

	return m, cmd
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	return renderHeader()
}

func (m model) footerView() string {
	if m.statusMessage != "" {
		return "\n" + m.statusMessage
	}

	return "\n" + m.help.View(keyMap{}) + "\n"
}
