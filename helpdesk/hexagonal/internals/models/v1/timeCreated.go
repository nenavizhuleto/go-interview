package models

import (
	"errors"
	"strings"
	"time"
)

type TimeCreated struct {
	Value string `json:"value,omitempty"`
}

var TimeFormatString = "2006-01-02 15:04:05Z07:00"

func (tc *TimeCreated) Scan(src interface{}) error {
	source, ok := src.(string)
	if !ok {
		return errors.New("Incompitable datetime type")
	}
	// FIXME: 'T' symbol occasionally returned in timeCreated from megaplan
	dt := strings.Replace(string(source), "T", " ", 1)
	if dt == "" {
		tc.Value = time.Now().Format(time.DateTime)
		return nil
	}
	date, err := time.Parse(TimeFormatString, dt)
	if err != nil {
		return errors.New("Couldn't parse datetime")
	}
	tc.Value = date.Format(time.DateTime)
	return nil
}
