package api

import "github.com/gofiber/fiber/v2"

func HandleError(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	return c.Status(code).JSON(Fail(err))
}
