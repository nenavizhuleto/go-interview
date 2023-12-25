package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/models/v1"
)

type Event struct {
	Data  models.Task `json:"data"`
	Event string      `json:"event"`
	Model string      `json:"model"`
}
type Task struct {
	ID     string `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Status string `json:"status" db:"status"`
}

func HandleMegaplanEvent(c *fiber.Ctx) error {
	var event Event
	if err := c.BodyParser(&event); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	if event.Model != "Task" {
		return c.SendStatus(200)
	}

	log.Printf("Event: %#v", event)

	task := event.Data
	if err := models.UpdateTaskStatus(task.ID, task.Status); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}
