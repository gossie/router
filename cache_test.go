package router_test

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gossie/router"
	"github.com/stretchr/testify/assert"
)

func TestCache_noCache(t *testing.T) {
	testRouter := router.New()
	testRouter.Get("/route", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	testRouter.FinishSetup()

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/route"},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 200, w.statusCode)
	assert.Equal(t, "", w.Header().Get("Cache-Control"))
}

func TestCache_cache(t *testing.T) {
	testRouter := router.New()
	testRouter.Get("/route", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	testRouter.Use(router.Cache(1 * time.Hour))

	testRouter.FinishSetup()

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/route"},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 200, w.statusCode)
	assert.Equal(t, "public, maxage=3600, s-maxage=3600, immutable", w.Header().Get("Cache-Control"))
}
