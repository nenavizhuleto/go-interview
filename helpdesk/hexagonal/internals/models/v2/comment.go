package models

import "time"

const (
	DirectionTo   = "to"
	DirectionFrom = "from"
)

type Comment struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	User        User      `json:"user"`
	Direction   string    `json:"direction"`
	TimeCreated time.Time `json:"timeCreated"`
}
