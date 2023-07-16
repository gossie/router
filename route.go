package router

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
