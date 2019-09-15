package min_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"github.com/arturovm/min"
	"github.com/arturovm/min/adapter"
)

func TestCreateMiddleware(t *testing.T) {
	mw := min.Middleware(func(next http.Handler) http.Handler {
		return nil
	})
	require.NotNil(t, mw)
}

func TestRunMiddleware(t *testing.T) {
	var result string
	mw := min.Middleware(func(next http.Handler) http.Handler {
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
	first := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "Hello, "
			next.ServeHTTP(w, r)
		})
	})
	second := min.Middleware(func(next http.Handler) http.Handler {
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

func TestUseMiddleware(t *testing.T) {
	h := &adapter.Httprouter{Router: httprouter.New()}
	m := min.New(h)

	var count int8
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count++
			next.ServeHTTP(w, r)
		})
	}
	secondMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			count++
			next.ServeHTTP(w, r)
		})
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	})

	m.Use(mw)
	m.Use(secondMw)
	m.Get("/test", handler)
	m.Post("/test", handler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	_, _ = http.Get(ts.URL + "/test")
	_, _ = http.Post(ts.URL+"/test", "text/plain", nil)

	require.Equal(t, int8(6), count)
}
