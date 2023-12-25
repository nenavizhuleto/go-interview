package company

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/repo"
	"helpdesk/internals/core/utils"
)

type CompanyRepo interface {
	core.EntityRepo[Company]
}

type Plugin struct {
	Repo   CompanyRepo
	logger core.Logger
	debug  bool
}

func New() *Plugin {
	r := repo.NewMemoryRepo[Company]("company")
	return &Plugin{
		Repo:  &r,
		debug: false,
	}
}

func (p *Plugin) SetDebug(v bool) {
	p.debug = v
}

func (p *Plugin) SetLogger(l core.Logger) {
	l.SetDebug(p.debug)
	p.logger = l
}

func (p *Plugin) GetCompanies() ([]Company, error) {
	companies, err := p.Repo.All()
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (p *Plugin) CreateCompany(name string) (Company, error) {
	slug := utils.ToSlug(name)
	company := NewCompany(name, slug)
	if err := p.Repo.Save(company); err != nil {
		return Company{}, err
	}

	return company, nil
}

func (p *Plugin) GetCompany(slug string) (Company, error) {
	company, err := p.Repo.Filter(func(c Company) bool {
		return c.Slug == slug
	})
	if err != nil {
		return Company{}, err
	}

	return company, nil
}

type Employee struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (p *Plugin) AddBranchToCompany(slug string, name, description, address, contacts string) (core.EntityID, error) {
	company, err := p.GetCompany(slug)
	if err != nil {
		return core.NilEntityID(), err
	}

	branch_id := company.AddBranch(name,
		description,
		address,
		contacts,
	)

	if err := p.Repo.Save(company); err != nil {
		return core.NilEntityID(), err
	}

	return branch_id, nil
}

func (p *Plugin) Setup(c *core.Core) error {
	p.logger.Printf("setting up")
	return nil
}

func (p *Plugin) Run(c *core.Core) error {
	p.logger.Printf("starting")
	return nil
}
