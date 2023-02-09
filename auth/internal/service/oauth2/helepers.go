package oauth2

import "github.com/golang-jwt/jwt/v4"

type AccessTokenManager[T jwt.Claims] interface {
	New(claims T) (string, error)
	Parse(token string) (T, error)
}

type RefreshTokenManager[T jwt.Claims] interface {
	New(claims T) (string, error)
	Parse(token string) (T, error)
}
