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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	})

	m.Use(mw)
	m.Get("/test", handler)
	m.Post("/test", handler)

	ts := httptest.NewServer(m)
	defer ts.Close()

	_, _ = http.Get(ts.URL + "/test")
	_, _ = http.Post(ts.URL+"/test", "text/plain", nil)

	require.Equal(t, int8(4), count)
}
