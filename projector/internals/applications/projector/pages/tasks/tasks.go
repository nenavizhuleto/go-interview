package tasks

import (
	"projector/internals/applications/projector/common"
	projectsPage "projector/internals/applications/projector/pages/projects"
	"projector/internals/core/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Tasks struct {
	common common.Common
	title  string
	tasks  list.Model
}

type BackMsg struct{}

func MockTasks() []list.Item {
	items := []list.Item{
		domain.NewTask("Task 1", "This is Task 1"),
		domain.NewTask("Task 2", "This is Task 2"),
		domain.NewTask("Task 3", "This is Task 3"),
	}
	return items
}

func New(common common.Common) *Tasks {
	return &Tasks{
		common: common,
		title:  "Tasks",
	}
}

// SetSize implements component
func (p *Tasks) SetSize(width, height int) {
	p.tasks.SetSize(width, height)
}

// Init implements tea.Model
func (p *Tasks) Init() tea.Cmd {
	p.tasks = list.New(MockTasks(), list.NewDefaultDelegate(), p.common.Width, p.common.Height)
	p.tasks.Title = p.title
	return nil
}

// Update implements tea.Model
func (p *Tasks) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "b":
			return p, Back
		}
	case projectsPage.ProjectSelectedMsg:
		// Open tasks page with selected project's tasks
		var items []list.Item
		for _, task := range msg.Tasks {
			items = append(items, task)
		}
		p.tasks.SetItems(items)
	}

	p.tasks, cmd = p.tasks.Update(msg)
	return p, cmd
}

func (p *Tasks) View() string {
	return p.tasks.View()
}

func Back() tea.Msg {
	return BackMsg{}
}
