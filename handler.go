package min

import "net/http"

// Handler represents a type that can register handlers for a given HTTP verb
// and path.
type Handler interface {
	http.Handler
	Handle(method, path string, handler http.Handler)
}
