package router_test

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/gossie/router"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuth_noAuthData(t *testing.T) {
	userChecker := func(us *router.UserData) bool {
		return us.Username() == "user2" && us.Password() == "password2"
	}

	testRouter := router.New()
	testRouter.Get("/protected", func(_ http.ResponseWriter, _ *http.Request, _ router.Context) {
		assert.Fail(t, "handler must not be called")
	})
	testRouter.Use(router.BasicAuth(userChecker))

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/protected"},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 401, w.statusCode)
}

func TestBasicAuth_wrongAuthData(t *testing.T) {
	userChecker := func(us *router.UserData) bool {
		return us.Username() == "user2" && us.Password() == "password2"
	}

	testRouter := router.New()
	testRouter.Get("/protected", func(_ http.ResponseWriter, _ *http.Request, _ router.Context) {
		assert.Fail(t, "handler must not be called")
	})
	testRouter.Use(router.BasicAuth(userChecker))

	userStr := base64.StdEncoding.EncodeToString([]byte("user2:wrong"))

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/protected"},
		Header: map[string][]string{"Authorization": {"Basic " + userStr}},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 401, w.statusCode)
}

func TestBasicAuth_correctAuthData(t *testing.T) {
	userChecker := func(us *router.UserData) bool {
		return us.Username() == "user2" && us.Password() == "password2"
	}

	testRouter := router.New()
	testRouter.Get("/protected", func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		assert.Equal(t, "user2", ctx.Username())
		w.WriteHeader(200)
	})
	testRouter.Use(router.BasicAuth(userChecker))

	userStr := base64.StdEncoding.EncodeToString([]byte("user2:password2"))

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/protected"},
		Header: map[string][]string{"Authorization": {"Basic " + userStr}},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 200, w.statusCode)
}
