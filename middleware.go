package min

import "net/http"

// Pluggable allows for middleware usage
type Pluggable interface {
	Use(func(http.Handler) http.Handler)
}

func compose(f, g func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return g(f(h))
	}
}
