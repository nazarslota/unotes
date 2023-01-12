package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func newCORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORS()
}
