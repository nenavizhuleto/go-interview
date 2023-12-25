package api

import (
	"encoding/json"
	"helpdesk/internals/megaplan"
	"helpdesk/internals/models/v3/task"
	"log"

	"github.com/gofiber/fiber/v2"
)

func HandleMegaplanEvent(c *fiber.Ctx) error {
	var event megaplan.TaskEvent
	if err := c.BodyParser(&event); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	if event.Model != "Task" {
		return c.SendStatus(200)
	}

	str, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("Event: %s", string(str))

	dto := event.Data
	task, err := task.Get(dto.ID)
	if err != nil {
		log.Printf("event: %s", err.Error())
		return c.SendStatus(200)
	}

	task.Status = dto.GetStatus()
	task.Comments = dto.GetComments()
	if dto.Activity != nil {
		task.LastActivity = dto.Activity.Value
	}

	if err := task.Save(); err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}
