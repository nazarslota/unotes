package handler

import (
	"net/http"
	"time"
)

type restLoggerMiddleware struct {
	Logger Logger
}

func newRESTLoggerMiddleware(logger Logger) *restLoggerMiddleware {
	return &restLoggerMiddleware{Logger: logger}
}

func (m *restLoggerMiddleware) Middleware(handler http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
		}
		logger.InfoFields("HTTP request handled.", fields)
	})
}
