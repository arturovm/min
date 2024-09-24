package min

import (
	"net/http"
	"path"
)

// Group represents a route group that shares a middleware chain and a common
// path.
type Group struct {
	Path       string
	handler    Handler
	parent     *Group
	entryChain Middleware
	exitChain  Middleware
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

// Use sets this group's middleware chain. Each call to Use appends to the
// chain.
func (g *Group) Use(m Middleware) {
	g.Entry(m)
}

// Entry sets this group's entry middleware chain. Entry middleware is executed
// from the start of the request and hands off to the handler. Each call to Entry
// appends to the entry chain.
func (g *Group) Entry(m Middleware) {
	if g.entryChain == nil {
		g.entryChain = m
		return
	}
	g.entryChain = g.entryChain.Then(m)
}

// Exit sets this group's exit middleware chain. Exit middleware is executed
// from after the handler has returned, up until the response is written. Each
// call to Exit appends to the exit chain.
func (g *Group) Exit(m Middleware) {
	if g.exitChain == nil {
		g.exitChain = m
		return
	}
	g.exitChain = g.exitChain.Then(m)
}

func (g *Group) handle(method, relativePath string, handler http.Handler) {
	entryChain := g.fullEntryChain()
	if entryChain != nil {
		handler = entryChain(handler)
	}

	exitChain := g.fullExitChain()
	if exitChain != nil {
		handler = connect(handler, exitChain)
	}
	g.handler.Handle(method, path.Join(g.FullPath(), relativePath), handler)
}

func (g *Group) fullEntryChain() Middleware {
	if g.parent == nil && g.entryChain == nil {
		return nil
	}
	if g.parent == nil {
		return g.entryChain
	}
	parentChain := g.parent.fullEntryChain()
	if parentChain == nil {
		return g.entryChain
	}
	return parentChain.Then(g.entryChain)
}

func (g *Group) fullExitChain() Middleware {
	if g.parent == nil && g.exitChain == nil {
		return nil
	}
	if g.parent == nil {
		return g.exitChain
	}
	parentChain := g.parent.fullExitChain()
	if g.exitChain == nil {
		return parentChain
	}
	return g.exitChain.Then(parentChain)
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
