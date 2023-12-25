package models

type Branch struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	CompanyID   string `json:"company_id" db:"company_id"`
	Address     string `json:"address" db:"address"`
	Contacts    string `json:"contacts" db:"contacts"`
	Description string `json:"description" db:"description"`
}
