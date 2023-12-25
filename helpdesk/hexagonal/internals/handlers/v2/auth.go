package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"helpdesk/internals/auth/v2"
	"helpdesk/internals/models/v2"
)

func HandleRegister(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	user.ID = uuid.NewString()

	device, err := auth.GetIdentity(c.IP())
	fmt.Print(err)
	if err != nil {

		device, err := auth.Register(c.IP(), &user)
		if err != nil {
			return fmt.Errorf("register: %w", err)
		}

		return c.JSON(device)
	}

	return c.JSON(device)
}
