# http-router

The module provides an HTTP router that integrates with Go's `net/http` package. That means the application code uses `http.ResponseWriter` and `http.Request` from Go's standard library.  
It supports
- path variables
- standard middleware functions
- custom middleware functions

## But why ...

There are already a lot of mux implementations. But after a brief search I only found implementations that did not match my requirements oder that overfullfill them.
The things I wanted were
- path variables
- built-in middleware support
- I wanted to stay as close to Go's `net/http` package as possible

And last but certainly not least: **It's a lot of fun to implement such a thing yourself!**

## Usage

A simple server could look like this:
```go
import (
    "net/http"

    "github.com/gossie/router"
)

func main() {
    httpRouter := router.New()

    httpRouter.Get("/books", getBooksHandler)
    httpRouter.Post("/books", createBookHandler)
    httpRouter.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```
The code creates two `GET` and one `POST` route to retrieve and create books. The first parameter is the path, that may contain path variables. Path variables start with a `:`. The second parameter is the handler function that handles the request. A handler function must be of the following type: `type HttpHandler func(http.ResponseWriter, *http.Request, *router.Context)`
The first and second parameter are the `ResponseWriter` and the `Request` of Go's `http` package. The third parameter is a `map` containing the path variables. The key is the name the way it was used in the route's path. In this example the third route would contain a value for the key `bookId`.

## Middleware

Middleware functions can be used to reuse behaviour that should be executed on every or a couple of request. Typical examples are authentication, request logging, etc.  
Middleware functions are added with the `Use` method. `Use` can be called directly on the router. Middleware that is added that way will be exexuted for every request. Alternatively `Use` can be called on a created route. That way the middleware will only be executed, if the specific route is called.

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func middleware1(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request, ctx *router.Context) {
        // ...
    }
}

func middleware2(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request, ctx *router.Context) {
        // ...
    }
}

func main() {
    httpRouter := router.New()

    httpRouter.Use(middleware1)

    httpRouter.Get("/test1", publicHandler)
    httpRouter.Post("/test2", protectedHanlder).Use(middleware2)

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```

There is a third way to add a middleware function. It is possible to define a middleware function for a certain path and HTTP method.

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func middleware(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request, ctx *router.Context) {
        // ...
    }
}

func main() {
    httpRouter := router.New()

    testRouter.UseRecursively(router.GET, "/tests", middleware)

    httpRouter.Get("/tests", testsHandler)
    httpRouter.Get("/tests/:testId", singleTestHandler)
    httpRouter.Get("/tests/:testId/assertions", assertionsHandler)
    httpRouter.Get("/other", otherHandler)

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```

The code makes sure that the middleware function is executed for `GET` request targeting `/tests`, `/tests/:testId` and `/tests/:testId/assertions`. It won't be executed when `/other` is called.

### Standard middleware functions

#### Basic auth

The module provides a standard middleware function for basic authentication. The line `testRouter.Use(router.BasicAuth(userChecker))` adds basic auth to the router. The `userChecker` is a function that checks if the authentication data is correct.

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func main() {
    userChecker := func(us *router.UserData) bool {
        // TODO: check the UserData and return true if username and password matches, false otherwise
    }

    httpRouter := router.New()

    httpRouter.Use(router.BasicAuth(userChecker))

    httpRouter.Get("/books", getBooksHandler)
    httpRouter.Post("/books", createBookHandler)
    httpRouter.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```

#### Cache headers

The module provides a standard middleware function to activate browser caching. The line `testRouter.Use(router.Cache(1 * time.Hour))` makes sure that the necessary headers are set, so that the reponse is cache one hour by the browser.

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func main() {
    httpRouter := router.New()

    httpRouter.Use(router.Cache(1 * time.Hour))

    httpRouter.Get("/books", getBooksHandler)
    httpRouter.Post("/books", createBookHandler)
    httpRouter.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```

### Add custom middleware

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func logRequestTime(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request, ctx *router.Context) {
        start := time.Now()
        defer func() {
            log.Default().Println("request took", time.Since(start).Milliseconds(), "ms")
        }()

        handler(w, r, m)
    }
}

func main() {
    httpRouter := router.New()

    httpRouter.Use(logRequestTime)

    httpRouter.Get("/books", getBooksHandler)
    httpRouter.Post("/books", createBookHandler)
    httpRouter.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```
