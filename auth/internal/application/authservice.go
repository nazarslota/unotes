package application

import (
	"time"

	"github.com/udholdenhed/unotes/auth/internal/application/oauth2"
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth/internal/domain/user"
)

type OAuth2Service struct {
	SignUpRequestHandler  oauth2.SignUpRequestHandler
	SingInRequestHandler  oauth2.SignInRequestHandler
	RefreshRequestHandler oauth2.RefreshRequestHandler
	SignOutRequestHandler oauth2.LogOutRequestHandler
}

type OAuth2ServiceOptions struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	UserRepository         user.Repository
	RefreshTokenRepository refreshtoken.Repository
}

func NewOAuth2Service(options OAuth2ServiceOptions) *OAuth2Service {
	return &OAuth2Service{
		SignUpRequestHandler: oauth2.NewSignUpRequestHandler(
			options.UserRepository,
		),
		SingInRequestHandler: oauth2.NewSignInRequestHandler(
			options.AccessTokenSecret,
			options.RefreshTokenSecret,
			options.AccessTokenExpiresIn,
			options.RefreshTokenExpiresIn,
			options.UserRepository,
			options.RefreshTokenRepository,
		),
		RefreshRequestHandler: oauth2.NewRefreshRequestHandler(
			options.AccessTokenSecret,
			options.RefreshTokenSecret,
			options.AccessTokenExpiresIn,
			options.RefreshTokenExpiresIn,
			options.RefreshTokenRepository,
		),
		SignOutRequestHandler: oauth2.NewSignOutRequestHandler(
			options.AccessTokenSecret,
			options.RefreshTokenRepository,
		),
	}
}
