package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/models/v2"
)

func SetDevicesRoutes(path string, router fiber.Router) {
	devices := router.Group(path)

	devices.Get("/", GetDevices)
	devices.Post("/", CreateDevice)
}

func GetDevices(c *fiber.Ctx) error {
	devices, err := models.DevicesAll()
	if err != nil {
		return err
	}

	return c.JSON(devices)
}

func CreateDevice(c *fiber.Ctx) error {
	var device models.Device
	if err := c.BodyParser(&device); err != nil {
		return fmt.Errorf("createDevice: %w", err)
	}

	if err := device.Create(); err != nil {
		return fmt.Errorf("createDevice: %w", err)
	}

	return c.JSON(device)
}
