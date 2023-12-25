package megaplan

import (
	"helpdesk/internals/domains/comment"
	"strings"
	"time"
)

type Employee struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Name      string `json:"name,omitempty"`
}

type TimeCreated struct {
	Value time.Time `json:"value,omitempty"`
}

type Comment struct {
	ID          string       `json:"id"`
	Content     string       `json:"content"`
	TimeCreated *TimeCreated `json:"timeCreated"`
}

// Piece of information to create task in Megaplan
type TaskDTO struct {
	ID            string       `json:"id,omitempty"`
	Name          string       `json:"name"`
	Subject       string       `json:"subject"`
	Responsible   Employee     `json:"responsible"`
	IsUrgent      bool         `json:"isUrgent"`
	IsTemplate    bool         `json:"isTemplate"`
	TimeCreated   *TimeCreated `json:"timeCreated,omitempty"`
	Activity      *TimeCreated `json:"activity,omitempty"`
	Status        string       `json:"status,omitempty"`
	LastComment   *Comment     `json:"lastComment,omitempty"`
	Comments      []Comment    `json:"comments,omitempty"`
	CommentsCount int          `json:"commentsCount,omitempty"`
}

const (
	StatusCreated   = "created"
	StatusAssigned  = "assigned"
	StatusAccepted  = "accepted"
	StatusDone      = "done"
	StatusCompleted = "completed"
	StatusRejected  = "rejected"
	StatusCancelled = "cancelled"
	StatusExpired   = "expired"
	StatusDelayed   = "delayed"
	StatusTemplate  = "template"
	StatusOverdue   = "overdue"
)

const (
	CommentTagTo   = "#[TOUSER]:"
	CommentTagFrom = "#[FROMUSER]:"
)

func (dto *TaskDTO) GetComments() []comment.Comment {
	if dto.CommentsCount == 0 {
		return []comment.Comment{}
	}

	var comments = make([]comment.Comment, 0)
	for _, _comment := range dto.Comments {
		// Parse comment.Content and determine which direction it is
		var comment_ comment.Comment
		content := _comment.Content
		if strings.Contains(content, CommentTagTo) {
			comment_.Direction = comment.DirectionTo
		} else if strings.Contains(content, CommentTagFrom) {
			comment_.Direction = comment.DirectionFrom
		} else {
			continue
		}

		comment_.ID = _comment.ID
		comment_.Content = _comment.Content
		if _comment.TimeCreated != nil {
			comment_.TimeCreated = _comment.TimeCreated.Value
		}

		comments = append(comments, comment_)
	}

	return comments
}

func (dto *TaskDTO) GetStatus() string {
	status := dto.Status
	switch status {
	case StatusAssigned:
		if dto.Responsible.ID == MP.Responsible {
			return StatusCreated
		} else {
			return StatusAssigned
		}
	case StatusAccepted:
		if dto.Responsible.ID == MP.Responsible {
			return StatusCreated
		} else {
			return StatusAccepted
		}
	default:
		return status
	}
}

type TaskEvent struct {
	Data  TaskDTO `json:"data"`
	Event string  `json:"event"`
	Model string  `json:"model"`
}
