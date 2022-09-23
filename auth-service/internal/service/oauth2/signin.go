package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth-service/internal/domain/user"
	"github.com/udholdenhed/unotes/auth-service/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type SignInRequest struct {
	Username string
	Password string
}

type SignInResult struct {
	AccessToken  string
	RefreshToken string
}

type SignInRequestHandler interface {
	Handle(ctx context.Context, r SignInRequest) (SignInResult, error)
}

type signInRequestHandler struct {
	AccessTokenSecret      string
	RefreshTokenSecret     string
	AccessTokenExpiresIn   time.Duration
	RefreshTokenExpiresIn  time.Duration
	UserRepository         user.Repository
	RefreshTokenRepository refreshtoken.Repository
}

func NewSignInRequestHandler(
	accessTokenSecret, refreshTokenSecret string,
	accessTokenExpiresIn, refreshTokenExpiresIn time.Duration,
	userRepository user.Repository, refreshTokenRepository refreshtoken.Repository,
) SignInRequestHandler {
	return &signInRequestHandler{
		AccessTokenSecret:      accessTokenSecret,
		RefreshTokenSecret:     refreshTokenSecret,
		AccessTokenExpiresIn:   accessTokenExpiresIn,
		RefreshTokenExpiresIn:  refreshTokenExpiresIn,
		UserRepository:         userRepository,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (h *signInRequestHandler) Handle(ctx context.Context, r SignInRequest) (SignInResult, error) {
	u, err := h.UserRepository.FindOne(ctx, r.Username)
	if err != nil {
		return SignInResult{}, errors.ErrInternalServerError.SetInternal(err)
	} else if u == nil {
		return SignInResult{}, errors.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("user with this username does not exist"),
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(r.Password)); err != nil {
		return SignInResult{}, errors.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("incorrect user password"),
		)
	}

	accessToken, err := func() (string, error) {
		now := time.Now()
		claims := jwt.MapClaims{
			"exp":     now.Add(h.AccessTokenExpiresIn).Unix(),
			"iat":     now.Unix(),
			"nbf":     now.Unix(),
			"user_id": u.ID,
		}
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.AccessTokenSecret))
	}()
	if err != nil {
		return SignInResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	refreshToken, err := func() (string, error) {
		now := time.Now()
		claims := jwt.MapClaims{
			"exp":     now.Add(h.RefreshTokenExpiresIn).Unix(),
			"iat":     now.Unix(),
			"nbf":     now.Unix(),
			"user_id": u.ID,
		}
		return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(h.RefreshTokenSecret))
	}()
	if err != nil {
		return SignInResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	if err = h.RefreshTokenRepository.SetRefreshToken(ctx, u.ID, refreshtoken.Token{Token: refreshToken}); err != nil {
		return SignInResult{}, errors.ErrInternalServerError.SetInternal(err)
	}

	return SignInResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
