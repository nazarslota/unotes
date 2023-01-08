package service

// Service is a struct that holds various services.
type Service struct {
	// OAuth2Service is a service for handling OAuth2 related requests.
	OAuth2Service *OAuth2Service
}

// NewService creates and returns a new Service with the given OAuth2 options.
func NewService(oAuth2Options *OAuth2ServiceOptions) *Service {
	return &Service{
		OAuth2Service: NewOAuth2Service(oAuth2Options),
	}
}
