package min

import "net/http"

// Router implements the groupable, handler, pluggable and http.Handler interfaces
type Router interface {
	Groupable
	Handler
	Pluggable
	http.Handler
}
