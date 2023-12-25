package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/domains/device"

)

func SetDevicesRoutes(path string, router fiber.Router) {
	devices := router.Group(path)

	devices.Get("/", GetDevices)
	devices.Post("/", CreateDevice)
}

func GetDevices(c *fiber.Ctx) error {
	devices, err := device.All()
	if err != nil {
		return err
	}

	return c.JSON(devices)
}

func CreateDevice(c *fiber.Ctx) error {
	var device device.Device
	if err := c.BodyParser(&device); err != nil {
		return fmt.Errorf("createDevice: %w", err)
	}

	if err := device.Save(); err != nil {
		return fmt.Errorf("createDevice: %w", err)
	}

	return c.JSON(device)
}
