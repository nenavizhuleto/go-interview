package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/django/v2"

	"golos/handlers"
)

func main() {
	engine := django.New("views", ".html")
	engine.Reload(true)

	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	// app.Static("/", "./public")
	app.Use(logger.New())
	log.Printf("%v", engine.Templates)

	// app.Get("/", handlers.HandleIndex)
	app.Get("/", handlers.HandleMain)
	app.Get("/info", handlers.HandleGetInfo)
	app.Get("/devices", handlers.HandleGetDevices)
	app.Get("/rooms", handlers.HandleGetRooms)
	app.Get("/porches", handlers.HandleGetPorches)
	app.Get("/buildings", handlers.HandleGetBuildings)
	app.Get("/blocks", handlers.HandleGetBlocks)
	app.Get("/streets", handlers.HandleGetStreets)
	app.Get("/config", handlers.HandleGetConfig)

	log.Fatal(app.Listen(":3000"))
}
