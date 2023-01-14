package grpc

import "github.com/nazarslota/unotes/auth/internal/service"

type HandlerOption func(h *Handler)

func WithService(services *service.Services) HandlerOption {
	return func(h *Handler) {
		h.services = services
	}
}

func WithAddress(address string) HandlerOption {
	return func(h *Handler) {
		h.address = address
	}
}

func WithLogger(logger Logger) HandlerOption {
	return func(h *Handler) {
		h.logger = logger
	}
}
