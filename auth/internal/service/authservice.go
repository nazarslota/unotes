package service

import (
	"time"

	"github.com/nazarslota/unotes/auth/internal/domain/refreshtoken"
	"github.com/nazarslota/unotes/auth/internal/domain/user"
	"github.com/nazarslota/unotes/auth/internal/service/oauth2"
)

// OAuth2Service is a struct that holds handlers for various OAuth2 related requests.
type OAuth2Service struct {
	// SignUpRequestHandler is a handler for OAuth2 sign up requests.
	SignUpRequestHandler oauth2.SignUpRequestHandler
	// SingInRequestHandler is a handler for OAuth2 sign in requests.
	SingInRequestHandler oauth2.SignInRequestHandler
	// RefreshRequestHandler is a handler for OAuth2 token refresh requests.
	RefreshRequestHandler oauth2.RefreshRequestHandler
	// SignOutRequestHandler is a handler for OAuth2 sign out requests.
	SignOutRequestHandler oauth2.LogOutRequestHandler
}

// OAuth2ServiceOptions is a struct that holds options for configuring an OAuth2Service.
type OAuth2ServiceOptions struct {
	// AccessTokenSecret is the secret used to sign the access tokens.
	AccessTokenSecret string
	// RefreshTokenSecret is the secret used to sign the refresh tokens.
	RefreshTokenSecret string
	// AccessTokenExpiresIn is the duration after which the access tokens expire.
	AccessTokenExpiresIn time.Duration
	// RefreshTokenExpiresIn is the duration after which the refresh tokens expire.
	RefreshTokenExpiresIn time.Duration
	// UserRepository is a repository for storing and retrieving users.
	UserRepository user.Repository
	// RefreshTokenRepository is a repository for storing and retrieving refresh tokens.
	RefreshTokenRepository refreshtoken.Repository
}

// NewOAuth2Service creates and returns a new OAuth2Service with the given options.
func NewOAuth2Service(options *OAuth2ServiceOptions) *OAuth2Service {
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
