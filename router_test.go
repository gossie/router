package router_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gossie/router"
)

type TestResponseWriter struct {
}

func (w *TestResponseWriter) Header() http.Header {
	return make(map[string][]string)
}

func (w *TestResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *TestResponseWriter) WriteHeader(statusCode int) {

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

	router.Put("/tests/:testString", func(w http.ResponseWriter, r *http.Request, pv map[string]string) {
		if testString == pv["testString"] {
			testString = "put-was-called"
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

	if testString != "post-was-called" {
		t.Fatalf("%s != post-was-called", testString)
	}

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

	if testString != "post-was-called" {
		t.Fatalf("%s != post-was-called", testString)
	}

	w4 := &TestResponseWriter{}
	r4 := &http.Request{
		Method: "PUT",
		URL:    &url.URL{Path: "/tests/post-was-called"},
	}
	router.ServeHTTP(w4, r4)

	if testString != "put-was-called" {
		t.Fatalf("%s != put-was-called", testString)
	}

	w5 := &TestResponseWriter{}
	r5 := &http.Request{
		Method: "DELETE",
		URL:    &url.URL{Path: "/tests/unknown"},
	}
	router.ServeHTTP(w5, r5)

	if testString != "put-was-called" {
		t.Fatalf("%s != put-was-called", testString)
	}

	w6 := &TestResponseWriter{}
	r6 := &http.Request{
		Method: "DELETE",
		URL:    &url.URL{Path: "/tests/put-was-called"},
	}
	router.ServeHTTP(w6, r6)

	if testString != "" {
		t.Fatalf("%s is not empty", testString)
	}
}
