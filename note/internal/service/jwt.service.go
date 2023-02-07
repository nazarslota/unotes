package service

import servicejwt "github.com/nazarslota/unotes/note/internal/service/jwt"

type JWTService struct {
	JWTVerifier servicejwt.Verifier
}

type JWTServiceOptions struct {
	AccessTokenSecret string
}

func NewJWTService(options JWTServiceOptions) JWTService {
	return JWTService{
		JWTVerifier: servicejwt.NewVerifier(options.AccessTokenSecret),
	}
}
