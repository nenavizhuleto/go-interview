package company

import (
	"helpdesk/internals/core"

	"github.com/google/uuid"
)

type BranchPropertyName string
type BranchProperty interface {
	GetPropertyName() BranchPropertyName
}

type Branch struct {
	ID          core.EntityID                         `json:"id"`
	CompanyID   core.EntityID                         `json:"company_id"`
	Name        string                                `json:"name"`
	Description string                                `json:"description"`
	Address     string                                `json:"address"`
	Contacts    string                                `json:"contacts"`
	Properties  map[BranchPropertyName]BranchProperty `json:"properties"`
}

// Implementing Entity interface
// -----------------------------
func (b Branch) EntityID() core.EntityID {
	return b.ID
}

// Implementing IChild interface
func (b Branch) Parent() core.Entity {
	return Company{ID: b.CompanyID}
}

func (b *Branch) SetCompany(c Company) {
	b.CompanyID = c.EntityID()
}

func NewBranch(name, description, address, contacts string, properties ...BranchProperty) Branch {
	branch := Branch{
		ID:          core.EntityID(uuid.NewString()),
		Name:        name,
		Description: description,
		Address:     address,
		Contacts:    contacts,
		Properties:  make(map[BranchPropertyName]BranchProperty),
	}

	if len(properties) > 0 {
		for _, prop := range properties {
			branch.AddProperty(prop)
		}
	}

	return branch
}

func (b *Branch) SetName(name string) {
	b.Name = name
}

func (b *Branch) SetDescription(description string) {
	b.Description = description
}

func (b *Branch) SetAddress(address string) {
	b.Address = address
}

func (b *Branch) SetContacts(contacts string) {
	b.Contacts = contacts
}

func (b *Branch) AddProperty(property BranchProperty) {
	add := func(p BranchProperty) {
		b.Properties[p.GetPropertyName()] = p
	}
	switch prop := property.(type) {
	case BranchEmployeesProperty:
		if FeatureEnabled(FeatureBranchEmployees) {
			add(prop)
		}
	case BranchNetworkProperty:
		if FeatureEnabled(FeatureBranchNetwork) {
			add(prop)
		}
	default:
		break
	}
}

func (b *Branch) GetProperty(name BranchPropertyName) BranchProperty {
	property, ok := b.Properties[name]
	if !ok {
		return nil
	}
	switch prop := property.(type) {
	case BranchEmployeesProperty:
		if FeatureEnabled(FeatureBranchEmployees) {
			return prop
		}
	case BranchNetworkProperty:
		if FeatureEnabled(FeatureBranchNetwork) {
			return prop
		}
	default:
		return nil
	}

	return nil
}
