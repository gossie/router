# http-router

The module provides an HTTP router that integrates with Go's `net/http` package. That means the application code uses `http.ResponseWriter` and `http.Request` from Go's standard library.  
It supports
- path variables
- standard middleware functions
- custom middleware functions

***Currently the module is not intended to be used in production.***

## Usage

A simple server could look like this:
```go
import (
    "net/http"

    "github.com/gossie/router"
)

func main() {
    router := router.New()

    router.Get("/books", getBooksHandler)
    router.Post("/books", createBookHandler)
    router.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
The code creates two `GET` and one `POST` route to retrieve and create books. The first parameter is the path, that may contain path variables. Path variables start with a `:`. The second parameter is the handler function that handles the request. A handler function must be of the following type: `type HttpHandler func(http.ResponseWriter, *http.Request, map[string]string)`
The first and second parameter are the `ResponseWriter` and the `Request` of Go's `http` package. The third parameter is a `map` containing the path variables. The key is the name the way it was used in the route's path. In this example the third route would contain a value for the key `bookId`.

## Standard middleware functions

### Basic auth

The module provides a standard middleware function for basic authentication. The line `testRouter.Use(middleware.BasicAuth(userChecker))` adds basic auth to the router. The `userChecker` is a function that checks if the authentication data is correct.

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func main() {
    userChecker := func(us *middleware.UserData) bool {
		// TODO: check the UserData and return true if username and password matches, false otherwise
	}

	router := router.New()

	router.Get("/books", getBooksHandler)
    router.Post("/books", createBookHandler)
    router.Get("/books/:bookId", getSingleBookHandler)

	router.Use(middleware.BasicAuth(userChecker))

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Add custom middleware

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func logRequestTime(handler router.HttpHandler) router.HttpHandler {
    return func(w http.ResponseWriter, r *http.Request, m map[string]string) {
        start := time.Now()
        defer func() {
            log.Default().Println("request took", time.Since(start).Milliseconds(), "ms")
        }()

        handler(w, r, m)
    }
}

func main() {
    router := router.New()

    router.Use(logRequestTime)

    router.Get("/books", getBooksHandler)
    router.Post("/books", createBookHandler)
    router.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
