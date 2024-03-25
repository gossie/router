# http-router

The module provides an HTTP router that integrates with Go's `net/http` package. That means the application code uses `http.ResponseWriter` and `http.Request` from Go's standard library.  
It supports
- standard middleware functions
- custom middleware functions

## But why ...

There are already a lot of mux implementations. But after a brief search I only found implementations that did not match my requirements oder that overfullfill them.
The things I wanted were
- path variables
- routes for certain methods
- built-in middleware support
- I wanted to stay as close to Go's `net/http` package as possible

And last but certainly not least: **It's a lot of fun to implement such a thing yourself!**

Since Go 1.22 the standard http package supports path variables and the possibility to specify a handler for a method-path combination. This HTTP router was migrated to make use of the new standard features. So now it just provides a different (but in my opinion clearer) API and the possibility to define middleware functions.

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
    httpRouter.Get("/books/{bookId}", getSingleBookHandler)

    httpRouter.FinishSetup()

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```
The code creates two `GET` and one `POST` route to retrieve and create books. The first parameter is the path, that may contain path variables. Path variables use Go 1.22's standard syntax. The second parameter is the handler function that handles the request. A handler function must be of the following type: `type HttpHandler func(http.ResponseWriter, *http.Request)`
The first and second parameter are the `ResponseWriter` and the `Request` of Go's `http` package.

## Middleware

Middleware functions can be used to reuse behaviour that should be executed on every or a couple of request. Typical examples are authentication, request logging, etc.  
Middleware functions are added with the `Use` method. `Use` can be called directly on the router. Middleware that is added that way will be exexuted for every request. Alternatively `Use` can be called on a created route. That way the middleware will only be executed, if the specific route is called.

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func middleware1(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        // ...
    }
}

func middleware2(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        // ...
    }
}

func main() {
    httpRouter := router.New()

    httpRouter.Use(middleware1)

    httpRouter.Get("/test1", publicHandler)
    httpRouter.Post("/test2", protectedHanlder).Use(middleware2)

    httpRouter.FinishSetup()

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```

### Standard middleware functions

#### Basic auth

The module provides a standard middleware function for basic authentication. The line `testRouter.Use(router.BasicAuth(userChecker))` adds basic auth to the router. The `userChecker` is a function that checks if the authentication data is correct. If the user was authenticated, the username will be added to the `context` of the request under the key `router.UsernameKey`.

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
    httpRouter.Get("/books/{bookId}", getSingleBookHandler)

    httpRouter.FinishSetup()

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
    httpRouter.Get("/books/{bookId}", getSingleBookHandler)

    httpRouter.FinishSetup()

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
    return func(w http.ResponseWriter, r *http.Request) {
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
    httpRouter.Get("/books/{bookId}", getSingleBookHandler)

    httpRouter.FinishSetup()

    log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
```
