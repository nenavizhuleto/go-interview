package memory

import (
	"errors"
	"helpdesk/internals/core/user"
)

type UserRepo struct {
	name string
	c    map[user.UserID]user.User
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func NewUserRepo(name string) UserRepo {
	return UserRepo{
		name: name,
		c:    make(map[user.UserID]user.User),
	}
}

func (r UserRepo) Get(id user.UserID) (user.User, error) {
	if c, ok := r.c[id]; ok {
		return c, nil
	} else {
		return user.User{}, ErrUserNotFound
	}
}

func (r UserRepo) All() ([]user.User, error) {
	users := make([]user.User, 0)
	for _, c := range r.c {
		users = append(users, c)
	}

	return users, nil
}

func (r UserRepo) Save(u user.User) error {
	r.c[u.ID] = u
	return nil
}

func (r UserRepo) Delete(id user.UserID) error {
	if _, ok := r.c[id]; ok {
		delete(r.c, id)
		return nil
	} else {
		return ErrUserNotFound
	}
}
