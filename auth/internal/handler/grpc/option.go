package grpc

import "github.com/nazarslota/unotes/auth/internal/service"

type HandlerOption func(h *Handler)

func WithAddress(address string) HandlerOption {
	return func(h *Handler) {
		h.address = address
	}
}

func WithService(services *service.Service) HandlerOption {
	return func(h *Handler) {
		h.services = services
	}
}

func WithLogger(logger Logger) HandlerOption {
	return func(h *Handler) {
		h.logger = logger
	}
}
