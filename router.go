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

type HttpHandler func(http.ResponseWriter, *http.Request, map[string]string)

type route struct {
	handler HttpHandler
}

type HttpRouter struct {
	routes map[string]*pathTree
}

func New() *HttpRouter {
	return &HttpRouter{make(map[string]*pathTree)}
}

func (h *HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathVariables := make(map[string]string)

	if tree, present := h.routes[r.Method]; present {
		currentNode := tree.root
		if r.URL.Path == "/" {
			currentNode.route.handler(w, r, pathVariables)
			return
		}

		for _, el := range strings.Split(r.URL.Path, "/") {
			if el != "" && currentNode != nil {
				currentNode = currentNode.getNode(el)
				if currentNode != nil && currentNode.nodeType == "var" {
					pathVariables[currentNode.pathElement] = el
				}
			}
		}

		if currentNode != nil && currentNode.route != nil {
			currentNode.route.handler(w, r, pathVariables)
			return
		}
	}

	log.Default().Println("no", r.Method, "pattern matched", r.URL.Path, "-> returning 404")
	http.NotFound(w, r)
}

func (h *HttpRouter) addRoute(path string, method string, handler HttpHandler) {
	if _, present := h.routes[method]; !present {
		h.routes[method] = createPathTree()
	}

	if path == "/" {
		h.routes[method].root.route = &route{handler}
	}

	currentNode := h.routes[method].root
	for _, el := range strings.Split(path, "/") {
		if el != "" {
			if strings.HasPrefix(el, ":") {
				currentNode, _ = currentNode.createOrGetVarChild(el[1:])
			} else {
				currentNode, _ = currentNode.createOrGetStaticChild(el)
			}
		}
	}
	currentNode.route = &route{handler}
}

func (h *HttpRouter) Handle(path string, handler http.Handler) {
	h.addRoute(path, GET, func(w http.ResponseWriter, r *http.Request, m map[string]string) {
		w.Header().Set("Cache-Control", "public, maxage=86400, s-maxage=86400, immutable")
		w.Header().Set("Expires", time.Now().Add(86400*time.Second).Local().Format("Mon, 02 Jan 2006 15:04:05 MST"))
		handler.ServeHTTP(w, r)
	})
}

func (h *HttpRouter) Get(path string, handler HttpHandler) {
	h.addRoute(path, GET, handler)
}

func (h *HttpRouter) Put(path string, handler HttpHandler) {
	h.addRoute(path, PUT, handler)
}

func (h *HttpRouter) Patch(path string, handler HttpHandler) {
	h.addRoute(path, PATCH, handler)
}

func (h *HttpRouter) Post(path string, handler HttpHandler) {
	h.addRoute(path, POST, handler)
}
func (h *HttpRouter) Delete(path string, handler HttpHandler) {
	h.addRoute(path, DELETE, handler)
}
