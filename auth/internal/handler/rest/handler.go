package rest

import (
	"io"
	"net/http"
	"time"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/nazarslota/unotes/auth/api/swagger"
	"github.com/nazarslota/unotes/auth/internal/service"
	swagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	address string
	logger  Logger

	services *service.Service
}

func NewHandler(options ...HandlerOption) *Handler {
	h := &Handler{}
	for _, option := range options {
		option(h)
	}
	return h
}

// @title       Auth
// @version     1.0
// @description Authentication service, developed for UNotes(notes system).

// @host     localhost:8081
// @BasePath /api/auth

func (h *Handler) S() *Server {
	e := echo.New()

	e.Logger.SetOutput(io.Discard)
	e.StdLogger.SetOutput(io.Discard)

	e.Validator = newValidator(valid.New())
	e.HTTPErrorHandler = newHTTPErrorHandler(e)

	e.Use(newLoggerMiddleware(h.logger))
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

	s := &http.Server{
		Addr:           h.address,
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes, // 1 MB
	}
	return &Server{server: s}
}
