package router

import (
	"log"
	"math"
	"net/http"
	"strings"
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

type HttpHandler = func(http.ResponseWriter, *http.Request, Context)

type Middleware = func(HttpHandler) HttpHandler

type HttpRouter struct {
	mutex             sync.RWMutex
	routes            map[string]*pathTree
	middleware        []Middleware
	pathVariableCount uint
}

func New() *HttpRouter {
	return &HttpRouter{routes: make(map[string]*pathTree)}
}

func (hr *HttpRouter) addRoute(path string, method string, handler HttpHandler) *route {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	hr.pathVariableCount = uint(math.Max(float64(hr.pathVariableCount), float64(strings.Count(path, PATH_VARIABLE_PREFIX))))

	rootHandler := func() {
		hr.routes[method].root.route = newRoute(handler)
	}

	currentNode := hr.getCreateOrGetNode(path, method, rootHandler)

	currentNode.route = newRoute(handler)
	return currentNode.route
}

func (hr *HttpRouter) Handle(path string, handler http.Handler) {
	hr.addRoute(path, GET, func(w http.ResponseWriter, r *http.Request, _ Context) {
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

func (hr *HttpRouter) UseRecursively(method, path string, middleware Middleware) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	rootHandler := func() {
		panic("use the Use() method")
	}

	currentNode := hr.getCreateOrGetNode(path, method, rootHandler)

	currentNode.middleware = append(currentNode.middleware, middleware)
}

func (hr *HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hr.mutex.RLock()
	defer hr.mutex.RUnlock()

	var pathVariables []pathParam

	if tree, present := hr.routes[r.Method]; present {
		currentNode := tree.root
		if r.URL.Path == SEPARATOR {
			currentNode.route.handler(w, r, newContext(nil))
			return
		}

		middlewareToExecute := appendMiddlewareIfNeeded(nil, hr.middleware)

		currentPath := r.URL.Path[1:]
		index := strings.Index(currentPath, SEPARATOR)
		for index > 0 || currentPath != "" {
			var el string
			if index < 0 {
				el = currentPath
				currentPath = ""
			} else {
				el = currentPath[0:index]
				currentPath = currentPath[index+1:]
			}

			if currentNode != nil {
				middlewareToExecute = appendMiddlewareIfNeeded(middlewareToExecute, currentNode.middleware)
				currentNode = currentNode.childNode(el)
				if currentNode != nil && currentNode.nodeType == NodeTypeVar {
					if pathVariables == nil {
						pathVariables = make([]pathParam, 0, hr.pathVariableCount)
					}
					pathVariables = append(pathVariables, pathParam{name: currentNode.pathElement, value: el})
				}
			}

			index = strings.Index(currentPath, SEPARATOR)
		}

		if currentNode == nil || currentNode.route == nil {
			log.Default().Println("no", r.Method, "pattern matched", r.URL.Path, "-> returning 404")
			http.NotFound(w, r)
			return
		}

		handlerToExceute := currentNode.route.handler
		middlewareToExecute = appendMiddlewareIfNeeded(middlewareToExecute, currentNode.route.middleware)
		for i := len(middlewareToExecute) - 1; i >= 0; i-- {
			handlerToExceute = middlewareToExecute[i](handlerToExceute)
		}

		handlerToExceute(w, r, newContext(pathVariables))
	}
}

func (hr *HttpRouter) getCreateOrGetNode(path string, method string, rootHandler func()) *node {
	if _, present := hr.routes[method]; !present {
		hr.routes[method] = newPathTree()
	}

	if path == SEPARATOR {
		rootHandler()
	}

	currentNode := hr.routes[method].root
	var err error
	for _, el := range strings.Split(path, SEPARATOR) {
		if el != "" {
			if strings.HasPrefix(el, PATH_VARIABLE_PREFIX) {
				currentNode, err = currentNode.createOrGetVarChild(el[1:])
			} else {
				currentNode, err = currentNode.createOrGetStaticChild(el)
			}

			if err != nil {
				panic(err.Error())
			}
		}
	}
	return currentNode
}

func appendMiddlewareIfNeeded(current []Middleware, source []Middleware) []Middleware {
	if len(source) > 0 {
		return append(current, source...)
	}
	return current
}
