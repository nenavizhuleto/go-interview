package application

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Todo struct {
	title       string
	description string
	done        bool
}

func (t Todo) FilterValue() string {
	return t.title
}

func (t Todo) Title() string {
	return t.title
}

func (t Todo) Description() string {
	return t.description
}

func (t *Todo) Toggle() {
	t.done = !t.done
}

type Status string

var (
	TODO = Status("todo")
	DONE = Status("done")
)

func (s Status) Switch() Status {
	if s == TODO {
		return DONE
	}
	return TODO
}

type Model struct {
	Lists    map[Status]*list.Model
	Focused  Status
	Quitting bool
	Loading  bool
}

func (m *Model) Init() tea.Cmd {
	m.Focused = TODO
	m.Loading = true
	return nil
}

func (m *Model) FocusedList() *list.Model {
	return m.Lists[m.Focused]
}

var (
	listStyle        = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.HiddenBorder())
	focusedListStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
)

func (m *Model) RenderList(status Status) string {
	if m.Focused == status {
		return focusedListStyle.Render(m.Lists[status].View())
	} else {
		return listStyle.Render(m.Lists[status].View())
	}

}

func (m *Model) ToggleTodo() tea.Msg {
	fromList := m.FocusedList()
	selectedTodo, ok := fromList.SelectedItem().(Todo)
	if !ok {
		return nil
	}

	selectedTodo.Toggle()
	fromList.RemoveItem(fromList.Index())

	toList := m.Lists[m.Focused.Switch()]
	cmd := toList.InsertItem(len(toList.Items()), selectedTodo)

	return cmd
}

func (m *Model) InitList(width, height int) {
	todos := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height/2)
	todos.SetShowHelp(false)
	todos.Title = "To Do"
	todos.SetItems([]list.Item{
		Todo{title: "Todo 1", description: "Important task to accomplish", done: false},
		Todo{title: "Todo 2", description: "Important task to accomplish", done: false},
		Todo{title: "Todo 3", description: "Important task to accomplish", done: false},
	})

	dones := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height/2)
	dones.SetShowHelp(false)
	dones.Title = "Done"
	dones.SetItems([]list.Item{
		Todo{title: "Todo 4", description: "Important task to accomplish", done: true},
		Todo{title: "Todo 4", description: "Important task to accomplish", done: true},
	})

	m.Lists = map[Status]*list.Model{
		TODO: &todos,
		DONE: &dones,
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		listStyle.Width(msg.Width / 3)
		focusedListStyle.Width(msg.Width / 3)
		if m.Loading {
			m.InitList(msg.Width, msg.Height)
			m.Loading = false
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit
		case "tab":
			m.Focused = m.Focused.Switch()
		case " ":
			return m, m.ToggleTodo
		}

	}

	updatedList, cmd := m.Lists[m.Focused].Update(msg)
	m.Lists[m.Focused] = &updatedList
	return m, cmd
}

func (m *Model) View() string {
	if len(m.Lists) != 0 {
		return lipgloss.JoinHorizontal(lipgloss.Center, m.RenderList(TODO), m.RenderList(DONE))
	}
	return ""
}

func NewModel() *Model {
	return &Model{}
}

func New() *tea.Program {
	model := NewModel()
	return tea.NewProgram(model, tea.WithAltScreen())
}
