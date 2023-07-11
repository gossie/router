# http-router

# Usage

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