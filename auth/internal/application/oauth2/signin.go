package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/udholdenhed/unotes/auth/internal/domain/refreshtoken"
	"github.com/udholdenhed/unotes/auth/internal/domain/user"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
	"github.com/udholdenhed/unotes/auth/pkg/utils"
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
	Handle(r SignInRequest) (SignInResult, error)
}

type signInRequestHandler struct {
	AccessTokenSecret     string
	RefreshTokenSecret    string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration

	UserRepository         user.Repository
	RefreshTokenRepository refreshtoken.Repository
}

func NewSignInRequestHandler(ats, rts string, ate, rte time.Duration, ur user.Repository, rtr refreshtoken.Repository) SignInRequestHandler {
	return &signInRequestHandler{
		AccessTokenSecret:      ats,
		RefreshTokenSecret:     rts,
		AccessTokenExpiresIn:   ate,
		RefreshTokenExpiresIn:  rte,
		UserRepository:         ur,
		RefreshTokenRepository: rtr,
	}
}

func (h *signInRequestHandler) Handle(r SignInRequest) (SignInResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := h.UserRepository.FindOne(ctx, r.Username)
	if err != nil {
		return SignInResult{}, errors.ErrInternalServerError.SetInternal(err)
	} else if u == nil {
		return SignInResult{}, errors.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("user with this username does not exist"),
		)
	}

	if err := utils.ComparePassword(u.PasswordHash, r.Password); err != nil {
		return SignInResult{}, errors.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("incorrect user password"),
		)
	}

	at, err := func() (string, error) {
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

	rt, err := func() (string, error) {
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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = h.RefreshTokenRepository.SetRefreshToken(ctx, u.ID, refreshtoken.RefreshToken{Token: rt}); err != nil {
		return SignInResult{}, errors.ErrInternalServerError.SetInternal(err)
	}
	return SignInResult{AccessToken: at, RefreshToken: rt}, nil
}
