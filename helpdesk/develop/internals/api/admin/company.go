package admin

import (
	"helpdesk/internals/api"
	"helpdesk/internals/models"
	"helpdesk/internals/models/company"

	"github.com/gofiber/fiber/v2"
)

// GET /company
func GetCompanies(c *fiber.Ctx) error {
	companies, err := company.All()
	if err != nil {
		return err
	}

	return c.JSON(api.Success(companies))
}

// GET /company/:id
func GetCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	_company, err := company.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(api.Success(_company))
}

// POST /company
func NewCompany(c *fiber.Ctx) error {
	var body struct {
		Itn  string `json:"itn"`
		Name string `json:"name"`
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("company", err)
	}

	_company, err := company.New(body.Itn, body.Name)
	if err != nil {
		return err
	}

	if err := _company.Save(); err != nil {
		return err
	}

	return c.JSON(api.Success(_company))
}

// DELETE /company/:id
func DeleteCompany(c *fiber.Ctx) error {
	id := c.Params("id")

	_company, err := company.Get(id)
	if err != nil {
		return err
	}

	if err := _company.Delete(); err != nil {
		return err
	}

	return c.JSON(api.Success(_company))
}
