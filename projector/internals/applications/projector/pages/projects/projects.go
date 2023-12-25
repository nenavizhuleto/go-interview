package projects

import (
	"fmt"
	"projector/internals/applications/projector/common"
	"projector/internals/applications/projector/components/form"
	"projector/internals/core/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Projects struct {
	common   common.Common
	title    string
	projects list.Model
	form     *form.Model
}

func MockProjects() []list.Item {
	t1 := domain.NewTask("Task 1", "This is Task 1")
	t2 := domain.NewTask("Task 2", "This is Task 2")
	t3 := domain.NewTask("Task 3", "This is Task 3")
	proj1 := domain.NewProject("Proj 1", "This is Proj 1")
	proj1.AddTask(t1)
	proj1.AddTask(t2)
	proj1.AddTask(t3)
	items := []list.Item{
		proj1,
		domain.NewProject("Proj 2", "This is Proj 2"),
		domain.NewProject("Proj 3", "This is Proj 3"),
	}
	return items
}

type ProjectSelectedMsg struct {
	Tasks map[string]*domain.Task
}

func New(common common.Common) *Projects {
	return &Projects{
		common: common,
		title:  "Projects",
	}
}

// SetSize implements component
func (p *Projects) SetSize(width, height int) {
	fmt.Println(width, height)
	p.projects.SetSize(width, height)
}

// Init implements tea.Model
func (p *Projects) Init() tea.Cmd {
	p.projects = list.New(MockProjects(), list.NewDefaultDelegate(), p.common.Width, p.common.Height)
	p.projects.Title = p.title

	p.form = form.New(p.common)
	p.form.AddInput("Project Name")
	p.form.AddInput("Project Description")
	return tea.Batch(
		p.form.Focus(),
	)
}

// Update implements tea.Model
func (p *Projects) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	if p.form.Active() {
		f, cmd := p.form.Update(msg)
		p.form = f.(*form.Model)
		cmds = append(cmds, cmd)
		return p, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			project := p.projects.SelectedItem().(*domain.Project)
			return p, SelectProject(project)
		}
		switch msg.String() {
		case "n":
			p.form.Show()
		}
	case form.SubmitFormMsg:
		project := domain.NewProject(msg.Values[0], msg.Values[1])
		p.projects.InsertItem(len(p.projects.Items()), project)
	}

	p.projects, cmd = p.projects.Update(msg)
	cmds = append(cmds, cmd)
	return p, tea.Batch(cmds...)
}

func (p *Projects) View() string {

	return lipgloss.JoinHorizontal(lipgloss.Left, p.projects.View(), p.form.View())
}

// Commands
func SelectProject(project *domain.Project) tea.Cmd {
	return func() tea.Msg {
		return ProjectSelectedMsg{Tasks: project.Tasks}
	}
}
