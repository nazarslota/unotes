package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/udholdenhed/unotes/user-service/internal/service/user"
)

type createUserModel struct {
	Username string `json:"username" validate:"required,min=4,max=32" example:"username"`
	Email    string `json:"email" validate:"required,email" example:"example@example.example"`
	Password string `json:"password" validate:"required,min=8,max=64" example:"password"`
}

// @Summary     Create User
// @Description Creates a new user account.
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       Input. body createUserModel true "Account info."
// @Success     204
// @Failure     400     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /create [post]
func (h *Handler) createUser(c echo.Context) error {
	input := new(createUserModel)
	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := user.CreateUserRequest{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := h.service.UserService.CreateUserRequestHandler.Handle(ctx, request); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

type findUserByIDModel struct {
	ID string `json:"id" validate:"required,uuid4" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
}

type findUserByIDResult struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// @Summary     Find User By ID
// @Description Returns the user by his ID.
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       Input.  path     string true "User ID."
// @Success     200     {object} findUserByIDResult
// @Failure     400     {object} errors.HTTPError
// @Failure     404     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /i/{id} [get]
func (h *Handler) findUserByID(c echo.Context) error {
	input := &findUserByIDModel{
		ID: c.Param("id"),
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := user.FindUserByIDRequest{
		ID: input.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.service.UserService.FindUserByIDRequestHandler.Handle(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, findUserByIDResult{
		ID:           result.ID,
		Username:     result.Username,
		Email:        result.Email,
		PasswordHash: result.PasswordHash,
	})
}

type findUserByUsernameModel struct {
	Username string `json:"username" validate:"required" example:"username"`
}

type findUserByUsernameResult struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// @Summary     Find User By Username
// @Description Returns the user by username.
// @Tags        User
// @Accept      json
// @Produce     json
// @Param       Input.  path     string true "User username."
// @Success     200     {object} findUserByUsernameResult
// @Failure     400     {object} errors.HTTPError
// @Failure     404     {object} errors.HTTPError
// @Failure     500     {object} errors.HTTPError
// @Failure     default {object} errors.HTTPError
// @Router      /u/{username} [get]
func (h *Handler) findUserByUsername(c echo.Context) error {
	input := &findUserByUsernameModel{
		Username: c.Param("username"),
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	request := user.FindUserByUsernameRequest{
		Username: input.Username,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := h.service.UserService.FindUserByUsernameRequestHandler.Handle(ctx, request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, findUserByUsernameResult{
		ID:           result.ID,
		Username:     result.Username,
		Email:        result.Email,
		PasswordHash: result.PasswordHash,
	})
}
