package api

import (
	"helpdesk/internals/services/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Service struct {
	router fiber.Router
	port   string
	controllers.CompanyController
	controllers.UserController
	controllers.TaskController
	companies controllers.CompanyRepo
}

func New(port string, memory bool) Service {

	var c controllers.CompanyController
	var u controllers.UserController
	var t controllers.TaskController
	if memory {
		c = controllers.NewMemoryCompanyController()
		u = controllers.NewMemoryUserController()
		t = controllers.NewMemoryTaskController()
	}

	return Service{
		port:              port,
		CompanyController: c,
		UserController:    u,
		TaskController:    t,
	}
}



func (s Service) Run() error {
	// initializing fiber application
	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	router := app.Group("/api")

	s.router = router

	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("api.")
	})

	s.HandleCompanyRoutes()
	s.HandleTaskRoutes()
	s.HandleUserRoutes()

	return app.Listen(":" + s.port)
}
