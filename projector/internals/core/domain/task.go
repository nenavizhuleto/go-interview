package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus int

const (
	StatusTodo TaskStatus = iota
	StatusInProgress
	StatusDone
)

type Task struct {
	ID          string
	title       string
	description string
	Assignee    *Person
	CreatedAt   time.Time
	Deadline    time.Duration
	Status      TaskStatus
}

func NewTask(title, description string) *Task {
	return &Task{
		ID:          uuid.NewString(),
		title:       title,
		description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
	}
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t *Task) ChangeStatus(status TaskStatus) {
	t.Status = status
}

func (t *Task) MarkTodo() {
	t.ChangeStatus(StatusTodo)
}

func (t *Task) MarkInProgress() {
	t.ChangeStatus(StatusInProgress)
}

func (t *Task) MarkDone() {
	t.ChangeStatus(StatusDone)
}

func (t *Task) AssignTo(person *Person) {
	t.Assignee = person
}

func (t *Task) SetDeadline(duration time.Duration) {
	t.Deadline = duration
}

func (t *Task) RemoveDeadline() {
	t.Deadline = 0
}
