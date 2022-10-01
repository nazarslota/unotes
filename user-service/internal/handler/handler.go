package handler

import (
	"io"
	"net/http"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/udholdenhed/unotes/user-service/internal/service"

	swagger "github.com/swaggo/echo-swagger"
	_ "github.com/udholdenhed/unotes/user-service/docs"
)

type Handler struct {
	service *service.Service
	logger  *zerolog.Logger
}

func NewHandler(service *service.Service, logger *zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// @title       User Service
// @version     1.0
// @description User Service - service for saving user data, developed for UNotes(notes system).

// @host     localhost:8082
// @BasePath /api/user

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
	api := e.Group("/api/user")
	{
		api.POST("/create", h.createUser)
		api.GET("/i/:id", h.findUserByID)
		api.GET("/u/:username", h.findUserByUsername)
	}

	return e
}
