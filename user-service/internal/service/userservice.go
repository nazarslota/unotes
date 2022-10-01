package service

import (
	"github.com/udholdenhed/unotes/user-service/internal/domain/user"
	userservice "github.com/udholdenhed/unotes/user-service/internal/service/user"
)

type UserService struct {
	CreateUserRequestHandler         userservice.CreateUserRequestHandler
	FindUserByIDRequestHandler       userservice.FindUserByIDRequestHandler
	FindUserByUsernameRequestHandler userservice.FindUserByUsernameRequestHandler
}

type UserServiceOptions struct {
	UserRepository user.Repository
}

func NewUserService(options *UserServiceOptions) *UserService {
	return &UserService{
		CreateUserRequestHandler: userservice.NewCreateUserRequestHandler(
			options.UserRepository,
		),
		FindUserByIDRequestHandler: userservice.NewFindUserByIDRequestHandler(
			options.UserRepository,
		),
		FindUserByUsernameRequestHandler: userservice.NewFindByUsernameRequestHandler(
			options.UserRepository,
		),
	}
}
