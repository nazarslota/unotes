package rest

import (
	"io"
	"net/http"
	"strings"
	"time"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nazarslota/unotes/auth/internal/service"
	"github.com/nazarslota/unotes/auth/pkg/errors"

	_ "github.com/nazarslota/unotes/auth/api/swagger"
	swagger "github.com/swaggo/echo-swagger"
)

type Handler struct {
	address  string
	logger   Logger
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
	router := echo.New()

	router.Logger.SetOutput(io.Discard)
	router.StdLogger.SetOutput(io.Discard)

	router.Validator = newValidator(valid.New())
	router.HTTPErrorHandler = newHTTPErrorHandler(router)

	router.Use(newLoggerMiddleware())
	router.Use(newRequestLoggerMiddleware(h.logger))
	router.Use(newCORSMiddleware())

	// router.Debug = true

	router.GET("/swagger/*", swagger.WrapHandler)
	api := router.Group("/api/auth")
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
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes, // 1 MB
	}
	return &Server{server: s}
}

// Validator.

type validator struct {
	validator *valid.Validate
}

func newValidator(v *valid.Validate) *validator {
	return &validator{validator: v}
}

func (v *validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Error handler.

func newHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
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
				"message": strings.ToLower(http.StatusText(http.StatusInternalServerError)),
			}
			if e.Debug {
				res["debug"] = err.Error()
			}
		}

		switch c.Request().Method {
		case http.MethodHead:
			err = c.NoContent(res["code"].(int))
		default:
			err = c.JSON(res["code"].(int), echo.Map{
				"error": res,
			})
		}

		if err != nil {
			e.Logger.Error(err)
		}
	}
}

// Logger.

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

func newLoggerMiddleware() echo.MiddlewareFunc {
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

// CORS.

func newCORSMiddleware() echo.MiddlewareFunc {
	return middleware.CORS()
}
