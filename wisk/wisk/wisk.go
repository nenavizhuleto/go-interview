package wisk

import "os/exec"

type Command struct {
	Name   string
	output string
}

func (c *Command) Exec() error {
	cmd := exec.Command("ls")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	c.output = string(stdout)
	return nil
}

func (c Command) Output() string {
	return c.output
}
