package api

import (
	"helpdesk/internals/util"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleCompanyRoutes() {
	// All
	s.router.Get("/company", func(c *fiber.Ctx) error {
		cs, err := s.CompanyController.Repo.All()
		if err != nil {
			return err
		}

		return c.JSON(cs)
	})

	// One
	s.router.Get("/company/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		com, err := s.CompanyController.Repo.GetBySlug(slug)
		if err != nil {
			return err
		}

		return c.JSON(com)
	})

	// Branches
	s.router.Get("/company/:slug/branch", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		com, err := s.CompanyController.Repo.GetBySlug(slug)
		if err != nil {
			return err
		}

		return c.JSON(com.Branches)
	})

	// Create One
	s.router.Post("/company", func(c *fiber.Ctx) error {
		var body struct {
			Name string
		}

		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if id, err := s.RegisterCompany(body.Name); err != nil {
			return err
		} else {
			return c.JSON(util.JSON{"id": id})
		}

	})

	// Create branch
	s.router.Post("/company/:slug/branch", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		com, err := s.CompanyController.Repo.GetBySlug(slug)
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
			return err
		}

		if b_id, err := s.AddBranchToCompany(com.ID, body); err != nil {
			return err
		} else {
			return c.JSON(util.JSON{"id": b_id})
		}
	})

}
