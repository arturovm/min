package min

import "net/http"

// Handler represents anything that can be used to register routes
type Handler interface {
	Post(string, http.HandlerFunc)
	Get(string, http.HandlerFunc)
	Put(string, http.HandlerFunc)
	Patch(string, http.HandlerFunc)
	Delete(string, http.HandlerFunc)
	Head(string, http.HandlerFunc)
	Options(string, http.HandlerFunc)
}
