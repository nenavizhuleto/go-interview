package company

import (
	"helpdesk/internals/data"
	"helpdesk/internals/models"
	"strings"
)

type Company struct {
	ID   string `json:"id" db:"id"` // Company's ITN Number
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

func New(itn string, name string) (*Company, error) {
	validITN, err := newITN(itn)
	if err != nil {
		return nil, err
	}
	validName, err := newName(name)
	if err != nil {
		return nil, err
	}

	return &Company{
		ID:   validITN,
		Name: validName,
		Slug: newSlug(validName),
	}, nil
}

func Get(id string) (*Company, error) {
	db := data.DB

	var company Company
	if err := db.Get(&company, "SELECT * FROM companies WHERE id = ?", id); err != nil {
		return nil, models.NewDatabaseError("company", "get", err)
	}

	return &company, nil
}

func All() ([]Company, error) {
	db := data.DB

	var companies []Company
	if err := db.Select(&companies, "SELECT * FROM companies"); err != nil {
		return nil, models.NewDatabaseError("company", "all", err)
	}

	return companies, nil
}

func (c *Company) Save() error {
	db := data.DB

	if _, err := Get(c.ID); err != nil {
		if _, err := db.NamedExec("INSERT INTO companies VALUES (:id, :name, :slug);", c); err != nil {
			return models.NewDatabaseError("company", "save", err)
		}
	} else {
		if _, err := db.NamedExec("UPDATE companies SET name = :name, slug = :slug WHERE id = :id;", c); err != nil {
			return models.NewDatabaseError("company", "update", err)
		}
	}

	return nil
}

func (c *Company) Delete() error {
	db := data.DB

	if _, err := db.Exec("DELETE FROM companies WHERE id = ?", c.ID); err != nil {
		return models.NewDatabaseError("company", "delete", err)
	}

	return nil
}

// Private functions

func newITN(itn string) (string, error) {
	if len(itn) != 10 {
		return "", models.NewValidationError("company", "itn", "ITN must be exaclty 10 digits long")
	}

	return itn, nil
}

func newName(name string) (string, error) {
	// Check if name exists in database

	// If name is empty
	if len(name) == 0 {
		return "", models.NewValidationError("company", "name")
	}

	return name, nil
}

func newSlug(validName string) string {
	slug := strings.ToLower(validName)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
