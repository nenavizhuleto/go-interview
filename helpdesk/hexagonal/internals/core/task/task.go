package task

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/company"
	"helpdesk/internals/core/user"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string
type TaskID string

const (
	StatusCreated   = TaskStatus("created")
	StatusAssigned  = TaskStatus("assigned")
	StatusAccepted  = TaskStatus("accepted")
	StatusDone      = TaskStatus("done")
	StatusCompleted = TaskStatus("completed")
	StatusRejected  = TaskStatus("rejected")
	StatusCancelled = TaskStatus("cancelled")
	StatusExpired   = TaskStatus("expired")
	StatusDelayed   = TaskStatus("delayed")
	StatusTemplate  = TaskStatus("template")
	StatusOverdue   = TaskStatus("overdue")
)

type Task struct {
	ID           TaskID            `json:"id"`
	Name         string            `json:"name"`
	Subject      string            `json:"subject"`
	Status       TaskStatus        `json:"status"`
	TimeCreated  time.Time         `json:"created_at"`
	LastActivity time.Time         `json:"activity_at"`
	CompanyID    company.CompanyID `json:"company_id"`
	BranchID     company.BranchID  `json:"branch_id"`
	UserID       user.UserID       `json:"user_id"`
	Comments     []Comment         `json:"comments"`
	owned        bool
}

func NewTask(name string, subject string) Task {
	created_at := time.Now()
	t := Task{
		ID:           TaskID(uuid.NewString()),
		Name:         name,
		Subject:      subject,
		Status:       StatusCreated,
		TimeCreated:  created_at,
		LastActivity: created_at,
		Comments:     make([]Comment, 0),
		owned:        false,
	}

	return t
}

// Implements ownable
func (t Task) Owned() bool {
	return t.owned
}

func (t *Task) SetOwner(o core.Owner) error {
	if u, ok := o.(user.User); ok {
		t.UserID = u.ID
		t.BranchID = u.BranchID
		t.CompanyID = u.CompanyID
		t.owned = true
		return nil
	}
	return core.ErrInvalidOwner
}
