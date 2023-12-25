package comment

import (
	"time"
	"helpdesk/internals/domains/user"
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
