package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	address string
	logger  Logger
}

func NewHandler(options ...HandlerOption) *Handler {
	h := &Handler{}
	for _, option := range options {
		option(h)
	}
	return h
}

func (h *Handler) S() *Server {
	router := gin.Default()

	router.Use()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	server := &http.Server{
		Addr:           h.address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	return &Server{server: server}
}

// Logger

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}
