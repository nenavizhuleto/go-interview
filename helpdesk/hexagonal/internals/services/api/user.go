package api

import (
	"helpdesk/internals/core/company"
	"helpdesk/internals/services/controllers"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleUserRoutes() {

	s.router.Get("/user", func(c *fiber.Ctx) error {
		users, err := s.GetUsers()
		if err != nil {
			return err
		}

		return c.JSON(users)
	})

	s.router.Post("/user", func(c *fiber.Ctx) error {
		var body struct {
			Name      string
			Phone     string
			CompanyID string `json:"company_id"`
			BranchID  string `json:"branch_id"`
		}
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		log.Printf("Body: %#v\n", body)

		u, err := s.CreateUser(body.Name, body.Phone)
		if err != nil {
			return err
		}

		if body.CompanyID != "" {
			c, err := s.CompanyController.Repo.Get(company.CompanyID(body.CompanyID))
			if err != nil {
				return err
			}
			u.SetCompany(c.ID)
			if body.BranchID != "" {
				var branch company.Branch
				for _, b := range c.Branches {
					if b.ID == company.BranchID(body.BranchID) {
						branch = b
					}
				}

				if branch.ID == "" {
					return controllers.ErrBranchNotFound
				}

				u.SetBranch(branch.ID)
			}

			if err := s.UserController.Repo.Save(u); err != nil {
				return err
			}
		}

		return c.JSON(u)
	})

}
