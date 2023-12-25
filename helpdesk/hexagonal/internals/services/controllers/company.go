package controllers

import (
	"errors"
	"helpdesk/internals/core"
	"helpdesk/internals/core/company"
	"helpdesk/internals/services/persistance/memory"
	"strings"
)

var (
	ErrCompanyAlreadyRegistered = errors.New("company already registered")
	ErrBranchNotFound           = errors.New("branch not found")
)

type CompanyRepo interface {
	core.Repo[company.CompanyID, company.Company]
	GetBySlug(string) (company.Company, error)
}

type CompanyController struct {
	Repo CompanyRepo
}

func NewCompanyController(repo CompanyRepo) CompanyController {
	return CompanyController{
		Repo: repo,
	}
}

func NewMemoryCompanyController() CompanyController {
	repo := memory.NewCompanyRepo("companies")
	return CompanyController{
		Repo: &repo,
	}
}

func (cc *CompanyController) RegisterCompany(name string) (company.CompanyID, error) {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	if c, err := cc.Repo.GetBySlug(slug); err == nil {
		return c.ID, ErrCompanyAlreadyRegistered
	}
	c := company.NewCompany(name, slug)
	if err := cc.Repo.Save(c); err != nil {
		return "", err
	}
	return c.ID, nil
}

func (cc *CompanyController) AddBranchToCompany(
	companyID company.CompanyID,
	b struct {
		Name        string
		Description string
		Address     string
		Contacts    string
	},
) (company.BranchID, error) {
	c, err := cc.Repo.Get(companyID)
	if err != nil {
		return "", err
	}
	branchID := c.AddBranch(b.Name, b.Description, b.Address, b.Contacts)
	return branchID, cc.Repo.Save(c)
}
