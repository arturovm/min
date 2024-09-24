package min

import "net/http"

// Middleware is a type alias to a function that takes a handler and returns
// another handler.
type Middleware func(http.Handler) http.Handler

// Then composes middleware m with middleware mw, returning a Middleware that
// first resolves m and then mw.
func (m Middleware) Then(mw Middleware) Middleware {
	if mw == nil {
		return m
	}
	return func(h http.Handler) http.Handler {
		return m(mw(h))
	}
}
func connect(handler http.Handler, mw Middleware) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Run handler and gather effects. In a future implementation, we
		// could catch writes to ResponseWriter by passing a bridge type
		// here with an internal buffer and then dumping that buffer onto
		// ResponseWritter.
		handler.ServeHTTP(w, r)
		// Commit effects. This noop is simply the chain's end point.
		mw(http.HandlerFunc(noop)).ServeHTTP(w, r)
	})
}

func noop(_ http.ResponseWriter, _ *http.Request) {}
