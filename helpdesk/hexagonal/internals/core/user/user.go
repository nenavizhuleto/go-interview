package user

import (
	"helpdesk/internals/core/company"

	"github.com/google/uuid"
)

type UserID string
type Login string // IP Address

type UserPropertyName string
type UserProperty interface {
	GetPropertyName() UserPropertyName
}

type User struct {
	ID         UserID                            `json:"id"`
	Name       string                            `json:"name"`
	Phone      string                            `json:"phone"`
	CompanyID  company.CompanyID                 `json:"company_id"`
	BranchID   company.BranchID                  `json:"branch_id"`
	Properties map[UserPropertyName]UserProperty `json:"properties"`
}

func NewUser(name string, phone string) User {
	return User{
		ID:         UserID(uuid.NewString()),
		Name:       name,
		Phone:      phone,
		Properties: make(map[UserPropertyName]UserProperty),
	}
}

func (u *User) SetCompany(id company.CompanyID) {
	u.CompanyID = id
}

func (u *User) SetBranch(id company.BranchID) {
	u.BranchID = id
}
