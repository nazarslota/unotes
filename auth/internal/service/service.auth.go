package service

import (
	"time"

	"github.com/nazarslota/unotes/auth/internal/service/oauth2"
)

type OAuth2Service struct {
	SignUpRequestHandler  oauth2.SignUpRequestHandler
	SingInRequestHandler  oauth2.SignInRequestHandler
	RefreshRequestHandler oauth2.RefreshRequestHandler
	SignOutRequestHandler oauth2.LogOutRequestHandler
}

type OAuth2ServiceOptions struct {
	AccessTokenCreator   oauth2.AccessTokenCreator
	AccessTokenParser    oauth2.AccessTokenParser
	AccessTokenExpiresIn time.Duration

	RefreshTokenCreator   oauth2.RefreshTokenCreator
	RefreshTokenParser    oauth2.RefreshTokenParser
	RefreshTokenExpiresIn time.Duration

	RefreshTokenSaver    oauth2.RefreshTokenSaver
	RefreshTokenDeleter  oauth2.RefreshTokenDeleter
	RefreshTokensDeleter oauth2.RefreshTokensDeleter
	RefreshTokenGetter   oauth2.RefreshTokenGetter

	UserSaver  oauth2.UserSaver
	UserFinder oauth2.UserFinder
}

func NewOAuth2Service(options OAuth2ServiceOptions) OAuth2Service {
	return OAuth2Service{
		SignUpRequestHandler: oauth2.NewSignUpRequestHandler(
			options.UserSaver,
		),
		RefreshRequestHandler: oauth2.NewRefreshRequestHandler(
			options.AccessTokenCreator,
			options.AccessTokenParser,
			options.AccessTokenExpiresIn,

			options.RefreshTokenCreator,
			options.RefreshTokenParser,
			options.RefreshTokenExpiresIn,

			options.RefreshTokenSaver,
			options.RefreshTokenDeleter,
			options.RefreshTokenGetter,
		),
		SingInRequestHandler: oauth2.NewSignInRequestHandler(
			options.AccessTokenCreator,
			options.AccessTokenExpiresIn,

			options.RefreshTokenCreator,
			options.RefreshTokenExpiresIn,

			options.RefreshTokenSaver,

			options.UserFinder,
		),
		SignOutRequestHandler: oauth2.NewSignOutRequestHandler(
			options.AccessTokenParser,

			options.RefreshTokensDeleter,
			options.RefreshTokenGetter,
		),
	}
}
