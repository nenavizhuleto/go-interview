package api

import (
	"github.com/gofiber/fiber/v2"

	"helpdesk/internals/models/v3/device"
)

func SetDevicesRoutes(path string, router fiber.Router) {
	devices := router.Group(path)
	devices.Get("/", GetDevices)
	devices.Delete("/:id", DeleteDevice)
}

func GetDevices(c *fiber.Ctx) error {
	devices, err := device.All()
	if err != nil {
		return err
	}

	return c.JSON(devices)
}

func DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	device, err := device.Get(id)
	if err != nil {
		return err
	}

	if err := device.Delete(); err != nil {
		return err
	}

	return c.JSON(device)
}
