package megaplan

import (
	models "helpdesk/internals/models/v3/comment"
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

func filterComment(_comment *Comment) *models.Comment {
	var comment models.Comment
	content := _comment.Content
	content = strings.ReplaceAll(content, "\u003cp\u003e", "")
	content = strings.ReplaceAll(content, "\u003c/p\u003e", "")

	if strings.HasPrefix(content, CommentTagTo) {
		comment.Direction = models.DirectionTo
		content = strings.Replace(content, CommentTagTo, "", 1)
	} else if strings.HasPrefix(content, CommentTagFrom) {
		comment.Direction = models.DirectionFrom
		content = strings.Replace(content, CommentTagFrom, "", 1)
	} else {
		return nil
	}
	comment.ID = _comment.ID
	comment.Content = content
	if _comment.TimeCreated != nil {
		comment.TimeCreated = _comment.TimeCreated.Value
	}

	return &comment
}

func (dto *TaskDTO) GetLastComment() *models.Comment {
	return filterComment(dto.LastComment)

}
func (dto *TaskDTO) GetComments() []models.Comment {
	var comments = make([]models.Comment, 0)
	if dto.LastComment != nil {
		comment := filterComment(dto.LastComment)
		if comment != nil {
			comments = append(comments, *comment)
		}
	}
	if dto.CommentsCount == 0 {
		return []models.Comment{}
	}

	for _, _comment := range dto.Comments {
		// Parse comment.Content and determine which direction it is
		comment := filterComment(&_comment)
		if comment != nil {
			comments = append(comments, *comment)
		}
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
