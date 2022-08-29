package oauth2

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/udholdenhed/unotes/auth/internal/application"
	"github.com/udholdenhed/unotes/auth/internal/application/oauth2"
)

type Handler struct {
	AuthService *application.OAuth2Service
}

func NewHandler(service *application.OAuth2Service) *Handler {
	return &Handler{AuthService: service}
}

func (h *Handler) SignIn(c echo.Context) error {
	type SignInUserModel struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	input := new(SignInUserModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := oauth2.SignInRequest{Username: input.Username, Password: input.Password}
	result, err := h.AuthService.SingInRequestHandler.Handle(request)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

func (h *Handler) SignOut(c echo.Context) error {
	type LogOutModel struct {
		AccessToken string `json:"access_token" validate:"required"`
	}

	input := new(LogOutModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := oauth2.LogOutRequest{AccessToken: input.AccessToken}
	if err := h.AuthService.SignOutRequestHandler.Handle(request); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) SignUp(c echo.Context) error {
	type SignUpUserModel struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	input := new(SignUpUserModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := oauth2.SignUpRequest{
		Username: input.Username,
		Password: input.Password,
	}

	if err := h.AuthService.SignUpRequestHandler.Handler(request); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Refresh(c echo.Context) error {
	type RefreshModel struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	input := new(RefreshModel)
	if err := c.Bind(input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := oauth2.RefreshRequest{
		RefreshToken: input.RefreshToken,
	}

	result, err := h.AuthService.RefreshRequestHandler.Handle(request)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}
