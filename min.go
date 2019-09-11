package min

import "net/http"

// Min is this package's main type. It contains the root route group.
type Min struct {
	*Group
}

// New takes a Handler and initializes a new Min instance with a root route
// group.
func New(handler Handler) *Min {
	return &Min{
		Group: &Group{Path: "/", handler: handler},
	}
}

func (m *Min) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r)
}
