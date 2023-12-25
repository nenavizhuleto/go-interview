package core

import (
	"context"
	"helpdesk/internals/core/utils"
	"sync"
)

type Core struct {
	ctx     context.Context
	plugins []Plugin
}

func New() Core {
	return Core{
		ctx:     context.Background(),
		plugins: make([]Plugin, 0),
	}
}

func (c *Core) SetContext(key interface{}, value interface{}) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

func (c *Core) GetContext(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *Core) Use(name string, p Plugin) {
	if pl, ok := p.(PluginLoggable); ok {
		pl.SetLogger(NewDefaultPluginLogger(utils.ToUpper(name)))
	}
	p.Setup(c)

	c.plugins = append(c.plugins, p)
}

func (c *Core) Run() {
	var wg sync.WaitGroup
	for _, p := range c.plugins {
		wg.Add(1)
		p := p
		go func() {
			defer wg.Done()
			// WARNING: Capturing variable in loop
			p.Run(c)
		}()
	}

	wg.Wait()
}
