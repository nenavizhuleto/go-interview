package api

import (
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/branch"
	"helpdesk/internals/models/v3/company"

	"github.com/gofiber/fiber/v2"
)

func SetBranchesRoutes(path string, router fiber.Router) {
	branches := router.Group(path)

	branches.Get("/", GetBranches)
	branches.Get("/:id", GetBranch)
	branches.Post("/", CreateBranch)
	branches.Delete("/:id", DeleteBranch)
}

func GetBranches(c *fiber.Ctx) error {
	branches, err := branch.All()
	if err != nil {
		return err
	}

	return c.JSON(branches)
}

func GetBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	branch, err := branch.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(branch)
}

func CreateBranch(c *fiber.Ctx) error {
	var body struct {
		CompanyID   string `json:"company_id"`
		Name        string
		Description string
		Address     string
		Contacts    string
	}
	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("create_branch", err)
	}

	company, err := company.Get(body.CompanyID)
	if err != nil {
		return err
	}

	branch, err := branch.New(company, body.Name, body.Description, body.Address, body.Contacts)
	if err != nil {
		return err
	}

	if err := branch.Save(); err != nil {
		return err
	}

	return c.JSON(branch)
}

func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	branch, err := branch.Get(id)
	if err != nil {
		return err
	}

	if err := branch.Delete(); err != nil {
		return err
	}

	return c.JSON(branch)
}
