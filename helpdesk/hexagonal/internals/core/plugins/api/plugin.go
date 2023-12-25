package api

import (
	"helpdesk/internals/core"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type ExternAPI interface {
	PrefixAPI() string
	IntoAPI(fiber.Router)
}

type ExternApp interface {
	fiber.Router
}

type Plugin struct {
	debug   bool
	logger  core.Logger
	app     *fiber.App
	externs []ExternAPI
}

func New() *Plugin {
	return &Plugin{
		debug:   false,
		app:     fiber.New(),
		externs: make([]ExternAPI, 0),
	}
}

func (p *Plugin) SetDebug(v bool) {
	p.debug = v
}

func (p *Plugin) SetLogger(l core.Logger) {
	l.SetDebug(p.debug)
	p.logger = l
}

func (p *Plugin) Use(api ExternAPI) {
	p.externs = append(p.externs, api)
}

func (p *Plugin) Setup(c *core.Core) error {
	p.logger.Printf("setting up")
	p.app.Use(logger.New())
	p.app.Use(recover.New())
	p.app.Use(cors.New())

	// Initializing external API's
	for _, ea := range p.externs {
		ea.IntoAPI(p.app.Group(ea.PrefixAPI()))
	}
	return nil
}

func (p *Plugin) Run(c *core.Core) error {
	p.logger.Printf("starting")

	p.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	return p.app.Listen(":3000")
}
