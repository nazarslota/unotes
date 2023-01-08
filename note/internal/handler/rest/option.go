package rest

type HandlerOption func(h *Handler)

func WithAddress(address string) HandlerOption {
	return func(h *Handler) {
		h.address = address
	}
}
