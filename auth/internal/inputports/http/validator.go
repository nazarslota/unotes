package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
)

type Validator struct {
	Validator *validator.Validate
}

var _ echo.Validator = (*Validator)(nil)

func NewValidator(validator *validator.Validate) *Validator {
	return &Validator{Validator: validator}
}

func (v *Validator) Validate(i any) error {
	if err := v.Validator.Struct(i); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
