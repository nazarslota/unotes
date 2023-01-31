package handler

import (
	"github.com/nazarslota/unotes/note/internal/service"
)

type Option func(h *Handler)

func WithServices(services service.Services) Option {
	return func(h *Handler) {
		h.services = services
	}
}

func WithGRPCServerAddr(addr string) Option {
	return func(h *Handler) {
		h.grpcAddr = addr
	}
}

func WithRESTServerAddr(addr string) Option {
	return func(h *Handler) {
		h.restAddr = addr
	}
}

func WithLogger(logger Logger) Option {
	return func(h *Handler) {
		h.logger = logger
	}
}
