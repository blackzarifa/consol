package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type fileItem struct {
	filename      string
	conflictCount int
}

func (f fileItem) FilterValue() string { return f.filename }
func (f fileItem) Title() string       { return f.filename }
func (f fileItem) Description() string {
	if f.conflictCount == 1 {
		return "1 conflict"
	}
	return fmt.Sprintf("%d conflicts", f.conflictCount)
}

type fileSelectorModel struct {
	list   list.Model
	choice string
	quit   bool
}

func RunFileSelector(files []string, conflictCounts []int) (string, error) {
	if len(files) == 0 {
		fmt.Println("No conflict files found.")
		return "", nil
	}

	m := newFileSelectorModel(files, conflictCounts)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m := finalModel.(fileSelectorModel); m.choice != "" {
		return m.choice, nil
	}

	return "", nil
}

func newFileSelectorModel(files []string, conflictCounts []int) fileSelectorModel {
	items := make([]list.Item, len(files))
	for i, file := range files {
		conflictCount := 0
		if i < len(conflictCounts) {
			conflictCount = conflictCounts[i]
		}
		items[i] = fileItem{filename: file, conflictCount: conflictCount}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = FileListNormalTitleStyle
	delegate.Styles.NormalDesc = FileListNormalDescStyle
	delegate.Styles.SelectedTitle = FileListSelectedTitleStyle
	delegate.Styles.SelectedDesc = FileListSelectedDescStyle

	l := list.New(items, delegate, 50, 25)
	l.Title = "Choose a file to resolve the conflicts:"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = BlankStyle
	l.Styles.HelpStyle = FooterStyle
	l.Styles.PaginationStyle = BlankStyle

	return fileSelectorModel{list: l}
}

func (m fileSelectorModel) Init() tea.Cmd {
	return nil
}

func (m fileSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(renderHeader())
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - headerHeight - 1)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(fileItem)
			if ok {
				m.choice = i.filename
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m fileSelectorModel) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quit {
		return FooterStyle.Render("No file selected.")
	}
	return renderHeader() + "\n" + m.list.View()
}
