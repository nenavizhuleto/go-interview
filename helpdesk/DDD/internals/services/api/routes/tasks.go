package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"helpdesk/internals/domains/task"
)

func SetTasksRoutes(path string, router fiber.Router) {
	tasks := router.Group(path)
	tasks.Get("/", GetTasks)
	tasks.Post("/", CreateTask)
	tasks.Put("/", UpdateTask)
	tasks.Delete("/", DeleteTask)
}

func GetTasks(c *fiber.Ctx) error {
	tasks, err := task.All()
	if err != nil {
		return fmt.Errorf("getTasks: %w", err)
	}

	return c.JSON(tasks)
}

func CreateTask(c *fiber.Ctx) error {
	var task task.Task
	if err := c.BodyParser(&task); err != nil {
		return fmt.Errorf("createTask: %w", err)
	}


	if err := task.Save(); err != nil {
		return fmt.Errorf("createTask: %w", err)
	}

	return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
	var task task.Task
	if err := c.BodyParser(&task); err != nil {
		return fmt.Errorf("updateTask: %w", err)
	}

	if err := task.Save(); err != nil {
		return fmt.Errorf("updateTask: %w", err)
	}

	return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	var task task.Task
	if err := c.BodyParser(&task); err != nil {
		return fmt.Errorf("deleteTask: %w", err)
	}

	if err := task.Delete(); err != nil {
		return fmt.Errorf("deleteTask: %w", err)
	}

	return c.JSON(task)
}
