package comment

import (
	"helpdesk/internals/models/v3/user"
	"time"

	"github.com/google/uuid"
)

const (
	DirectionTo   = "to"
	DirectionFrom = "from"
)

type Comment struct {
	ID          string     `json:"id"`
	Content     string     `json:"content"`
	User        *user.User `json:"user"`
	Direction   string     `json:"direction"`
	TimeCreated time.Time  `json:"timeCreated"`
}

func NewComment(content string, dir string) Comment {
	return Comment{
		ID:          uuid.NewString(),
		Content:     content,
		Direction:   dir,
		TimeCreated: time.Now(),
	}
}
