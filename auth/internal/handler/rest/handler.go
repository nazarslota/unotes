package rest

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/nazarslota/unotes/auth/api/swagger"
	"github.com/nazarslota/unotes/auth/internal/service"
	swagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	addr   string
	logger Logger
	debug  bool

	services service.Services
}

func NewHandler(options ...HandlerOption) *Handler {
	h := new(Handler)
	for _, option := range options {
		option(h)
	}
	return h
}

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func (h *Handler) Server() Server {
	e := echo.New()

	e.Debug = h.debug

	e.Validator = newValidate(validator.New())
	e.HTTPErrorHandler = newHTTPErrorHandler(e, h.logger)

	e.Logger.SetOutput(io.Discard)
	e.StdLogger.SetOutput(io.Discard)

	e.Use(newLoggerMiddleware(h.logger))
	e.Use(newRequestLoggerMiddleware(h.logger))
	e.Use(newCORSMiddleware())

	h.registerEndpoints(e)

	server := &http.Server{
		Addr:           h.addr,
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	return newServer(h.addr, server)
}

func (h *Handler) registerEndpoints(e *echo.Echo) {
	api := e.Group("/api")
	{
		api.GET("/swagger/*", swagger.WrapHandler)
		oAuth2 := api.Group("/oauth2")
		{
			oAuth2.POST("/sign-up", h.oAuth2SignUp)
			oAuth2.POST("/sign-in", h.oAuth2SignIn)
			oAuth2.POST("/sign-out", h.oAuth2SignOut)
			oAuth2.GET("/refresh", h.oAuth2Refresh)
		}
	}
}
