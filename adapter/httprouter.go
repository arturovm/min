package adapter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Httprouter struct {
	Router *httprouter.Router
}

func (h *Httprouter) Handle(method, path string, handler http.Handler) {
	h.Router.Handler(method, path, handler)
}

func (h *Httprouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}
