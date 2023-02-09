package service

import (
	"time"

	"github.com/nazarslota/unotes/auth/internal/domain/refresh"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/nazarslota/unotes/auth/internal/service/oauth2"
	"github.com/nazarslota/unotes/auth/pkg/jwt"
)

type OAuth2Service struct {
	SignUpRequestHandler  oauth2.SignUpRequestHandler
	SingInRequestHandler  oauth2.SignInRequestHandler
	RefreshRequestHandler oauth2.RefreshRequestHandler
	SignOutRequestHandler oauth2.LogOutRequestHandler
}

type OAuth2ServiceOptions struct {
	AccessTokenManager     oauth2.AccessTokenManager[jwt.AccessTokenClaims]
	AccessTokenExpiresIn   time.Duration
	RefreshTokenManager    oauth2.RefreshTokenManager[jwt.RefreshTokenClaims]
	RefreshTokenExpiresIn  time.Duration
	UserRepository         user.Repository
	RefreshTokenRepository refresh.Repository
}

func NewOAuth2Service(options OAuth2ServiceOptions) OAuth2Service {
	return OAuth2Service{
		SignUpRequestHandler: oauth2.NewSignUpRequestHandler(
			options.UserRepository,
		),
		SingInRequestHandler: oauth2.NewSignInRequestHandler(
			options.AccessTokenManager,
			options.AccessTokenExpiresIn,
			options.RefreshTokenManager,
			options.RefreshTokenExpiresIn,
			options.UserRepository,
			options.RefreshTokenRepository,
		),
		RefreshRequestHandler: oauth2.NewRefreshRequestHandler(
			options.AccessTokenManager,
			options.AccessTokenExpiresIn,
			options.RefreshTokenManager,
			options.RefreshTokenExpiresIn,
			options.UserRepository,
			options.RefreshTokenRepository,
		),
		SignOutRequestHandler: oauth2.NewSignOutRequestHandler(
			options.AccessTokenManager,
			options.RefreshTokenRepository,
		),
	}
}
