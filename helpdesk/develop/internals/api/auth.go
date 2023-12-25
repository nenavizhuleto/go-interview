package api

import (
	"helpdesk/internals/models"
	"helpdesk/internals/models/device"
	"helpdesk/internals/models/user"

	"github.com/gofiber/fiber/v2"
)

func GetToken(c *fiber.Ctx) error {
	dev, err := device.Get(c.IP())
	if err != nil {
		return err
	}
	user, err := user.Get(dev.OwnerID)
	if err != nil {
		return err
	}

	return c.JSON(Success(
		fiber.Map{
			"token":         user.ID,
			"refresh_token": "refresh_token",
		},
	))
}

func Register(c *fiber.Ctx) error {
	var body struct {
		Name  string
		Phone string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("register", err)
	}

	user, err := user.New(body.Name, body.Phone, c.IP())
	if err != nil {
		return err
	}

	dev, err := device.New(c.IP(), device.PC)
	if err != nil {
		return err
	}

	if err := user.AddDevice(dev); err != nil {
		return err
	}

	if err := user.Save(); err != nil {
		return err
	}

	return c.JSON(Success(
		fiber.Map{
			"token":         user.ID,
			"refresh_token": "refresh_token",
		},
	))

}
