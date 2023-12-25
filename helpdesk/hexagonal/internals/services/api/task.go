package api

import (
	"helpdesk/internals/core/user"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleTaskRoutes() {

	s.router.Get("/task", func(c *fiber.Ctx) error {
		tasks, err := s.GetTasks()
		if err != nil {
			return err
		}

		return c.JSON(tasks)
	})

	s.router.Post("/task", func(c *fiber.Ctx) error {
		var body struct {
			Name    string
			Subject string
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		t, err := s.CreateTask(body.Name, body.Subject)
		if err != nil {
			return err
		}

		return c.JSON(t)
	})

	s.router.Post("/user/:id/task", func(c *fiber.Ctx) error {
		user_id := c.Params("id")

		u, err := s.UserController.Repo.Get(user.UserID(user_id))
		if err != nil {
			return err
		}

		var body struct {
			Name    string
			Subject string
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		t, err := s.CreateTask(body.Name, body.Subject)
		if err != nil {
			return err
		}

		if err := s.TaskController.SetTaskOwner(t, u); err != nil {
			return err
		}

		return c.JSON(t)
	})

}
