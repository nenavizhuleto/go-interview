package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/auth/v2"
	"helpdesk/internals/models/v2"
)

func GetIdentity(c *fiber.Ctx) error {
	dev, err := auth.GetIdentity(c.IP())
	if err != nil {
		return fmt.Errorf("identity: %w", err)
	}

	return c.JSON(dev)
}

func Register(c *fiber.Ctx) error {
	dev, err := models.NewDevice(c.IP())
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	if err := dev.Create(); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	var user = models.NewUser()

	if err := c.BodyParser(user); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	user.Devices = append(user.Devices, c.IP())
	user.OnAfterCreate = models.DeviceUserCreateHook
	if err := user.Create(); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	return c.JSON(user)
}
