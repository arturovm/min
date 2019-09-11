package middleware

import "net/http"

// Middleware is a type alias to a function that takes a handler and returns
// another handler.
type Middleware func(http.Handler) http.Handler

// Then composes middleware m with middleware mw, returning a Middleware that
// first resolves m and then mw.
func (m Middleware) Then(mw Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		return m(mw(h))
	}
}
