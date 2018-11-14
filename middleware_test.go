package min_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/arturovm/min"
)

type contextKey uint8

const (
	insertorKey contextKey = iota + 1
	userIDKey
	postIDKey
)

var _ = Describe("Middleware", func() {
	insertor := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), insertorKey, "inserted value")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	insertor2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), insertorKey, "inserted value 2")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	insertor3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), insertorKey, "inserted value 3")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	userExtractor := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := min.GetParam(r, "userID")
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	postExtractor := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			postID := min.GetParam(r, "postID")
			ctx := context.WithValue(r.Context(), postIDKey, postID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	postHandler := func(w http.ResponseWriter, r *http.Request) {
		inserted := r.Context().Value(insertorKey).(string)
		userID := r.Context().Value(userIDKey).(string)
		postID := r.Context().Value(postIDKey).(string)
		w.Header().Set("x-inserted", inserted)
		w.Header().Set("x-user-id", userID)
		w.Header().Set("x-post-id", postID)
		values := map[string]string{
			"inserted": inserted,
			"user_id":  userID,
			"post_id":  postID,
		}
		enc := json.NewEncoder(w)
		enc.Encode(values)
	}

	insertorHandler := func(w http.ResponseWriter, r *http.Request) {
		inserted := r.Context().Value(insertorKey).(string)
		fmt.Fprint(w, inserted)
	}

	Given("a router with middleware at the root and middleware in a group", func() {
		r := min.New(nil)

		r.Use(insertor)

		user := r.Group("/user/:userID")
		user.Use(userExtractor)

		post := user.Group("/post/:postID")
		post.Use(postExtractor)
		post.Get("/", postHandler)
		post.Post("/", postHandler)
		post.Put("/", postHandler)
		post.Patch("/", postHandler)
		post.Delete("/", postHandler)
		post.Head("/", postHandler)
		post.Options("/", postHandler)

		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a GET request to the route /user/1/post/1", func() {
			resp, err := http.Get(ts.URL + "/user/1/post/1")
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				var values map[string]string
				dec := json.NewDecoder(resp.Body)
				dec.Decode(&values)
				Expect(values["inserted"]).To(Equal("inserted value"))
				Expect(values["user_id"]).To(Equal("1"))
				Expect(values["post_id"]).To(Equal("1"))
			})
		})

		When("I make a POST request to the route /user/1/post/2", func() {
			resp, err := http.Post(ts.URL+"/user/1/post/2", "text/plain", nil)
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				var values map[string]string
				dec := json.NewDecoder(resp.Body)
				dec.Decode(&values)
				Expect(values["inserted"]).To(Equal("inserted value"))
				Expect(values["user_id"]).To(Equal("1"))
				Expect(values["post_id"]).To(Equal("2"))
			})
		})

		When("I make a PUT request to the route /user/2/post/3", func() {
			req, _ := http.NewRequest(http.MethodPut, ts.URL+"/user/2/post/3", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				var values map[string]string
				dec := json.NewDecoder(resp.Body)
				dec.Decode(&values)
				Expect(values["inserted"]).To(Equal("inserted value"))
				Expect(values["user_id"]).To(Equal("2"))
				Expect(values["post_id"]).To(Equal("3"))
			})
		})

		When("I make a PATCH request to the route /user/3/post/test", func() {
			req, _ := http.NewRequest(http.MethodPatch, ts.URL+"/user/3/post/test", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				var values map[string]string
				dec := json.NewDecoder(resp.Body)
				dec.Decode(&values)
				Expect(values["inserted"]).To(Equal("inserted value"))
				Expect(values["user_id"]).To(Equal("3"))
				Expect(values["post_id"]).To(Equal("test"))
			})
		})

		When("I make a DELETE request to the route /user/3/post/test2", func() {
			req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/user/3/post/test2", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				var values map[string]string
				dec := json.NewDecoder(resp.Body)
				dec.Decode(&values)
				Expect(values["inserted"]).To(Equal("inserted value"))
				Expect(values["user_id"]).To(Equal("3"))
				Expect(values["post_id"]).To(Equal("test2"))
			})
		})

		When("I make a HEAD request to the route /user/1/post/head", func() {
			resp, err := http.Head(ts.URL + "/user/1/post/head")
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Header.Get("x-inserted")).To(Equal("inserted value"))
				Expect(resp.Header.Get("x-user-id")).To(Equal("1"))
				Expect(resp.Header.Get("x-post-id")).To(Equal("head"))
			})
		})

		When("I make an OPTIONS request to the route /user/3/post/options", func() {
			req, _ := http.NewRequest(http.MethodOptions, ts.URL+"/user/3/post/options", nil)
			resp, err := http.DefaultClient.Do(req)
			Then("an HTTP response with the correct parameters should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				var values map[string]string
				dec := json.NewDecoder(resp.Body)
				dec.Decode(&values)
				Expect(values["inserted"]).To(Equal("inserted value"))
				Expect(values["user_id"]).To(Equal("3"))
				Expect(values["post_id"]).To(Equal("options"))
			})
		})
	})

	Given("a router with middleware at the root and middleware in groups", func() {
		r := min.New(nil)

		r.Use(insertor)
		r.Use(insertor2)
		r.Get("/", insertorHandler)

		widgets := r.Group("/widgets")
		widgets.Use(insertor)
		widgets.Use(insertor3)
		widgets.Get("/", insertorHandler)

		sprockets := r.Group("/sprockets")
		sprockets.Get("/", insertorHandler)

		ts := httptest.NewServer(r)
		defer ts.Close()
		When("I make a GET request to the route /", func() {
			resp, err := http.Get(ts.URL)
			Then("an HTTP response with a body with the value 'inserted value 2 should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("inserted value 2"))
			})
		})

		When("I make a GET request to the route /widgets", func() {
			resp, err := http.Get(ts.URL + "/widgets")
			Then("an HTTP response with a body with the value 'inserted value 3 should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("inserted value 3"))
			})
		})

		When("I make a GET request to the route /sprockets", func() {
			resp, err := http.Get(ts.URL + "/sprockets")
			Then("an HTTP response with a body with the value 'inserted value 2 should be returned", func() {
				Expect(err).ToNot(HaveOccurred())
				b, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(b)).To(Equal("inserted value 2"))
			})
		})
	})
})
