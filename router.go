package router

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	GET                  = "GET"
	POST                 = "POST"
	PUT                  = "PUT"
	DELETE               = "DELETE"
	PATCH                = "PATCH"
	SEPARATOR            = "/"
	PATH_VARIABLE_PREFIX = ":"
)

type HttpHandler = func(w http.ResponseWriter, r *http.Request)

type Middleware = func(HttpHandler) HttpHandler

type HttpRouter struct {
	mutex      sync.RWMutex
	routes     []*route
	middleware []Middleware
}

func New() *HttpRouter {
	http.DefaultServeMux = &http.ServeMux{}
	return &HttpRouter{routes: make([]*route, 0)}
}

func (hr *HttpRouter) addRoute(path string, method string, handler HttpHandler) *route {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	newRoute := newRoute(method, path, handler)
	hr.routes = append(hr.routes, newRoute)

	return newRoute
}

func (hr *HttpRouter) Handle(path string, handler http.Handler) {
	hr.addRoute(path, GET, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, maxage=86400, s-maxage=86400, immutable")
		w.Header().Set("Expires", time.Now().Add(86400*time.Second).Local().Format("Mon, 02 Jan 2006 15:04:05 MST"))
		handler.ServeHTTP(w, r)
	})
}

func (hr *HttpRouter) Get(path string, handler HttpHandler) *route {
	return hr.addRoute(path, GET, handler)
}

func (hr *HttpRouter) Put(path string, handler HttpHandler) *route {
	return hr.addRoute(path, PUT, handler)
}

func (hr *HttpRouter) Patch(path string, handler HttpHandler) *route {
	return hr.addRoute(path, PATCH, handler)
}

func (hr *HttpRouter) Post(path string, handler HttpHandler) *route {
	return hr.addRoute(path, POST, handler)
}

func (hr *HttpRouter) Delete(path string, handler HttpHandler) *route {
	return hr.addRoute(path, DELETE, handler)
}

func (hr *HttpRouter) Use(middleware Middleware) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.middleware = append(hr.middleware, middleware)
}

func (hr *HttpRouter) FinishSetup() {
	for _, r := range hr.routes {
		handler := r.handler

		for i := len(r.middleware) - 1; i >= 0; i-- {
			handler = r.middleware[i](handler)
		}

		for i := len(hr.middleware) - 1; i >= 0; i-- {
			handler = hr.middleware[i](handler)
		}

		http.Handle(fmt.Sprintf("%v %v", r.method, r.path), http.HandlerFunc(handler))
	}
}

func (hr *HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
