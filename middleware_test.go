package min_test

import (
	"io"
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

func TestUseMiddlewareWithGroups(t *testing.T) {
	h := &adapter.Httprouter{Router: httprouter.New()}
	m := min.New(h)

	var result string
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "first, "
			next.ServeHTTP(w, r)
		})
	}
	secondMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "second, "
			next.ServeHTTP(w, r)
		})
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result += "handler"
	})

	group := m.NewGroup("/group")
	{
		group.Use(mw)
		emptyGroup := group.NewGroup("/")
		{
			anotherEmptyGroup := emptyGroup.NewGroup("/")
			{
				anotherEmptyGroup.Use(secondMw)
				anotherEmptyGroup.Get("/test", handler)
			}
		}
	}

	ts := httptest.NewServer(m)
	defer ts.Close()

	_, _ = http.Get(ts.URL + "/group/test")

	require.Equal(t, "first, second, handler", result)
}

func TestEntryMiddleware(t *testing.T) {
	var result string
	first := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "Hello, "
			next.ServeHTTP(w, r)
		})
	})
	second := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "world, "
			next.ServeHTTP(w, r)
		})
	})
	mw := first.Then(second)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result += "once more!"
	})

	h := &adapter.Httprouter{Router: httprouter.New()}
	m := min.New(h)

	m.Entry(mw)
	m.Get("/hello", handler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	_, _ = http.Get(ts.URL + "/hello")

	require.Equal(t, "Hello, world, once more!", result)
}

func TestExitMiddleware(t *testing.T) {
	var result string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result += "Goodbye"

		w.Header().Set("x-greeting", "hello :)")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("a little greeting"))
	})
	first := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += ", "
			next.ServeHTTP(w, r)
		})
	})
	second := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "my "
			next.ServeHTTP(w, r)
		})
	})
	third := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "good "
			next.ServeHTTP(w, r)
		})
	})
	fourth := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "friends"
			next.ServeHTTP(w, r)
		})
	})
	fifth := min.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result += "!"
			next.ServeHTTP(w, r)
		})
	})

	h := &adapter.Httprouter{Router: httprouter.New()}
	m := min.New(h)

	m.Exit(fourth)
	m.Exit(fifth)
	group := m.NewGroup("/")
	{
		mw := second.Then(third)
		group.Exit(mw)
		emptyGroup := group.NewGroup("/")
		{
			anotherEmptyGroup := emptyGroup.NewGroup("/")
			{
				anotherEmptyGroup.Exit(first)
				anotherEmptyGroup.Get("/hello", handler)
			}
		}
	}

	ts := httptest.NewServer(m)
	defer ts.Close()

	response, _ := http.Get(ts.URL + "/hello")
	responseBody, _ := io.ReadAll(response.Body)

	require.Equal(t, "Goodbye, my good friends!", result)
	require.Equal(t, http.StatusCreated, response.StatusCode)
	require.Equal(t, "hello :)", response.Header.Get("x-greeting"))
	require.Equal(t, []byte("a little greeting"), responseBody)
}
