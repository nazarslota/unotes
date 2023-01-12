package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

func newLoggerMiddleware(_ Logger) echo.MiddlewareFunc {
	return middleware.Logger()
}

func newRequestLoggerMiddleware(logger Logger) echo.MiddlewareFunc {
	logValuesFunc := func(c echo.Context, v middleware.RequestLoggerValues) error {
		logger.Infof("REST: \"%s %s %s\" %s", v.Method, v.URI, v.Protocol, http.StatusText(v.Status))
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
