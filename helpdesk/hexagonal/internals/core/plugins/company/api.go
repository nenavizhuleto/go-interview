package company

import (
	"github.com/gofiber/fiber/v2"
)

func (p *Plugin) PrefixAPI() string {
	return "/company"
}

func (p *Plugin) IntoAPI(r fiber.Router) {
	r.Get("/", func(c *fiber.Ctx) error {
		companies, err := p.GetCompanies()
		if err != nil {
			return err
		}

		return c.JSON(companies)
	})

	r.Post("/", func(c *fiber.Ctx) error {
		var body struct {
			Name string
		}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		company, err := p.CreateCompany(body.Name)
		if err != nil {
			return err
		}

		return c.JSON(company)
	})

	r.Get("/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		company, err := p.GetCompany(slug)
		if err != nil {
			return err
		}

		return c.JSON(company)
	})

	r.Post("/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		var body struct {
			Name        string
			Description string
			Address     string
			Contacts    string
		}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		branch_id, err := p.AddBranchToCompany(slug, body.Name, body.Description, body.Address, body.Contacts)
		if err != nil {
			return err
		}

		return c.JSON(branch_id)
	})
}
