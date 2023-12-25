package api

import (
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/company"

	"github.com/gofiber/fiber/v2"
)

func SetCompaniesRoutes(path string, router fiber.Router) {
	companies := router.Group(path)

	companies.Get("/", GetCompanies)
	companies.Get("/:id", GetCompany)
	companies.Post("/", CreateCompany)
	companies.Put("/:id", UpdateCompany)
	companies.Delete("/:id", DeleteCompany)
}

func GetCompanies(c *fiber.Ctx) error {
	companies, err := company.All()
	if err != nil {
		return err
	}

	return c.JSON(companies)
}

func GetCompany(c *fiber.Ctx) error {
	id := c.Params("id")
	company, err := company.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(company)
}

func CreateCompany(c *fiber.Ctx) error {
	var body struct {
		ITN  string
		Name string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("create_company", err)
	}

	company, err := company.New(body.ITN, body.Name)
	if err != nil {
		return err
	}

	if err := company.Save(); err != nil {
		return err
	}

	return c.JSON(company)
}

func UpdateCompany(c *fiber.Ctx) error {
	id := c.Params("id")

	company, err := company.Get(id)
	if err != nil {
		return err
	}

	var body struct {
		ITN  string
		Name string
	}

	if err := c.BodyParser(&body); err != nil {
		return models.NewParseError("update_company", err)
	}

	company.ID = body.ITN
	company.Name = body.Name

	if err := company.Save(); err != nil {
		return err
	}

	return c.JSON(company)
}

func DeleteCompany(c *fiber.Ctx) error {
	id := c.Params("id")

	company, err := company.Get(id)
	if err != nil {
		return err
	}

	if err := company.Delete(); err != nil {
		return err
	}

	return c.JSON(company)
}
