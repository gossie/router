package router_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gossie/router"
	"github.com/stretchr/testify/assert"
)

type TestResponseWriter struct {
	statusCode int
	headers    map[string][]string
}

func (w *TestResponseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(map[string][]string)
	}

	return w.headers
}

func (w *TestResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *TestResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func TestRouting(t *testing.T) {
	var testString string

	testRouter := router.New()

	testRouter.Post("/tests", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		testString = "post-was-called"
	})

	testRouter.Get("/tests/:testString", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		if testString != ctx.PathParameter("testString") {
			t.Fatalf("%s != %s", testString, ctx.PathParameter("testString"))
		}
	})

	testRouter.Get("/tests/:testString/:detailId", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		if testString != ctx.PathParameter("testString") {
			t.Fatalf("%s != %s", testString, ctx.PathParameter("testString"))
		}
	})

	testRouter.Put("/tests/:testString", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		if testString == ctx.PathParameter("testString") {
			testString = "put-was-called"
		}
	})

	testRouter.Patch("/tests/:testString", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		if testString == ctx.PathParameter("testString") {
			testString = "patch-was-called"
		}
	})

	testRouter.Delete("/tests/:testString", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		if testString == ctx.PathParameter("testString") {
			testString = ""
		}
	})

	w1 := &TestResponseWriter{}
	r1 := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/tests"},
	}
	testRouter.ServeHTTP(w1, r1)

	assert.Equal(t, "post-was-called", testString)

	w2 := &TestResponseWriter{}
	r2 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/tests/post-was-called"},
	}
	testRouter.ServeHTTP(w2, r2)

	w3 := &TestResponseWriter{}
	r3 := &http.Request{
		Method: "PUT",
		URL:    &url.URL{Path: "/tests/not-found"},
	}
	testRouter.ServeHTTP(w3, r3)

	assert.Equal(t, "post-was-called", testString)

	w4 := &TestResponseWriter{}
	r4 := &http.Request{
		Method: "PUT",
		URL:    &url.URL{Path: "/tests/post-was-called"},
	}
	testRouter.ServeHTTP(w4, r4)

	assert.Equal(t, "put-was-called", testString)

	w5 := &TestResponseWriter{}
	r5 := &http.Request{
		Method: "PATCH",
		URL:    &url.URL{Path: "/tests/put-was-called"},
	}
	testRouter.ServeHTTP(w5, r5)

	assert.Equal(t, "patch-was-called", testString)

	w6 := &TestResponseWriter{}
	r6 := &http.Request{
		Method: "DELETE",
		URL:    &url.URL{Path: "/tests/unknown"},
	}
	testRouter.ServeHTTP(w6, r6)

	assert.Equal(t, "patch-was-called", testString)

	w7 := &TestResponseWriter{}
	r7 := &http.Request{
		Method: "DELETE",
		URL:    &url.URL{Path: "/tests/patch-was-called"},
	}
	testRouter.ServeHTTP(w7, r7)

	assert.Empty(t, testString)
}

func TestHasRootRoute(t *testing.T) {
	emptyHandler := func(w http.ResponseWriter, _ *http.Request, _ router.Context) {
		w.WriteHeader(200)
	}

	testRouter := router.New()
	testRouter.Get("/", emptyHandler)

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/"},
	}

	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 200, w.statusCode)
}

func TestReturnsStatus404(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ router.Context) {}

	testRouter := router.New()
	testRouter.Get("/tests/:id", emptyHandler)

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/tests/test-id/details"},
	}

	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 404, w.statusCode)
}

func TestCreatesVariableAndStaticElementAtTheSamePosition(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ router.Context) {}

	testRouter := router.New()

	assert.Panics(t, func() {
		testRouter.Get("/tests/:id", emptyHandler)
		testRouter.Get("/tests/green", emptyHandler)
	})
}

func TestCreatesStaticElementAndVariableAtTheSamePosition(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ router.Context) {}

	testRouter := router.New()

	assert.Panics(t, func() {
		testRouter.Get("/tests/green", emptyHandler)
		testRouter.Get("/tests/:id", emptyHandler)
	})
}

func TestCreatesTwoVariablesAtTheSamePosition(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ router.Context) {}

	testRouter := router.New()

	assert.Panics(t, func() {
		testRouter.Get("/tests/:testId", emptyHandler)
		testRouter.Get("/tests/:id", emptyHandler)
	})
}

