package application

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) HandleKeyMsg(key tea.KeyMsg) tea.Cmd  {
	switch key.String() {
	case "ctrl+c", "q":
		return tea.Quit
	case "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "j":
		if m.cursor < len(m.choices) - 1{
			m.cursor++
		}
	case " ":
		_, ok := m.selected[m.cursor]
		if ok {
			delete(m.selected, m.cursor)
		} else {
			m.selected[m.cursor] = struct{}{}
		}
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if cmd := m.HandleKeyMsg(msg); cmd != nil {
			return m, cmd
		}
	}

	return m, nil
}

func (m model) RenderChoices() string {
	output := ""
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		output += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	return output
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"
	s += m.RenderChoices()
	s += "\nPress q to quit.\n"
	return s
}

func New() *tea.Program {
	return tea.NewProgram(initialModel())
}



