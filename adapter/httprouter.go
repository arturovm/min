package adapter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Httprouter is an adapter around httprouter.Router that implements
// min.Handler.
type Httprouter struct {
	Router *httprouter.Router
}

// Handle implements Handler.Handle.
func (h *Httprouter) Handle(method, path string, handler http.Handler) {
	h.Router.Handler(method, path, handler)
}

func (h *Httprouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}