func TestMiddleware(t *testing.T) {
	executed := make([]string, 0)

	middleware1 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware1")
			in(w, r, ctx)
		}
	}

	middleware2 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware2")
			in(w, r, ctx)
		}
	}

	testRouter := router.New()

	testRouter.Use(middleware1)
	testRouter.Use(middleware2)

	testRouter.Get("/test", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "get")
	})

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test"},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 3, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "middleware2", executed[1])
	assert.Equal(t, "get", executed[2])
}

func TestMiddlewareForSingleRoute(t *testing.T) {
	executed := make([]string, 0)

	middleware1 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware1")
			in(w, r, ctx)
		}
	}

	middleware2 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware2")
			in(w, r, ctx)
		}
	}

	testRouter := router.New()

	testRouter.Use(middleware1)

	testRouter.Get("/test1", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "test1")
	})

	testRouter.Get("/test2", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "test2")
	}).Use(middleware2)

	w := &TestResponseWriter{}
	r1 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test1"},
	}
	r2 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/test2"},
	}

	testRouter.ServeHTTP(w, r1)

	assert.Equal(t, 2, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "test1", executed[1])

	testRouter.ServeHTTP(w, r2)

	assert.Equal(t, 5, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "test1", executed[1])
	assert.Equal(t, "middleware1", executed[2])
	assert.Equal(t, "middleware2", executed[3])
	assert.Equal(t, "test2", executed[4])
}

func TestMiddlewareForGroupOfRoutes(t *testing.T) {
	executed := make([]string, 0)

	middleware1 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware1")
			in(w, r, ctx)
		}
	}

	middleware2 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware2")
			in(w, r, ctx)
		}
	}

	middleware3 := func(in router.HttpHandler) router.HttpHandler {
		return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
			executed = append(executed, "middleware3")
			in(w, r, ctx)
		}
	}

	testRouter := router.New()

	testRouter.Use(middleware1)
	testRouter.UseRecursively(router.GET, "/tests", middleware2)

	testRouter.Get("/tests/test1", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "GET test1")
	})

	testRouter.Get("/tests/test2", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "GET test2")
	}).Use(middleware3)

	testRouter.Post("/tests/test2", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "POST test2")
	})

	testRouter.Get("/other", func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		executed = append(executed, "GET other")
	})

	w := &TestResponseWriter{}
	r1 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/tests/test1"},
	}
	testRouter.ServeHTTP(w, r1)

	assert.Equal(t, 3, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "middleware2", executed[1])
	assert.Equal(t, "GET test1", executed[2])

	r2 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/tests/test2"},
	}
	testRouter.ServeHTTP(w, r2)

	assert.Equal(t, 7, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "middleware2", executed[1])
	assert.Equal(t, "GET test1", executed[2])
	assert.Equal(t, "middleware1", executed[3])
	assert.Equal(t, "middleware2", executed[4])
	assert.Equal(t, "middleware3", executed[5])
	assert.Equal(t, "GET test2", executed[6])

	r3 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/other"},
	}
	testRouter.ServeHTTP(w, r3)

	assert.Equal(t, 9, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "middleware2", executed[1])
	assert.Equal(t, "GET test1", executed[2])
	assert.Equal(t, "middleware1", executed[3])
	assert.Equal(t, "middleware2", executed[4])
	assert.Equal(t, "middleware3", executed[5])
	assert.Equal(t, "GET test2", executed[6])
	assert.Equal(t, "middleware1", executed[7])
	assert.Equal(t, "GET other", executed[8])

	r4 := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/tests/test2"},
	}
	testRouter.ServeHTTP(w, r4)

	assert.Equal(t, 11, len(executed))
	assert.Equal(t, "middleware1", executed[0])
	assert.Equal(t, "middleware2", executed[1])
	assert.Equal(t, "GET test1", executed[2])
	assert.Equal(t, "middleware1", executed[3])
	assert.Equal(t, "middleware2", executed[4])
	assert.Equal(t, "middleware3", executed[5])
	assert.Equal(t, "GET test2", executed[6])
	assert.Equal(t, "middleware1", executed[7])
	assert.Equal(t, "GET other", executed[8])
	assert.Equal(t, "middleware1", executed[9])
	assert.Equal(t, "POST test2", executed[10])
}

// func TestRouteCaseInsensitivity(t *testing.T) {
// 	executed := false

// 	testRouter := router.New()

// 	testRouter.Get("/TEST1/:id/test2", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
// 		assert.Equal(t, "aBc", ctx.PathParameter("id"))
// 		executed = true
// 	})

// 	w := &TestResponseWriter{}
// 	r := &http.Request{
// 		Method: "GET",
// 		URL:    &url.URL{Path: "/tEsT1/aBc/TeSt2"},
// 	}

// 	testRouter.ServeHTTP(w, r)

// 	assert.True(t, executed)
// }
