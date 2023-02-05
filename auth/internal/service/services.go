package service

type Services struct {
	OAuth2Service OAuth2Service
}

func NewServices(oAuth2Options OAuth2ServiceOptions) Services {
	return Services{OAuth2Service: NewOAuth2Service(oAuth2Options)}
}
