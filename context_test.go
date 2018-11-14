package min_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/arturovm/min"
)

var _ = Describe("Context", func() {
	paramHandler := func(w http.ResponseWriter, r *http.Request) {
		val := min.GetParam(r, "param")
		fmt.Fprint(w, val)
	}

	Given("a vanilla router with a parameter handler for GET requests on the router's root", func() {
		s := http.NewServeMux()
		s.HandleFunc("/", paramHandler)
		ts := httptest.NewServer(s)
		defer ts.Close()
		When("I make a GET request to the router's root", func() {
			resp, err := http.Get(ts.URL)
			Then("an HTTP response with a body containing an empty string should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal(""))
			})
		})
	})
})
