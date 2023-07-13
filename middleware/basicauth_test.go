package middleware_test

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/gossie/router"
	"github.com/gossie/router/middleware"
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

func TestBasicAuth_noAuthData(t *testing.T) {
	users := []*middleware.UserData{
		middleware.NewUserData("user1", "password1"),
		middleware.NewUserData("user2", "password2"),
		middleware.NewUserData("user3", "password3"),
	}

	testRouter := router.New()
	testRouter.Get("/protected", func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		w.WriteHeader(200)
	})
	testRouter.Use(middleware.BasicAuth(users))

	w := &TestResponseWriter{}
	r := &http.Request{
		Method: "GET",
		URL:    &url.URL{Path: "/protected"},
	}
	testRouter.ServeHTTP(w, r)

	assert.Equal(t, 401, w.statusCode)
}

func TestBasicAuth_wrongAuthData(t *testing.T) {
	users := []*middleware.UserData{
		middleware.NewUserData("user1", "password1"),
		middleware.NewUserData("user2", "password2"),
		middleware.NewUserData("user3", "password3"),
	}

	testRouter := router.New()
	testRouter.Get("/protected", func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		w.WriteHeader(200)
	})
	testRouter.Use(middleware.BasicAuth(users))

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
	users := []*middleware.UserData{
		middleware.NewUserData("user1", "password1"),
		middleware.NewUserData("user2", "password2"),
		middleware.NewUserData("user3", "password3"),
	}

	testRouter := router.New()
	testRouter.Get("/protected", func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		w.WriteHeader(200)
	})
	testRouter.Use(middleware.BasicAuth(users))

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
