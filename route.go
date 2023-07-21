package router

type route struct {
	handler    HttpHandler
	middleware []Middleware
}

func newRoute(handler HttpHandler) *route {
	return &route{handler: handler}
}

func (r *route) Use(middleware Middleware) *route {
	r.middleware = append(r.middleware, middleware)
	return r
}
