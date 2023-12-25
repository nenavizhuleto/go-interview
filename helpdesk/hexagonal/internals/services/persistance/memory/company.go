package memory

import (
	"errors"
	"helpdesk/internals/core/company"
)

var (
	ErrCompanyNotFound = errors.New("company not found")
)

type CompanyRepo struct {
	name string
	c    map[company.CompanyID]company.Company
}

func NewCompanyRepo(name string) CompanyRepo {
	return CompanyRepo{
		name: name,
		c:    make(map[company.CompanyID]company.Company),
	}
}

func (r CompanyRepo) Get(id company.CompanyID) (company.Company, error) {
	if c, ok := r.c[id]; ok {
		return c, nil
	} else {
		return company.Company{}, ErrCompanyNotFound
	}
}

func (r CompanyRepo) All() ([]company.Company, error) {
	companies := make([]company.Company, 0)
	for _, c := range r.c {
		companies = append(companies, c)
	}

	return companies, nil
}

func (r *CompanyRepo) Save(c company.Company) error {
	r.c[c.ID] = c
	return nil
}

func (r *CompanyRepo) Delete(id company.CompanyID) error {
	if _, ok := r.c[id]; ok {
		delete(r.c, id)
		return nil
	} else {
		return ErrCompanyNotFound
	}
}

func (r *CompanyRepo) GetBySlug(slug string) (company.Company, error) {
	for _, c := range r.c {
		if c.Slug == slug {
			return c, nil
		}
	}

	return company.Company{}, ErrCompanyNotFound
}
