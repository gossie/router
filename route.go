package router

type route struct {
	handler    HttpHandler
	middleware []Middleware
}

func newRoute(handler HttpHandler) *route {
	return &route{handler, make([]Middleware, 0)}
}

func (r *route) Use(middleware Middleware) *route {
	r.middleware = append(r.middleware, middleware)
	return r
}
