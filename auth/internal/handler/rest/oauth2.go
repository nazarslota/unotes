package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	serviceoauth2 "github.com/udholdenhed/unotes/auth/internal/service/oauth2"
)

type oAuth2SignUpUserModel struct {
	Username string `json:"username" validate:"required,min=4,max=32" example:"username"`
	Password string `json:"password" validate:"required,min=8,max=64" example:"password"`
}

// @Summary     oAuth2 sign up
// @Description create account
// @Tags        oAuth2
// @Accept      json
// @Produce     json
// @Param       input body oAuth2SignUpUserModel true "account info"
// @Success     204
// @Failure     400     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /serviceoauth2/sign-up [post]
func (h *Handler) oAuth2SignUp(c echo.Context) error {
	input := new(oAuth2SignUpUserModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := &serviceoauth2.SignUpRequest{
		Username: input.Username,
		Password: input.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := h.services.OAuth2Service.SignUpRequestHandler.Handler(ctx, request)
	if err != nil {
		if errors.Is(err, serviceoauth2.ErrUserAlreadyExist) {
			return echo.NewHTTPError(http.StatusBadRequest, "user already exist").SetInternal(err)
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}
	return c.NoContent(http.StatusNoContent)
}

type oAuth2SignInUserModel struct {
	Username string `json:"username" validate:"required" example:"username"`
	Password string `json:"password" validate:"required" example:"password"`
}

type oAuth2SignInUserResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary     oAuth2 sign in
// @Description sign in
// @Tags        oAuth2
// @Accept      json
// @Produce     json
// @Param       input   body     oAuth2SignInUserModel true "account info"
// @Success     200     {object} oAuth2SignInUserResult
// @Failure     400     {object} errors.HTTPError
// @Failure     404     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /serviceoauth2/sign-in [post]
func (h *Handler) oAuth2SignIn(c echo.Context) error {
	input := new(oAuth2SignInUserModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := &serviceoauth2.SignInRequest{
		Username: input.Username,
		Password: input.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.services.OAuth2Service.SingInRequestHandler.Handle(ctx, request)
	if err != nil {
		if errors.Is(err, serviceoauth2.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "user not found")
		} else if errors.Is(err, serviceoauth2.ErrInvalidPassword) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user password")
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, oAuth2SignInUserResult{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}

type oAuth2LogOutModel struct {
	AccessToken string `json:"access_token" validate:"required"`
}

// @Summary     oAuth2 sign out
// @Description sign out
// @Tags        oAuth2
// @Accept      json
// @Produce     json
// @Param       input body oAuth2LogOutModel true "account info"
// @Success     204
// @Failure     400     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /serviceoauth2/sign-out [post]
func (h *Handler) oAuth2SignOut(c echo.Context) error {
	input := new(oAuth2LogOutModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := &serviceoauth2.SignOutRequest{
		AccessToken: input.AccessToken,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := h.services.OAuth2Service.SignOutRequestHandler.Handle(ctx, request)
	if err != nil {
		if errors.Is(err, serviceoauth2.ErrInvalidOrExpiredToken) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired token")
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

type oAuth2RefreshModel struct {
	RefreshToken string `json:"refresh_token" validator:"required"`
}

type oAuth2RefreshResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary     oAuth2 refresh
// @Description refresh
// @Tags        oAuth2
// @Accept      json
// @Produce     json
// @Param       input   body     oAuth2RefreshModel true "account info"
// @Success     200     {object} oAuth2RefreshResult
// @Failure     400     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /serviceoauth2/refresh [post]
func (h *Handler) oAuth2Refresh(c echo.Context) error {
	input := new(oAuth2RefreshModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := &serviceoauth2.RefreshRequest{
		RefreshToken: input.RefreshToken,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.services.OAuth2Service.RefreshRequestHandler.Handle(ctx, request)
	if err != nil {
		if errors.Is(err, serviceoauth2.ErrInvalidOrExpiredToken) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid or expired token")
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, oAuth2RefreshResult{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	})
}
