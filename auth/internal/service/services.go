package service

// Services is a struct that holds various services.
type Services struct {
	// OAuth2Service is a service for handling OAuth2 related requests.
	OAuth2Service *OAuth2Service
}

// NewServices creates and returns a new Services with the given OAuth2 options.
func NewServices(oAuth2Options *OAuth2ServiceOptions) *Services {
	return &Services{
		OAuth2Service: NewOAuth2Service(oAuth2Options),
	}
}
