package company

import (
	"github.com/google/uuid"
)

type CompanyID string

type Company struct {
	ID         CompanyID  `json:"id"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	Branches   []Branch   `json:"branches"`
	Properties []Property `json:"properties"`
}

func NewCompany(name, slug string) Company {
	return Company{
		ID:         CompanyID(uuid.NewString()),
		Name:       name,
		Slug:       slug,
		Branches:   make([]Branch, 0),
		Properties: make([]Property, 0),
	}
}

func (c *Company) SetName(name string) {
	c.Name = name
}

func (c *Company) AddBranch(b Branch) BranchID {
	for _, prop := range c.Properties {
		b.AddProperty(prop)
	}
	c.Branches = append(c.Branches, b)
	return b.ID
}

func (c *Company) AddProperty(property Property) {
	for _, b := range c.Branches {
		b.AddProperty(property)
	}
	c.Properties = append(c.Properties, property)
}

func (c Company) GetBranch(id BranchID) Branch {
	var branch Branch
	for _, b := range c.Branches {
		if b.ID == id {
			branch = b
			return branch
		}
	}

	return branch
}
