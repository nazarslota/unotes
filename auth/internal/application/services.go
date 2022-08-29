package application

type Services struct {
	AuthService *OAuth2Service
}

func NewServices(options OAuth2ServiceOptions) *Services {
	return &Services{AuthService: NewOAuth2Service(options)}
}
