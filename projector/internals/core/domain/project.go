package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type Project struct {
	ID          string
	Name        string
	description string
	Tasks       map[string]*Task
}

func NewProject(name, description string) *Project {
	return &Project{
		ID:          uuid.NewString(),
		Name:        name,
		description: description,
		Tasks:       make(map[string]*Task),
	}
}

func (p *Project) AddTask(task *Task) {
	p.Tasks[task.ID] = task
}

func (p *Project) RemoveTask(task_id string) {
	delete(p.Tasks, task_id)
}

// FilterValue implements list.Item
func (p Project) FilterValue() string {
	return p.Name
}

// Title implements list.Item
func (p Project) Title() string {
	return p.Name
}

// Description implements list.Item
func (p Project) Description() string {
	return fmt.Sprintf("%s. Tasks: %d", p.description, len(p.Tasks))
}
