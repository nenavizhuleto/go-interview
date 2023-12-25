package application

import (
	"projector/internals/applications/projector/common"
	personsPage "projector/internals/applications/projector/pages/persons"
	projectsPage "projector/internals/applications/projector/pages/projects"
	tasksPage "projector/internals/applications/projector/pages/tasks"

	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	projects page = iota
	tasks
	persons
	lastPage
)

type Model struct {
	common     common.Common
	activePage page
	pages      []common.Component
}

func NewModel() Model {
	c := common.New()
	return Model{
		common:     c,
		activePage: projects,
		pages:      make([]common.Component, 3),
	}
}

func (m *Model) NextPage() {
	m.activePage++
	if m.activePage == lastPage {
		m.activePage = projects
	}
}

// SetSize implements component
func (m *Model) SetSize(width, height int) {
	m.common.SetSize(width, height)
	for _, p := range m.pages {
		p.SetSize(width, height)
	}
}

// Init implements tea.Model interface
func (m *Model) Init() tea.Cmd {
	m.pages[projects] = projectsPage.New(m.common)
	m.pages[tasks] = tasksPage.New(m.common)
	m.pages[persons] = personsPage.New(m.common)
	var cmds []tea.Cmd
	for _, p := range m.pages {
		cmds = append(cmds, p.Init())
	}
	return tea.Batch(cmds...)
}

// Update implements tea.Model interface
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		for i, p := range m.pages {
			model, cmd := p.Update(msg)
			m.pages[i] = model.(common.Component)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.NextPage()
		}
	case projectsPage.ProjectSelectedMsg:
		// Open tasks page with selected project's tasks
		t, cmd := m.pages[tasks].Update(msg)
		m.pages[tasks] = t.(*tasksPage.Tasks)
		m.activePage = tasks
		cmds = append(cmds, cmd)
	case tasksPage.BackMsg:
		m.activePage = projects
	}

	p, cmd := m.pages[m.activePage].Update(msg)
	m.pages[m.activePage] = p.(common.Component)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.pages[m.activePage].View()
}

func New() *tea.Program {
	model := NewModel()
	program := tea.NewProgram(&model, tea.WithAltScreen())
	return program
}
