package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ErrorDTO struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func HandleError(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	return c.Status(code).JSON(ErrorDTO{
		Status:  code,
		Message: fmt.Sprintf("error: %s", err.Error()),
	})
}
