package service

type Service struct {
	UserService *UserService
}

func NewService(userServiceOptions *UserServiceOptions) *Service {
	return &Service{
		UserService: NewUserService(userServiceOptions),
	}
}
