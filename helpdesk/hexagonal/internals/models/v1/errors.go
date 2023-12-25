package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrUserNotFound          = fiber.NewError(fiber.StatusServiceUnavailable, "system: user not found for device")
	ErrUnsupportedDevice     = fiber.NewError(fiber.StatusServiceUnavailable, "system: device's ip not found in subnets")
	ErrNotIdentified         = fiber.NewError(fiber.StatusUnauthorized, "user not identified")
	ErrAPIInvalidRequestBody = fiber.NewError(fiber.StatusNotAcceptable, "request body to api is invalid")
	ErrEntityNotFound        = fiber.NewError(fiber.StatusNotFound, "entities not found")
	ErrMegaplan              = fiber.NewError(fiber.StatusServiceUnavailable, "megaplan api error")
	ErrDatabase              = fiber.NewError(fiber.StatusInternalServerError, "database error")
)

func ValidationError(message string) error {
	return errors.New(message)
}
