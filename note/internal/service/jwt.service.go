package service

import servicejwt "github.com/nazarslota/unotes/note/internal/service/jwt"

type JWTService struct {
	AccessTokenValidator servicejwt.AccessTokenValidator
}

type JWTServiceOptions struct {
	AccessTokenSecret string
}

func NewJWTService(options JWTServiceOptions) JWTService {
	return JWTService{
		AccessTokenValidator: servicejwt.NewAccessTokenValidator(options.AccessTokenSecret),
	}
}
