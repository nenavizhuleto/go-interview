package main

import (
	"time"
)

type Transaction struct {
	ID          int64
	AccountID   int64
	Timestamp   time.Time
	Description string
	Value       float64
}

type Account struct {
	ID      int64
	Name    string
	Balance float64
}
