package min_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/arturovm/min"
)

var _ = Describe("Route Grouping", func() {
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

	Given("a router with a sub-group at path /api and handlers for each method", func() {
		r := min.New(nil)
		api := r.Group("/api")
		{
			api.Get("/", noContentHandler)
			api.Post("/hello", noContentHandler)
			api.Put("/:param", paramHandler)
			api.Patch("/catch/*catchall", catchAllHandler)
			api.Delete("/resource1", noContentHandler)
			api.Head("/", noContentHandler)
			api.Options("/", noContentHandler)
		}
		ts := httptest.NewServer(api)
		defer ts.Close()
		When("I make a GET request to the route /api", func() {
			resp, err := http.Get(ts.URL + "/api")
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		When("I make a POST request to the route /api/hello", func() {
			resp, err := http.Post(ts.URL+"/api/hello", "text/plain", nil)
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		When("I make a PUT request to the route /api/echo", func() {
			req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/echo", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a body containing 'echo' should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("echo"))
			})
		})

		When("I make a PATCH request to the route /api/catch/this/is/echo", func() {
			req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/api/catch/this/is/echo", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a body containing '/this/is/echo' should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("/this/is/echo"))
			})
		})

		When("I make a DELETE request to the route /api/resource1", func() {
			req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/api/resource1", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		When("I make a HEAD request to the route /api", func() {
			resp, err := http.Head(ts.URL + "/api")
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		When("I make an OPTIONS request to the route /api", func() {
			req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/api", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with a 204 status code should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})
})
