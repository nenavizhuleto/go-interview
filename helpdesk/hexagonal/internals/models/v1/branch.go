package models

import "helpdesk/internals/data"

type Branch struct {
	ID          string `json:"id" db:"id"`
	CompanyID   string `json:"company_id" db:"company_id"`
	Name        string `json:"name" db:"name"`
	Address     string `json:"address" db:"address"`
	Contacts    string `json:"contacts" db:"contacts"`
	Description string `json:"description" db:"description"`
}

func GetBranchFromSubnet(s *Subnet) (*Branch, error) {
	db := data.DB
	var branch Branch
	if err := db.Get(&branch, "SELECT * FROM branches WHERE id = $1", s.BranchID); err != nil {
		return nil, err
	}
	return &branch, nil
}
