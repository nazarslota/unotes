package service

type Service struct {
	OAuth2Service *OAuth2Service
}

func NewService(oAuth2Options *OAuth2ServiceOptions) *Service {
	return &Service{OAuth2Service: NewOAuth2Service(oAuth2Options)}
}
