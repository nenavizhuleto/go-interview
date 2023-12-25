package branch

import (
	"helpdesk/internals/data"
	"helpdesk/internals/models"
	"helpdesk/internals/models/v3/company"

	"github.com/google/uuid"
)

type Branch struct {
	ID          string `json:"id" db:"id"`
	CompanyID   string `json:"company_id" db:"company_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Address     string `json:"address" db:"address"`
	Contacts    string `json:"contacts" db:"contacts"`
}

func New(company *company.Company, name, description, address, contacts string) (*Branch, error) {
	company_id := company.ID

	return &Branch{
		ID:          newID(),
		CompanyID:   company_id,
		Name:        name,
		Description: description,
		Address:     address,
		Contacts:    contacts,
	}, nil
}

func Get(id string) (*Branch, error) {
	db := data.DB

	var branch Branch
	if err := db.Get(&branch, "SELECT * FROM branches WHERE id = ?", id); err != nil {
		return nil, models.NewDatabaseError("branch", "get", err)
	}

	return &branch, nil
}

func All() ([]Branch, error) {
	db := data.DB

	var branches []Branch
	if err := db.Select(&branches, "SELECT * FROM branches"); err != nil {
		return nil, models.NewDatabaseError("branch", "all", err)
	}

	return branches, nil
}

func (b *Branch) Save() error {
	db := data.DB

	if _, err := Get(b.ID); err != nil {
		// Not exists
		if _, err := db.NamedExec("INSERT INTO branches VALUES (:id, :company_id, :name, :description, :address, :contacts)", b); err != nil {
			return models.NewDatabaseError("branch", "create", err)
		}
	} else {
		// Exists
		if _, err := db.NamedExec("UPDATE branches SET name = :name, description = :description, address = :address, contacts = :contacts", b); err != nil {
			return models.NewDatabaseError("branch", "update", err)
		}
	}

	return nil
}

func (b *Branch) Delete() error {
	db := data.DB

	if _, err := db.Exec("DELETE branches WHERE id = ?", b.ID); err != nil {
		return models.NewDatabaseError("branch", "delete", err)
	}

	return nil
}

// Private functions

func newID() string {
	return uuid.NewString()
}
