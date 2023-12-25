package admin

import (
	"helpdesk/internals/api"
	"helpdesk/internals/models"
	"helpdesk/internals/models/branch"
	"helpdesk/internals/models/company"

	"github.com/gofiber/fiber/v2"
)

// GET /branch
func GetBranches(c *fiber.Ctx) error {
	branches, err := branch.All()
	if err != nil {
		return err
	}

	return c.JSON(api.Success(branches))
}

// GET /branch/:id
func GetBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	_branch, err := branch.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(api.Success(_branch))
}

// POST /branch/:company_id
func NewBranch(c *fiber.Ctx) error {
	company_id := c.Params("company_id")
	_company, err := company.Get(company_id)
	if err != nil {
		return err
	}

	var body struct {
		Name        string
		Description string
		Address     string
		Contacts    string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("branch", err)
	}

	_branch, err := branch.New(_company, body.Name, body.Description, body.Address, body.Contacts)
	if err != nil {
		return err
	}

	if err := _branch.Save(); err != nil {
		return err
	}

	return c.JSON(api.Success(_branch))
}

// DELETE /branch/:id
func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")

	_branch, err := branch.Get(id)
	if err != nil {
		return err
	}

	if err := _branch.Delete(); err != nil {
		return err
	}

	return c.JSON(api.Success(_branch))
}
