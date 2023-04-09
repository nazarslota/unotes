package oauth2

import (
	"context"

	domainrefresh "github.com/nazarslota/unotes/auth/internal/domain/refresh"
	domainuser "github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/nazarslota/unotes/auth/pkg/jwt"
)

type AccessTokenCreator interface {
	New(claims jwt.AccessTokenClaims) (string, error)
}

type AccessTokenParser interface {
	Parse(token string) (jwt.AccessTokenClaims, error)
}

type RefreshTokenCreator interface {
	New(claims jwt.RefreshTokenClaims) (string, error)
}

type RefreshTokenParser interface {
	Parse(token string) (jwt.RefreshTokenClaims, error)
}

type RefreshTokenSaver interface {
	SaveRefreshToken(ctx context.Context, userID string, token domainrefresh.Token) error
}

type RefreshTokenDeleter interface {
	DeleteRefreshToken(ctx context.Context, userID string, token domainrefresh.Token) error
}

type RefreshTokensDeleter interface {
	DeleteRefreshTokens(ctx context.Context, userID string, tokens []domainrefresh.Token) error
}

type RefreshTokenGetter interface {
	GetRefreshTokens(ctx context.Context, userID string) ([]domainrefresh.Token, error)
}

type UserSaver interface {
	SaveUser(ctx context.Context, user domainuser.User) error
}

type UserFinder interface {
	FindUserByUsername(ctx context.Context, username string) (domainuser.User, error)
}
