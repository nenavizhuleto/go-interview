package models

import "errors"

type Model interface {
	Create() error
	Save() error
	Update() error
	Delete() error
}

func NewError(message string) error {
	return errors.New("model: " + message)
}
