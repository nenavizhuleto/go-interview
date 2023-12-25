package common

import "context"

type Common struct {
	ctx    context.Context
	Width  int
	Height int
}

func New() Common {
	return Common{
		ctx: context.TODO(),
	}
}

// SetValue sets value in the context
func (c *Common) SetValue(key, value interface{}) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

// SetSize sets the width and the height
func (c *Common) SetSize(width, height int) {
	c.Width = width
	c.Height = height
}

// Context return context
func (c *Common) Context() context.Context {
	return c.ctx
}
