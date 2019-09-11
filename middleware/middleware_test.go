package middleware_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/min/middleware"
)

func TestCreateMiddleware(t *testing.T) {
	mw := middleware.Middleware(func(next http.Handler) http.Handler {
		return nil
	})
	require.NotNil(t, mw)
}

func TestUseMiddleware(t *testing.T) {
	var result string
	mw := middleware.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "Hello, "
			next.ServeHTTP(w, r)
		})
	})
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result += "world!"
	})).ServeHTTP(nil, nil)

	require.Equal(t, "Hello, world!", result)
}

func TestComposeMiddleware(t *testing.T) {
	var result string
	first := middleware.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "Hello, "
			next.ServeHTTP(w, r)
		})
	})
	second := middleware.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "world "
			next.ServeHTTP(w, r)
		})
	})
	mw := first.Then(second)
	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result += "again!"
	})).ServeHTTP(nil, nil)

	require.Equal(t, "Hello, world again!", result)
}
