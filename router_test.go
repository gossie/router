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

	router := router.New()

	router.Post("/tests", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		testString = "post-was-called"
	})

	router.Get("/tests/:testString", func(w http.ResponseWriter, r *http.Request, pv map[string]string) {
		if testString != pv["testString"] {
			t.Fatalf("%s != %s", testString, pv["testString"])
		}
	})

	router.Get("/tests/:testString/:detailId", func(w http.ResponseWriter, r *http.Request, pv map[string]string) {
		if testString != pv["testString"] {
			t.Fatalf("%s != %s", testString, pv["testString"])
		}
	})

	router.Put("/tests/:testString", func(w http.ResponseWriter, r *http.Request, pv map[string]string) {
		if testString == pv["testString"] {
			testString = "put-was-called"
		}
	})

	router.Patch("/tests/:testString", func(w http.ResponseWriter, r *http.Request, pv map[string]string) {
		if testString == pv["testString"] {
			testString = "patch-was-called"
		}
	})

	router.Delete("/tests/:testString", func(w http.ResponseWriter, r *http.Request, pv map[string]string) {
		if testString == pv["testString"] {
			testString = ""
		}
	})

	w1 := &TestResponseWriter{}
	r1 := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/tests"},
	}
	router.ServeHTTP(w1, r1)

	assert.Equal(t, "post-was-called", testString)

	w2 := &TestResponseWriter{}
	r2 := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/tests/post-was-called"},
	}
	router.ServeHTTP(w2, r2)

	w3 := &TestResponseWriter{}
	r3 := &http.Request{
		Method: "PUT",
		URL:    &url.URL{Path: "/tests/not-found"},
	}
	router.ServeHTTP(w3, r3)

	assert.Equal(t, "post-was-called", testString)

	w4 := &TestResponseWriter{}
	r4 := &http.Request{
		Method: "PUT",
		URL:    &url.URL{Path: "/tests/post-was-called"},
	}
	router.ServeHTTP(w4, r4)

	assert.Equal(t, "put-was-called", testString)

	w5 := &TestResponseWriter{}
	r5 := &http.Request{
		Method: "PATCH",
		URL:    &url.URL{Path: "/tests/put-was-called"},
	}
	router.ServeHTTP(w5, r5)

	assert.Equal(t, "patch-was-called", testString)

	w6 := &TestResponseWriter{}
	r6 := &http.Request{
		Method: "DELETE",
		URL:    &url.URL{Path: "/tests/unknown"},
	}
	router.ServeHTTP(w6, r6)

	assert.Equal(t, "patch-was-called", testString)

	w7 := &TestResponseWriter{}
	r7 := &http.Request{
		Method: "DELETE",
		URL:    &url.URL{Path: "/tests/patch-was-called"},
	}
	router.ServeHTTP(w7, r7)

	assert.Empty(t, testString)
}

func TestHasRootRoute(t *testing.T) {
	emptyHandler := func(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
		w.WriteHeader(200)
	}

	router := router.New()
	router.Get("/", emptyHandler)

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/"},
	}

	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.statusCode)
}

func TestReturnsStatus404(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {}

	router := router.New()
	router.Get("/tests/:id", emptyHandler)

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/tests/test-id/details"},
	}

	router.ServeHTTP(w, r)

	assert.Equal(t, 404, w.statusCode)
}

func TestCreatesVariableAndStaticElementAtTheSamePosition(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {}

	router := router.New()

	assert.Panics(t, func() {
		router.Get("/tests/:id", emptyHandler)
		router.Get("/tests/green", emptyHandler)
	})
}

func TestCreatesStaticElementAndVariableAtTheSamePosition(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {}

	router := router.New()

	assert.Panics(t, func() {
		router.Get("/tests/green", emptyHandler)
		router.Get("/tests/:id", emptyHandler)
	})
}

func TestCreatesTwoVariablesAtTheSamePosition(t *testing.T) {
	emptyHandler := func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {}

	router := router.New()

	assert.Panics(t, func() {
		router.Get("/tests/:testId", emptyHandler)
		router.Get("/tests/:id", emptyHandler)
	})
}
