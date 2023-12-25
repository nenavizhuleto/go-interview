package application

import (

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	help help.Model
	cols map[status]*column
	focused status
}

func NewApp() *App {
	help := help.New()
	help.ShowAll = true
	return &App{
		help:  help,
	}
}

func (m *App) Init() tea.Cmd {
	m.focused = StatusTODO
	return nil
}

func (m *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		var cmds []tea.Cmd
		for _, val := range m.cols {
			_, cmd := val.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Right):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		case key.Matches(msg, keys.Left):
			m.cols[m.focused].Blur()
			m.focused = m.focused.getNext()
			m.cols[m.focused].Focus()
		}
	case Form:
		return m, m.cols[m.focused].Set(msg.index, msg.CreateTodo())
	case moveMsg:
		return m, m.cols[m.focused.getNext()].Set(APPEND, msg.Todo)
	}

	res, cmd := m.cols[m.focused].Update(msg)
	if col, ok := res.(column); ok {
		m.cols[m.focused] = &col
	} else {
		return res, cmd
	}
	return m, cmd

}

func (m *App) View() string {
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[StatusTODO].View(),
		m.cols[StatusDone].View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left,	board, m.help.View(keys))
}
