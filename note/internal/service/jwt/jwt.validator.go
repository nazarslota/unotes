package jwt

import "github.com/nazarslota/unotes/auth/pkg/jwt"

type AccessTokenValidator interface {
	Validate(token string) (jwt.AccessTokenClaims, error)
}

type accessTokenValidator struct {
	AccessTokenManager *jwt.AccessTokenManagerHMAC
}

func NewAccessTokenValidator(accessTokenSecret string) AccessTokenValidator {
	return &accessTokenValidator{
		AccessTokenManager: jwt.NewAccessTokenManagerHMAC(accessTokenSecret),
	}
}

func (v accessTokenValidator) Validate(token string) (jwt.AccessTokenClaims, error) {
	return v.AccessTokenManager.Parse(token)
}
