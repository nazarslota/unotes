package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/udholdenhed/unotes/auth/internal/service/oauth2"
)

func (h *Handler) oAuth2SignUp(c echo.Context) error {
	type SignUpUserModel struct {
		Username string `json:"username" validate:"required,min=4,max=32"`
		Password string `json:"password" validate:"required,min=8,max=64"`
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.services.OAuth2Service.SignUpRequestHandler.Handler(ctx, request); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) oAuth2SignIn(c echo.Context) error {
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

	request := oauth2.SignInRequest{
		Username: input.Username,
		Password: input.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.services.OAuth2Service.SingInRequestHandler.Handle(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

func (h *Handler) oAuth2SignOut(c echo.Context) error {
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

	request := oauth2.LogOutRequest{
		AccessToken: input.AccessToken,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.services.OAuth2Service.SignOutRequestHandler.Handle(ctx, request); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) oAuth2Refresh(c echo.Context) error {
	type RefreshModel struct {
		RefreshToken string `json:"refresh_token" validator:"required"`
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.services.OAuth2Service.RefreshRequestHandler.Handle(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}
