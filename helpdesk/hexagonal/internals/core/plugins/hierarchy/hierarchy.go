package hierarchy

import (
	"helpdesk/internals/core"
)

type Plugin struct {
	Root   *Node
	logger core.Logger
	debug  bool
}

func New() *Plugin {

	return &Plugin{
		Root: NewNode("root"),
	}
}

func (p *Plugin) SetDebug(v bool) {
	p.debug = v
}

func (p *Plugin) SetLogger(l core.Logger) {
	l.SetDebug(p.debug)
	p.logger = l
}

func (p *Plugin) Setup(c *core.Core) error {
	p.logger.Printf("setting up")
	n := NewNode("devices").AddChild(
		NewNode("172.16.222.31"),
	)

	p.logger.Printf("Setup: %#v", n)
	p.Root.AddChild(
		NewNode("company 1").AddChild(
			NewNode("departments").AddChild(
				NewNode("depart 1").AddChild(
					NewNode("employees").AddChild(
						NewNode("employee_1"),
						NewNode("employee_2"),
					),
					NewNode("manager"),
				),
				NewNode("depart 2").AddChild(
					NewNode("employees").AddChild(
						NewNode("employee_3").AddChild(n),
						NewNode("employee_4"),
					),
					NewNode("manager"),
				),
			),
			NewNode("director"),
		),
		NewNode("company 2").AddChild(
			NewNode("departments").AddChild(
				NewNode("depart 1").AddChild(
					NewNode("employees").AddChild(
						NewNode("employee_1"),
						NewNode("employee_2"),
					),
					NewNode("manager"),
				),
			),
		),
	)

	return nil
}

func (p *Plugin) Run(c *core.Core) error {
	p.logger.Printf("starting")
	p.logger.Printf("Tree \n%s\n", PrintHierarchy(p.Root, 0))

	n := FindByID_DFS(p.Root, "employee_3")
	p.logger.Printf("Search: %v", n)
	return nil
}
