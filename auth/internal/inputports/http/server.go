package http

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/udholdenhed/unotes/auth/internal/application"
	"github.com/udholdenhed/unotes/auth/internal/inputports/http/oauth2"
)

type Server struct {
	E           *echo.Echo
	Address     string
	AppServices *application.Services
}

type ServerOptions struct {
	Address string
	Debug   bool
}

func NewServer(services *application.Services, options ServerOptions) *Server {
	s := &Server{
		E:           echo.New(),
		Address:     options.Address,
		AppServices: services,
	}
	s.E.HTTPErrorHandler = NewErrorHandler(s.E)
	s.E.Validator = NewValidator(validator.New())
	s.E.Debug = options.Debug

	s.E.Use(LoggerMiddleware())
	s.E.Use(RecoverMiddleware())
	s.E.Use(CORSMiddleware())

	s.RegisterOAuth2Routes()
	return s
}

func (s *Server) RegisterOAuth2Routes() {
	oa2g := s.E.Group("api/v1/oauth2")
	{
		oa2g.POST("/sign-up", oauth2.NewHandler(s.AppServices.AuthService).SignUp)
		oa2g.POST("/sign-in", oauth2.NewHandler(s.AppServices.AuthService).SignIn)
		oa2g.POST("/sign-out", oauth2.NewHandler(s.AppServices.AuthService).SignOut)
		oa2g.POST("/refresh", oauth2.NewHandler(s.AppServices.AuthService).Refresh)
	}
}

func (s *Server) Run() error {
	return s.E.Start(s.Address)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.E.Shutdown(ctx)
}

func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.Logger()
}

func RecoverMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}

func CORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORS()
}
