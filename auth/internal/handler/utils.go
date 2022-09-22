package handler

import (
	"net/http"
	"strings"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/udholdenhed/unotes/auth/pkg/errors"
)

type validator struct {
	validator *valid.Validate
}

func newValidator(v *valid.Validate) *validator {
	return &validator{validator: v}
}

func (v *validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func newHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		res := echo.Map{}
		switch err := err.(type) {
		case *echo.HTTPError:
			res = echo.Map{"code": err.Code, "message": err.Message}
			if s, ok := err.Message.(string); ok {
				res["message"] = strings.ToLower(s)
				if e.Debug && err.Internal != nil {
					res["debug"] = strings.ToLower(err.Internal.Error())
				}
			}
		case *errors.HTTPError:
			res = echo.Map{"code": err.Code, "message": err.Message}
			if s, ok := err.Message.(string); ok {
				res["message"] = strings.ToLower(s)
				if e.Debug && err.Internal != nil {
					res["debug"] = strings.ToLower(err.Internal.Error())
				}
			}
		default:
			res = echo.Map{
				"code":    http.StatusInternalServerError,
				"message": strings.ToLower(http.StatusText(http.StatusInternalServerError)),
			}
			if e.Debug {
				res["debug"] = err.Error()
			}
		}

		switch c.Request().Method {
		case http.MethodHead:
			err = c.NoContent(res["code"].(int))
		default:
			err = c.JSON(res["code"].(int), echo.Map{
				"error": res,
			})
		}

		if err != nil {
			e.Logger.Error(err)
		}
	}
}

func newLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.Logger()
}

func newRequestLoggerMiddleware(logger *zerolog.Logger) echo.MiddlewareFunc {
	logValuesFunc := func(c echo.Context, v middleware.RequestLoggerValues) error {
		logger.Info().Msgf("\"%s %s %s\" %d", v.Method, v.URI, v.Protocol, v.Status)
		return nil
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: logValuesFunc,
		LogLatency:    true,
		LogProtocol:   true,
		LogRemoteIP:   true,
		LogHost:       true,
		LogMethod:     true,
		LogURI:        true,
		LogStatus:     true,
	})
}
