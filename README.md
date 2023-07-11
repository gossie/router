# http-router

# Usage

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

# Add custom middleware

```go
import (
    "net/http"

    "github.com/gossie/router"
)

func trackRequestRuntime(handler router.HttpHandler) router.HttpHandler {
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

    router.Use(trackRequestRuntime)

    router.Get("/books", getBooksHandler)
    router.Post("/books", createBookHandler)
    router.Get("/books/:bookId", getSingleBookHandler)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```
