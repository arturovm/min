package min_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/arturovm/min"
	"github.com/julienschmidt/httprouter"
)

var _ = Describe("Router", func() {
	noContentHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	paramHandler := func(w http.ResponseWriter, r *http.Request) {
		val := min.GetParam(r, "param")
		fmt.Fprint(w, val)
	}

	catchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		val := min.GetParam(r, "catchall")
		fmt.Fprint(w, val)
	}

	Given("a router with a root-level handler for GET requests", func() {
		r := min.New(nil)
		r.Get("/", noContentHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a GET request to the server's root", func() {
			resp, err := http.Get(ts.URL)
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	Given("a router with a handler for POST requests with the route /hello", func() {
		r := min.New(nil)
		r.Post("/hello", noContentHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a POST request to the route /hello", func() {
			resp, err := http.Post(ts.URL+"/hello", "text/plain", nil)
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	Given("a router with a handler for PUT requests with the route /:param", func() {
		r := min.New(nil)
		r.Put("/:param", paramHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a PUT request to the route /echo", func() {
			req, _ := http.NewRequest(http.MethodPut, ts.URL+"/echo", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a body containing 'echo' should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("echo"))
			})
		})
	})

	Given("a router with a handler for PATCH requests with the route /catch/*catchall", func() {
		r := min.New(nil)
		r.Patch("/catch/*catchall", catchAllHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a PATCH request to the route /catch/this/is/echo", func() {
			req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/catch/this/is/echo", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a body containing '/this/is/echo' should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("/this/is/echo"))
			})
		})
	})

	Given("a router with a root-level handler for DELETE requests", func() {
		r := min.New(nil)
		r.Delete("/", noContentHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a DELETE request to the route /resource1", func() {
			req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/resource1", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a 404 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Given("a router with a root-level handler for HEAD requests", func() {
		r := min.New(nil)
		r.Head("/", noContentHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a HEAD request to the server's root", func() {
			resp, err := http.Get(ts.URL)
			Then("an HTTP response with a 405 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
			})
		})
	})

	Given("a router with a root-level handler for OPTIONS requests", func() {
		r := min.New(nil)
		r.Options("/", noContentHandler)
		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make an OPTIONS request to the server's root", func() {
			req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	Given("a router initialized with an existing httprouter.Router", func() {
		rr := httprouter.New()
		r := min.New(rr)
		When("I access the underlying router", func() {
			base := r.GetBaseRouter()
			Then("the underlying router should be the same as the prexisting router", func() {
				Expect(base).To(Equal(rr))
			})
		})
	})
})
