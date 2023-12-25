package core

import "errors"

var (
	ErrInvalidOwner = errors.New("invalid owner specified")
)

type Ownable interface {
	Owned() bool
	SetOwner(Owner) error
}

type Owner interface{}
