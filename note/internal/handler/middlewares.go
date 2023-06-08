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

type corsMiddleware struct{}

type corsMiddlewareOptions struct{}

func newCORSMiddleware(_ corsMiddlewareOptions) *corsMiddleware {
	return &corsMiddleware{}
}

func (m *corsMiddleware) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.allowedOrigin(r.Header.Get("Origin")) {
			//w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		}

		if r.Method == "OPTIONS" {
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func (m *corsMiddleware) allowedOrigin(_ string) bool { return true }
