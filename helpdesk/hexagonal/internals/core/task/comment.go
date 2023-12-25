package task

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/user"
	"time"
)

const (
	DirectionTo   = "to"
	DirectionFrom = "from"
)

type Comment struct {
	ID          string      `json:"id"`
	Content     string      `json:"content"`
	UserID      user.UserID `json:"user_id"`
	Direction   string      `json:"direction"`
	TimeCreated time.Time   `json:"timeCreated"`
	owned       bool
}

func (c Comment) Owned() bool {
	return c.owned
}

func (c *Comment) SetOwner(o core.Owner) error {
	if u, ok := o.(user.User); ok {
		c.UserID = u.ID
		c.owned = true
		return nil
	}
	return core.ErrInvalidOwner
}
