package api

import (
	"helpdesk/internals/models"
	"helpdesk/internals/models/user"

	"github.com/gofiber/fiber/v2"
)

const UserTokenHeader = "token"

func UserMiddleware(c *fiber.Ctx) error {
	token := c.Get(UserTokenHeader)

	u, err := user.Get(token)
	if err != nil {
		return models.NewTokenError(token, err)
	}

	c.Locals("user", *u)

	return c.Next()
}
