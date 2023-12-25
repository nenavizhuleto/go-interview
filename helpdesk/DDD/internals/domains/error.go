package models

import (
	"encoding/json"
	"fmt"
)

type ErrorType string

var (
	ErrTypeValidation = ErrorType("validation")
	ErrTypeDatabase   = ErrorType("database")
	ErrTypeParse      = ErrorType("parse")
)

type ErrorBody map[string]any

type Error struct {
	Type ErrorType `json:"type"`
	Body ErrorBody `json:"body"`
}

func (e Error) Error() string {
	body, err := json.Marshal(e.Body)
	if err != nil {
		return fmt.Sprintf("type: %s body: undefined", e.Type)
	}
	return fmt.Sprintf("type: %s body: %s", e.Type, string(body))
}

func NewValidationError(entity string, field string, message ...string) Error {
	return Error{
		Type: ErrTypeValidation,
		Body: ErrorBody{
			"entity": entity,
			"field":  field,
			"message": message,
		},
	}
}

func NewDatabaseError(entity string, action string, errors ...error) Error {
	var errorMsgs []string
	if len(errors) > 0 {
		for _, error := range errors {
			errorMsgs = append(errorMsgs, error.Error())
		}
	}
	return Error{
		Type: ErrTypeDatabase,
		Body: ErrorBody{
			"entity": entity,
			"action": action,
			"errors": errorMsgs,
		},
	}
}

func NewParseError(endpoint string, err error) Error {
	return Error{
		Type: ErrTypeParse,
		Body: ErrorBody{
			"endpoint": endpoint,
			"message":  err.Error(),
		},
	}

}
