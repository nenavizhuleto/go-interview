package admin

import (
	"helpdesk/internals/models"

	"github.com/gofiber/fiber/v2"
)

const AdminTokenHeader = "token"
const AdminTokenValue = "password"

func AdminMiddleware(c *fiber.Ctx) error {

	token := c.Get(AdminTokenHeader)

	if token != AdminTokenValue {
		return models.NewTokenError(token, nil)
	}

	return c.Next()
}
