package handler

import (
	"net/http"
	"time"
)

type RESTLogger interface {
	InfoFields(msg string, fields map[string]any)
}

type loggerMiddleware struct {
	Logger RESTLogger
}

type loggerMiddlewareOptions struct {
	Logger RESTLogger
}

func newLoggerMiddleware(options loggerMiddlewareOptions) *loggerMiddleware {
	return &loggerMiddleware{Logger: options.Logger}
}

func (m *loggerMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(start)

		fields := map[string]any{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
		}
		m.Logger.InfoFields("HTTP request handled.", fields)
	})
}
