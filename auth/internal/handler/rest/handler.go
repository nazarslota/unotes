package rest

import (
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
	address  string
	services *service.Services

	logger Logger
	debug  bool
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
// @description Authentication service.

// @host     localhost:8081
// @BasePath /api

func (h *Handler) S() *Server {
	mux := echo.New()

	mux.Debug = h.debug

	mux.Logger.SetOutput(io.Discard)
	mux.StdLogger.SetOutput(io.Discard)

	mux.Validator = newValidate(validator.New())
	mux.HTTPErrorHandler = newHTTPErrorHandler(mux, h.logger)

	mux.Use(newLoggerMiddleware(h.logger))
	mux.Use(newRequestLoggerMiddleware(h.logger))
	mux.Use(corsMiddleware())

	api := mux.Group("/api")
	{
		api.GET("/swagger/*", swagger.WrapHandler)
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
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes, // 1 MB
	}

	return &Server{server: s}
}
