package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type LizingReporter struct {
	title  string
	app    fyne.App
	master fyne.Window
}

func NewLizingReporter() *LizingReporter {
	l := &LizingReporter{
		title: "Lizing Reporter",
		app:   app.New(),
	}

	l.master = l.app.NewWindow(l.title)


	top := container.NewVBox()
	bottom := container.NewStack()

	main := container.NewVSplit(top, bottom)
	debug := container.NewStack(debugLog)

	tabs := container.NewAppTabs(
		container.NewTabItem("Main", main),
		container.NewTabItem("Debug", debug),
	)

	l.master.SetContent(tabs)
	return l
}

func (l *LizingReporter) Run() {
	l.master.Show()
	l.app.Run()

}
