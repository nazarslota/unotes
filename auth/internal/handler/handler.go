package handler

import (
	"io"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/udholdenhed/unotes/auth/internal/service"
)

type Handler struct {
	services *service.Service
	logger   *zerolog.Logger
}

func NewHandler(services *service.Service, logger *zerolog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()

	router.Logger.SetOutput(io.Discard)
	router.StdLogger.SetOutput(io.Discard)

	router.Validator = newValidator(valid.New())
	router.HTTPErrorHandler = newHTTPErrorHandler(router)

	router.Use(newLoggerMiddleware())
	router.Use(newRequestLoggerMiddleware(h.logger))

	router.Debug = true

	api := router.Group("/api")
	{
		oAuth2 := api.Group("/oauth2")
		{
			oAuth2.POST("/sign-up", h.oAuth2SignUp)
			oAuth2.POST("/sign-in", h.oAuth2SignIn)
			oAuth2.POST("/sign-out", h.oAuth2SignOut)
			oAuth2.POST("/refresh", h.oAuth2Refresh)
		}
	}

	return router
}
