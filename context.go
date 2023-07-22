package router

type pathParam struct {
	name, value string
}

type Context struct {
	pathParameters []pathParam
	username       string
}

func newContext(pathParameters []pathParam) *Context {
	return &Context{
		pathParameters: pathParameters,
	}
}

func (ctx *Context) PathParameter(name string) string {
	for _, p := range ctx.pathParameters {
		if name == p.name {
			return p.value
		}
	}
	return ""
}

func (ctx *Context) Username() string {
	return ctx.username
}
