package min

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// key is an alias of int, for use with request contexts
type key int

// paramsKey is the key that identifies httprouter.Params in a request's
// context
var paramsKey key

func wrapHandler(f http.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, paramsKey, params)
		f.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetParam gets a parameter from a given request
func GetParam(r *http.Request, key string) string {
	if val := r.Context().Value(paramsKey); val != nil {
		if p, ok := val.(httprouter.Params); ok {
			return p.ByName(key)
		}
	}
	return ""
}
