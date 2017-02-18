package min

import (
	"net/http"
	_path "path"

	"github.com/julienschmidt/httprouter"
)

// Min is the main type. It contains an httprouter.Router, and wraps around it
// to provide extra functionality.
type Min struct {
	router *httprouter.Router
	root   string
	chain  func(http.Handler) http.Handler
}

// New returns a new Min with the provided httprouter.Router as the base router.
// If router is nil, a new, default one is created.
func New(router *httprouter.Router) *Min {
	if router == nil {
		router = httprouter.New()
	}
	return &Min{
		router: router,
		root:   "/",
	}
}

// GetBaseRouter returns a pointer to the contained httprouter.Router, which enables
// users to change and configure the router as they please
func (m *Min) GetBaseRouter() *httprouter.Router {
	return m.router
}

// Group defines a route group on top of the root router
func (m *Min) Group(root string) Router {
	return &group{
		root:  root,
		super: m,
	}
}

// Use adds a middleware to the global middleware chain
func (m *Min) Use(f func(http.Handler) http.Handler) {
	if m.chain == nil {
		m.chain = f
		return
	}
	m.chain = compose(f, m.chain)
}

func (m *Min) handleWithMethod(method, path string, handler http.HandlerFunc) {
	if m.chain != nil {
		handler = m.chain(handler).ServeHTTP
	}
	m.router.Handle(method, path, wrapHandler(handler))
}

// Post registers a handler for POST requests
func (m *Min) Post(path string, handler http.HandlerFunc) {
	m.handleWithMethod("POST", _path.Join(m.root, path), handler)
}

// Get registers a handler for GET requests
func (m *Min) Get(path string, handler http.HandlerFunc) {
	m.handleWithMethod("GET", _path.Join(m.root, path), handler)
}

// Put registers a handler for PUT requests
func (m *Min) Put(path string, handler http.HandlerFunc) {
	m.handleWithMethod("PUT", _path.Join(m.root, path), handler)
}

// Patch registers a handler for PATCH requests
func (m *Min) Patch(path string, handler http.HandlerFunc) {
	m.handleWithMethod("PATCH", _path.Join(m.root, path), handler)
}

// Delete registers a handler for DELETE requests
func (m *Min) Delete(path string, handler http.HandlerFunc) {
	m.handleWithMethod("DELETE", _path.Join(m.root, path), handler)
}

// Head registers a handler for HEAD requests
func (m *Min) Head(path string, handler http.HandlerFunc) {
	m.handleWithMethod("HEAD", _path.Join(m.root, path), handler)
}

// Options registers a handler for OPTIONS requests
func (m *Min) Options(path string, handler http.HandlerFunc) {
	m.handleWithMethod("OPTIONS", _path.Join(m.root, path), handler)
}

func (m *Min) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.router.ServeHTTP(w, req)
}
