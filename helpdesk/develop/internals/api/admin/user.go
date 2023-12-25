package admin

import (
	"helpdesk/internals/api"
	"helpdesk/internals/models"
	"helpdesk/internals/models/user"

	"github.com/gofiber/fiber/v2"
)

// GET /user
func GetUsers(c *fiber.Ctx) error {
	users, err := user.All()
	if err != nil {
		return err
	}

	return c.JSON(api.Success(users))
}

// GET /user/:id
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_user, err := user.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(api.Success(_user))
}

// POST /user
func NewUser(c *fiber.Ctx) error {
	var body struct {
		IP    string `json:"ip"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("user", err)
	}

	_user, err := user.New(body.Name, body.Phone, body.IP)
	if err != nil {
		return err
	}

	if err := _user.Save(); err != nil {
		return err
	}

	return c.JSON(api.Success(_user))
}

// DELETE /user/:id
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_user, err := user.Get(id)
	if err != nil {
		return err
	}

	if err := _user.Delete(); err != nil {
		return err
	}

	return c.JSON(api.Success(_user))
}
