package min

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// contextKey is an alias of string, for use with request contexts
type contextKey string

func wrapHandler(f http.Handler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		for i := range params {
			ctx = context.WithValue(ctx, contextKey(params[i].Key), params[i].Value)
		}
		f.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetParam gets a parameter from a given request
func GetParam(r *http.Request, key string) string {
	if ctx := r.Context(); ctx != nil {
		if val := ctx.Value(contextKey(key)); val != nil {
			if s, ok := val.(string); ok {
				return s
			}
		}
	}
	return ""
}
