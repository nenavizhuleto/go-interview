package components

import (
	"projector/internals/applications/projector/common"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type List struct {
	common common.Common
	self   list.Model
}

// SetSize implements component
func (l List) SetSize(width, height int) {
	l.self.SetSize(width, height)
}

func NewList(common common.Common) List {
	items, _ := common.Context().Value("items").([]list.Item)
	ls := list.New(items, list.NewDefaultDelegate(), common.Width, common.Height)
	return List{
		common: common,
		self:   ls,
	}
}

// Init implements tea.Model
func (l List) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (l List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	l.self, cmd = l.self.Update(msg)
	return l, cmd
}

// View implements tea.Model
func (l List) View() string {
	return l.self.View()
}
