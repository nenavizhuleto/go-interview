package form

import (
	"projector/internals/applications/projector/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	common     common.Common
	focusIndex int
	inputs     []textinput.Model
	show       bool
}

func New(c common.Common) *Model {
	return &Model{
		common: c,
		inputs: make([]textinput.Model, 0),
	}
}

func (m *Model) AddInput(placeholder string) {
	ti := textinput.New()
	ti.Placeholder = placeholder

	m.inputs = append(m.inputs, ti)
}

func (m *Model) Focus() tea.Cmd {
	if len(m.inputs) > 0 {
		return m.inputs[0].Focus()
	}
	return nil
}

func (m *Model) Show() {
	m.show = true
}

func (m Model) Active() bool {
	return m.show
}

func (m *Model) Hide() {
	m.show = false
}

// SetSize implements component
func (m *Model) SetSize(width, height int) {
}

// Init implements tea.Model
func (m *Model) Init() tea.Cmd {
	return nil
}

type SubmitFormMsg struct {
	Values []string
}

func (m *Model) SubmitForm() tea.Msg {
	msg := SubmitFormMsg{
		Values: make([]string, len(m.inputs)),
	}
	for _, i := range m.inputs {
		msg.Values = append(msg.Values, i.Value())
	}

	m.Hide()

	return msg
}

// Update implements tea.Model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down", "enter":
			s := msg.String()

			if s == "enter" {
				return m, m.SubmitForm
			}

			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			for i := range m.inputs {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					continue
				}

				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.show {
		var inputs []string
		for _, i := range m.inputs {
			inputs = append(inputs, i.View())
		}

		return lipgloss.JoinVertical(lipgloss.Top, inputs...)
	}

	return ""
}
