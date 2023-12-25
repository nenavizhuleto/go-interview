package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/domains/device"
	"helpdesk/internals/domains/user"
)

func GetIdentity(c *fiber.Ctx) error {
	dev, err := device.Get(c.IP())
	if err != nil {
		return err
	}

	return c.JSON(dev)
}

func Register(c *fiber.Ctx) error {
	if err := device.Identify(c.IP()); err != nil {
		return err
	}

	var body struct {
		Name string
		Phone string
	}

	if err := c.BodyParser(&body); err != nil {
		return fmt.Errorf("register: %w", err)
	}

	user, err := user.New(body.Name, body.Phone)
	if err != nil {
		return err
	}

	user.Devices = append(user.Devices, c.IP())
	user.OnAfterCreate = device.DeviceUserCreateHook
	if err := user.Save(); err != nil {
		return err
	}

	return c.JSON(user)
}
