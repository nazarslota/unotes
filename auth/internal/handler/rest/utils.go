package rest

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nazarslota/unotes/auth/pkg/errors"
)

type validate struct {
	validate *validator.Validate
}

func newValidate(v *validator.Validate) *validate {
	return &validate{validate: v}
}

func (v *validate) Validate(i any) error {
	if err := v.validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func newHTTPErrorHandler(e *echo.Echo, logger Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		res := echo.Map{}
		switch err := err.(type) {
		case *echo.HTTPError:
			res = echo.Map{"code": err.Code, "message": err.Message}
			if s, ok := err.Message.(string); ok {
				res["message"] = s
				if e.Debug && err.Internal != nil {
					res["debug"] = err.Internal.Error()
				}
			}
		case *errors.HTTPError:
			res = echo.Map{"code": err.Code, "message": err.Message}
			if s, ok := err.Message.(string); ok {
				res["message"] = s
				if e.Debug && err.Internal != nil {
					res["debug"] = err.Internal.Error()
				}
			}
		default:
			res = echo.Map{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			}
			if e.Debug {
				res["debug"] = err.Error()
			}
		}

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(res["code"].(int))
		} else {
			err = c.JSON(res["code"].(int), echo.Map{"error": res})
		}

		if err != nil {
			logger.WarnFields("An error occurred while processing a HTTP error.", map[string]any{"error": err})
		}
	}
}

type Logger interface {
	InfoFields(msg string, fields map[string]any)
	WarnFields(msg string, fields map[string]any)
}

func newLoggerMiddleware(logger Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req, res := c.Request(), c.Response()
			fields := map[string]any{
				"ip":         req.RemoteAddr,
				"method":     req.Method,
				"uri":        req.RequestURI,
				"user_agent": req.UserAgent(),
			}

			s := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}
			dur := time.Since(s)

			fields["status"] = res.Status
			fields["duration"] = dur.String()
			logger.InfoFields("REST, request handled.", fields)

			return nil
		}
	}
}

func newRequestLoggerMiddleware(_ Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:       func(_ echo.Context) bool { return true },
		LogValuesFunc: func(_ echo.Context, _ middleware.RequestLoggerValues) error { return nil },
	})
}

func newCORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Link"},
		MaxAge:           300,
	})
}
