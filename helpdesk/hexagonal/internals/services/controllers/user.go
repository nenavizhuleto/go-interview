package controllers

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/user"
	"helpdesk/internals/services/persistance/memory"
)

type UserRepo core.Repo[user.UserID, user.User]

type UserController struct {
	Repo UserRepo
}

func NewUserController(repo UserRepo) UserController {
	return UserController{
		Repo: repo,
	}
}

func NewMemoryUserController() UserController {
	repo := memory.NewUserRepo("users")
	return UserController{
		Repo: &repo,
	}
}

func (cc *UserController) GetUsers() ([]user.User, error) {
	return cc.Repo.All()
}

func (cc *UserController) CreateUser(name string, phone string) (user.User, error) {
	u := user.NewUser(name, phone)

	if err := cc.Repo.Save(u); err != nil {
		return user.User{}, err
	}

	return u, nil

}
