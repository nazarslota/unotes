package handler

import (
	"io"
	"net/http"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/udholdenhed/unotes/auth-service/internal/service"

	swagger "github.com/swaggo/echo-swagger"
	_ "github.com/udholdenhed/unotes/auth-service/docs"
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

// @title       Auth Service
// @version     1.0
// @description Authentication service, developed for UNotes(notes system).

// @host     localhost:8081
// @BasePath /api/auth

func (h *Handler) H() http.Handler {
	e := echo.New()

	e.Logger.SetOutput(io.Discard)
	e.StdLogger.SetOutput(io.Discard)

	e.Validator = newValidator(valid.New())
	e.HTTPErrorHandler = newHTTPErrorHandler(e)

	e.Use(newLoggerMiddleware())
	e.Use(newRequestLoggerMiddleware(h.logger))
	e.Use(newCORSMiddleware())

	// e.Debug = true

	e.GET("/swagger/*", swagger.WrapHandler)
	api := e.Group("/api/auth")
	{
		oAuth2 := api.Group("/oauth2")
		{
			oAuth2.POST("/sign-up", h.oAuth2SignUp)
			oAuth2.POST("/sign-in", h.oAuth2SignIn)
			oAuth2.POST("/sign-out", h.oAuth2SignOut)
			oAuth2.POST("/refresh", h.oAuth2Refresh)
		}
	}

	return e
}
