package min

import (
	"net/http"
	_path "path"
)

// Groupable enables route composition
type Groupable interface {
	Group(string) Router
}

// group is a concrete route group
type group struct {
	root  string
	super Router
	chain func(http.Handler) http.Handler
}

// Group defines a route group on top of the current group
func (g *group) Group(root string) Router {
	return &group{
		root:  root,
		super: g,
	}
}

// Use adds a middleware to the group's middleware chain
func (g *group) Use(f func(http.Handler) http.Handler) {
	if g.chain == nil {
		g.chain = f
		return
	}
	g.chain = compose(f, g.chain)
}

// Post registers a handler for POST requests
func (g *group) Post(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Post(_path.Join(g.root, path), handler)
}

// Get registers a handler for GET requests
func (g *group) Get(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Get(_path.Join(g.root, path), handler)
}

// Put registers a handler for PUT requests
func (g *group) Put(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Put(_path.Join(g.root, path), handler)
}

// Patch registers a handler for PATCH requests
func (g *group) Patch(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Patch(_path.Join(g.root, path), handler)
}

// Delete registers a handler for DELETE requests
func (g *group) Delete(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Delete(_path.Join(g.root, path), handler)
}

// Head registers a handler for HEAD requests
func (g *group) Head(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Head(_path.Join(g.root, path), handler)
}

// Options registers a handler for OPTIONS requests
func (g *group) Options(path string, handler http.HandlerFunc) {
	if g.chain != nil {
		handler = g.chain(handler).ServeHTTP
	}
	g.super.Options(_path.Join(g.root, path), handler)
}

func (g *group) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.super.ServeHTTP(w, req)
}
