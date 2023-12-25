package persons

import (
	"projector/internals/applications/projector/common"
	"projector/internals/core/domain"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Persons struct {
	common  common.Common
	title   string
	persons list.Model
}

func MockPersons() []list.Item {
	items := []list.Item{
		domain.NewPerson("John Doe", domain.Role("Manager")),
		domain.NewPerson("Marian Frog", domain.Role("SEO")),
		domain.NewPerson("P's Pad's", domain.Role("Developer")),
	}
	return items
}

func New(common common.Common) *Persons {
	return &Persons{
		common: common,
		title:  "Persons",
	}
}

// SetSize implements component
func (p *Persons) SetSize(width, height int) {
	p.persons.SetSize(width, height)
}

// Init implements tea.Model
func (p *Persons) Init() tea.Cmd {
	p.persons = list.New(MockPersons(), list.NewDefaultDelegate(), p.common.Width, p.common.Height)
	p.persons.Title = p.title
	return nil
}

// Update implements tea.Model
func (p *Persons) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	p.persons, cmd = p.persons.Update(msg)
	return p, cmd
}

func (p *Persons) View() string {
	return p.persons.View()
}
