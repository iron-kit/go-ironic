package assistant

type Handler interface {
	SetInstance(Handler)
	getInstance() Handler
}

type BaseHandler struct {
	instance Handler
}

func (h *BaseHandler) getInstance() Handler {
	return h.instance
}

func (h *BaseHandler) SetInstance(instance Handler) {
	h.instance = instance
}
