package router

type route struct {
	method     string
	path       string
	handler    HttpHandler
	middleware []Middleware
}

func newRoute(method, path string, handler HttpHandler) *route {
	return &route{
		method:  method,
		path:    path,
		handler: handler,
	}
}

func (r *route) Use(middleware Middleware) *route {
	r.middleware = append(r.middleware, middleware)
	return r
}
