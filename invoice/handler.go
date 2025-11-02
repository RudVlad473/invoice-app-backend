package invoice

import "net/http"

type Handler struct {
	http.Handler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var methodHandler http.HandlerFunc

	switch r.Method {
	case http.MethodGet:
		methodHandler = h.getAll
	}

	methodHandler(w, r)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {

}
