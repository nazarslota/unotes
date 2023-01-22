package grpc

import "github.com/nazarslota/unotes/note/internal/service"

type HandlerOption func(h *Handler)

func WithService(services service.Services) HandlerOption {
	return func(h *Handler) {
		h.service = services
	}
}

func WithAddress(address string) HandlerOption {
	return func(h *Handler) {
		h.address = address
	}
}
