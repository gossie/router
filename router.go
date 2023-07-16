package router

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type HttpHandler = func(http.ResponseWriter, *http.Request, *Context)

type Middleware = func(HttpHandler) HttpHandler

type route struct {
	handler             HttpHandler
	middlewareFunctions []Middleware
}

func newRoute(handler HttpHandler) *route {
	return &route{handler, make([]Middleware, 0)}
}

func (r *route) Use(middleware Middleware) *route {
	r.middlewareFunctions = append(r.middlewareFunctions, middleware)
	return r
}

type HttpRouter struct {
	routes              map[string]*pathTree
	middlewareFunctions []Middleware
}

func New() *HttpRouter {
	return &HttpRouter{make(map[string]*pathTree), make([]Middleware, 0)}
}

func (h *HttpRouter) addRoute(path string, method string, handler HttpHandler) *route {
	if _, present := h.routes[method]; !present {
		h.routes[method] = newPathTree()
	}

	if path == "/" {
		h.routes[method].root.route = newRoute(handler)
	}

	currentNode := h.routes[method].root
	var err error
	for _, el := range strings.Split(path, "/") {
		if el != "" {
			if strings.HasPrefix(el, ":") {
				currentNode, err = currentNode.createOrGetVarChild(el[1:])
			} else {
				currentNode, err = currentNode.createOrGetStaticChild(el)
			}

			if err != nil {
				panic(err.Error())
			}
		}
	}

	currentNode.route = newRoute(handler)
	return currentNode.route
}

func (h *HttpRouter) Handle(path string, handler http.Handler) {
	h.addRoute(path, GET, func(w http.ResponseWriter, r *http.Request, _ *Context) {
		w.Header().Set("Cache-Control", "public, maxage=86400, s-maxage=86400, immutable")
		w.Header().Set("Expires", time.Now().Add(86400*time.Second).Local().Format("Mon, 02 Jan 2006 15:04:05 MST"))
		handler.ServeHTTP(w, r)
	})
}

func (h *HttpRouter) Get(path string, handler HttpHandler) *route {
	return h.addRoute(path, GET, handler)
}

func (h *HttpRouter) Put(path string, handler HttpHandler) *route {
	return h.addRoute(path, PUT, handler)
}

func (h *HttpRouter) Patch(path string, handler HttpHandler) *route {
	return h.addRoute(path, PATCH, handler)
}

func (h *HttpRouter) Post(path string, handler HttpHandler) *route {
	return h.addRoute(path, POST, handler)
}

func (h *HttpRouter) Delete(path string, handler HttpHandler) *route {
	return h.addRoute(path, DELETE, handler)
}

func (h *HttpRouter) Use(middleware Middleware) {
	h.middlewareFunctions = append(h.middlewareFunctions, middleware)
}

func (h *HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathVariables := make(map[string]string)

	if tree, present := h.routes[r.Method]; present {
		currentNode := tree.root
		if r.URL.Path == "/" {
			currentNode.route.handler(w, r, newContext(pathVariables))
			return
		}

		for _, el := range strings.Split(r.URL.Path, "/") {
			if el != "" && currentNode != nil {
				currentNode = currentNode.childNode(el)
				if currentNode != nil && currentNode.nodeType == "var" {
					pathVariables[currentNode.pathElement] = el
				}
			}
		}

		if currentNode != nil && currentNode.route != nil {
			handlerToExceute := currentNode.route.handler
			allMiddlewareFunctions := make([]Middleware, 0, len(h.middlewareFunctions)+len(currentNode.route.middlewareFunctions))
			allMiddlewareFunctions = append(allMiddlewareFunctions, h.middlewareFunctions...)
			allMiddlewareFunctions = append(allMiddlewareFunctions, currentNode.route.middlewareFunctions...)
			for i := len(allMiddlewareFunctions) - 1; i >= 0; i-- {
				handlerToExceute = allMiddlewareFunctions[i](handlerToExceute)
			}
			handlerToExceute(w, r, newContext(pathVariables))
			return
		}
	}

	log.Default().Println("no", r.Method, "pattern matched", r.URL.Path, "-> returning 404")
	http.NotFound(w, r)
}
