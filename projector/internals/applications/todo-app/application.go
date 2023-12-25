package application

import tea "github.com/charmbracelet/bubbletea"

var app *App

func New() *tea.Program {

	app = NewApp()
	app.initLists()
	return tea.NewProgram(app)
}
