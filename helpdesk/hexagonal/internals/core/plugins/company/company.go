package company

import (
	"helpdesk/internals/core"

	"github.com/google/uuid"
)

type Branches core.Children[Branch]

type Company struct {
	ID       core.EntityID `json:"id"`
	Name     string        `json:"name"`
	Slug     string        `json:"slug"`
	Branches Branches      `json:"branches"`
}

// Implement Entity interface
func (c Company) EntityID() core.EntityID {
	return c.ID
}

// Implementing IParent interface
func (c Company) Children() Branches {
	return c.Branches
}

func NewCompany(name, slug string) Company {
	return Company{
		ID:       core.EntityID(uuid.NewString()),
		Name:     name,
		Slug:     slug,
		Branches: make(Branches),
	}
}

func (c *Company) SetName(name string) {
	c.Name = name
}

func (c *Company) AddBranch(name, description, address, contacts string) core.EntityID {
	b := NewBranch(
		name, 
		description, 
		address, 
		contacts,
		NewBranchEmployeesProperty(),
		NewBranchNetworkProperty("127.0.0.1"),
	)
	b.SetCompany(*c)
	key := b.EntityID()
	c.Children()[key] = b
	return b.ID
}

func (c Company) GetBranch(id core.EntityID) Branch {
	if b, ok := c.Children()[id]; ok {
		return b
	}

	return Branch{}
}
