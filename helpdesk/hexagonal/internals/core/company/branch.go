package company

import (
	"github.com/google/uuid"
)

type BranchID string
type Property interface {
	PropertyName() string
	Value() any
}

type BranchAddress string

func (p BranchAddress) PropertyName() string {
	return "address"
}

func (p BranchAddress) Value() string {
	return string(p)
}

type BranchContacts string

func (p BranchContacts) PropertyName() string {
	return "contacts"
}

func (p BranchContacts) Value() string {
	return string(p)
}

type Branch struct {
	ID          BranchID            `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Employees   []Employee          `json:"employees"`
	Properties  map[string]Property `json:"properties"`
}

func NewBranch(name, description, address, contacts string) Branch {
	return Branch{
		ID:          BranchID(uuid.NewString()),
		Name:        name,
		Description: description,
		Properties:  make(map[string]Property),
	}
}

func (b *Branch) SetName(name string) {
	b.Name = name
}

func (b *Branch) SetDescription(description string) {
	b.Description = description
}

func (b *Branch) AddEmployee(e Employee) {
	b.Employees = append(b.Employees, e)
}

func (b Branch) GetEmployee(id EmployeeID) Employee {
	var employee Employee
	for _, e := range b.Employees {
		if e.ID == id {
			employee = e
			return employee
		}
	}
	return employee
}

func (b *Branch) AddProperty(property Property) {
	b.Properties[property.PropertyName()] = property
}
