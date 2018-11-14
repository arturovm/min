# min

v0.1.1

[![GoDoc](https://godoc.org/github.com/arturovm/min?status.svg)](https://godoc.org/github.com/arturovm/min)
[![Go Report Card](https://goreportcard.com/badge/github.com/arturovm/min)](https://goreportcard.com/report/github.com/arturovm/min)
[![Build Status](https://travis-ci.com/arturovm/min.svg?branch=master)](https://travis-ci.com/arturovm/min)
[![Codecov](https://img.shields.io/codecov/c/github/arturovm/min.svg)](https://codecov.io/gh/arturovm/min)
![GitHub](https://img.shields.io/github/license/arturovm/min.svg)

`min` is a BYO\*, minimalistic web framework that harnesses the power of
`context` and [`httprouter`](https://github.com/julienschmidt/httprouter), and
adds some functionality on top—namely, middleware chaining and route grouping.
It's meant to be used on projects large and small that require flexibility, and
varying degrees of custom code and architecture. `min` provides the routing, you
provide the app.

This package takes some inspiration from design decisions in [`chi`](https://github.com/pressly/chi)
and [`gin`](https://github.com/gin-gonic/gin).

## Usage

`min` is designed to be as elegant and as close to "the right way to do things"
as possible. Which means that it doesn't implement a lot of custom types, or
does a lot of magic. It relies heavily on `context` and regular types from
`net/http`.

### Hello World

You can initialize a new instance of the `Min` type with an `httprouter.Router`
you provide, or you can pass `nil` to `min.New` and a new, default router will
be created.

``` go
import (
    "fmt"
    "net/http"

    "github.com/arturovm/min"
)

func main() {
    m := min.New(nil)

    m.Get("/", helloWorld)

    http.ListenAndServe(":8080", m)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world!")
}
```

### Route Parameters

`min` supports all the syntax variations for defining route parameters that
`httprouter` does. You then simply use the helper function `min.GetParam` to
access them.

``` go
import (
    "fmt"
    "net/http"

    "github.com/arturovm/min"
)

func main() {
    m := min.New(nil)

    m.Get("/:name", greet)

    http.ListenAndServe(":8080", m)
}

func greet(w http.ResponseWriter, r *http.Request) {
    name := min.GetParam(r, "name")
    fmt.Fprintf(w, "hello %s!", name)
}
```

### Route Grouping

``` go
import (
    "fmt"
    "net/http"

    "github.com/arturovm/min"
)

func main() {
    m := min.New(nil)

    apiRouter := m.Group("/api")
    {
        // GET /api
        apiRouter.Get("/", apiRoot)
        // GET /api/ignacio
        apiRouter.Get("/:name", greet)
    }

    http.ListenAndServe(":8080", m)
}

func apiRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "api root")
}

func greet(w http.ResponseWriter, r *http.Request) {
    name := min.GetParam(r, "name")
    fmt.Fprintf(w, "hello %s!", name)
}
```

### Middleware

Middleware in `min` are simply functions that take an `http.Handler` (the one
next in the chain) and return another one. They are resolved in the order that
they are declared.

`min` users are meant to take advantage of `context` to make better use of
middleware.

``` go
import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/arturovm/min"
)

func main() {
    m := min.New(nil)
    m.Use(logger)
    m.Use(printer)

    apiRouter := m.Group("/api")
    {
        apiRouter.Get("/", apiRoot)
        nameRouter := apiRouter.Group("/:name")
        {
            // Every request sent to routes defined on this sub-router will now
            // have a reference to a name in its context.
            // Useful for RESTful design.
            nameRouter.Use(nameExtractor)

            // GET /api/ignacio
            nameRouter.Get("/", greet)
            // GET /api/ignacio/goodbye
            nameRouter.Get("/goodbye", goodbye)
        }
    }

    http.ListenAndServe(":8080", m)
}

// -- Middleware --

// a simple logger
func logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("| %s %s", r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}

// a useless middleware that prints text
func printer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("this prints some text")
        next.ServeHTTP(w, r)
    })
}

// extracts a name from the URL and injects it into the request's context
func nameExtractor(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        name := min.GetParam(r, "name")
        ctx := context.WithValue(r.Context(), "name", name)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// -- Handlers --

func apiRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "api root")
}

// greets the user with :name
func greet(w http.ResponseWriter, r *http.Request) {
    name := r.Context().Value("name").(string)
    fmt.Fprintf(w, "hello %s!", name)
}

// says "bye" to the user with :name
func goodbye(w http.ResponseWriter, r *http.Request) {
    name := r.Context().Value("name").(string)
    fmt.Fprintf(w, "bye %s!", name)
}
```

### Base Router

If you need access to the underlying `httprouter.Router`, you can use the
`GetBaseRouter` method.

``` go
import (
    "github.com/arturovm/min"
)

func main() {
    m := min.New(nil)

    _ = m.GetBaseRouter()
}
```
