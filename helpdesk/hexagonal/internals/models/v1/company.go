package models

import "helpdesk/internals/data"

type Company struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
}

func GetCompanyFromBranch(b *Branch) (*Company, error) {
	db := data.DB
	var company Company
	if err := db.Get(&company, "SELECT * FROM companies WHERE id = $1", b.CompanyID); err != nil {
		return nil, err
	}
	return &company, nil
}
