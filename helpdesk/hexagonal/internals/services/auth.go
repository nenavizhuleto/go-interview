package services

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/user"
)

type UserRepo core.Repo[user.UserID, user.User]

type AuthService struct {
	u UserRepo
}

func NewAuthService(userRepo UserRepo) AuthService {
	return AuthService{
		u: userRepo,
	}
}

func (s AuthService) Signup(name string, phone string) (user.User, error) {
	return user.User{}, nil
}

func (s AuthService) Signin(ip user.Login) (user.User, error) {
	return user.User{}, nil
}

func (s AuthService) Signout(ip user.Login) error {
	return nil
}
