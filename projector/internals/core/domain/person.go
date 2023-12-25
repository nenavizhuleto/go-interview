package domain

import (
	"github.com/google/uuid"
)

type Role string

type Person struct {
	ID   string
	Name string
	Role Role
}

func (p Person) FilterValue() string {
	return p.Name
}

func (p Person) Title() string {
	return p.Name
}

func (p Person) Description() string {
	return string(p.Role)
}

func NewPerson(name string, role Role) *Person {
	return &Person{
		ID:   uuid.NewString(),
		Name: name,
		Role: role,
	}
}

func (p *Person) CreateTask(project *Project, title, description string) {
	task := NewTask(title, description)
	task.AssignTo(p)
	project.AddTask(task)
}
