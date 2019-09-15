package min

import (
	"net/http"
	"path"
)

// Group represents a route group that shares a middleware chain and a common
// path.
type Group struct {
	Path    string
	handler Handler
	parent  *Group
	chain   Middleware
}

// NewGroup creates a new subgroup of group g.
func (g *Group) NewGroup(path string) *Group {
	return &Group{
		Path:    path,
		handler: g.handler,
		parent:  g,
	}
}

// Parent gets the group's parent in the group tree.
func (g *Group) Parent() *Group {
	return g.parent
}

// FullPath returns the group's full path in the group tree (as opposed to this
// group's sub-path)
func (g *Group) FullPath() string {
	if g.parent == nil {
		return g.Path
	}
	return path.Join(g.parent.FullPath(), g.Path)
}

// Use sets this group's middleware chain. Each call to Use replaces the entire chain.
func (g *Group) Use(m Middleware) {
	g.chain = m
}

func (g *Group) handle(method, relativePath string, handler http.Handler) {
	if g.chain != nil {
		handler = g.chain(handler)
	}
	g.handler.Handle(method, path.Join(g.FullPath(), relativePath), handler)
}

// Get registers a handler for GET requests on the given relative path.
func (g *Group) Get(relativePath string, handler http.Handler) {
	g.handle(http.MethodGet, relativePath, handler)
}

// Post registers a handler for POST requests on the given relative path.
func (g *Group) Post(relativePath string, handler http.Handler) {
	g.handle(http.MethodPost, relativePath, handler)
}

// Put registers a handler for PUT requests on the given relative path.
func (g *Group) Put(relativePath string, handler http.Handler) {
	g.handle(http.MethodPut, relativePath, handler)
}

// Patch registers a handler for PATCH requests on the given relative path.
func (g *Group) Patch(relativePath string, handler http.Handler) {
	g.handle(http.MethodPatch, relativePath, handler)
}

// Delete registers a handler for DELETE requests on the given relative path.
func (g *Group) Delete(relativePath string, handler http.Handler) {
	g.handle(http.MethodDelete, relativePath, handler)
}
