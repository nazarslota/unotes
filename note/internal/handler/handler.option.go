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

func WithGRPCLogger(logger GRPCLogger) Option {
	return func(h *Handler) {
		h.grpcLogger = logger
	}
}

func WithRESTLogger(logger RESTLogger) Option {
	return func(h *Handler) {
		h.restLogger = logger
	}
}
