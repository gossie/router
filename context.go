package router

type Context struct {
	pathParameters map[string]string
	username       string
}

func newContext(pathParameters map[string]string) *Context {
	return &Context{
		pathParameters: pathParameters,
	}
}

func (ctx *Context) PathParameter(name string) string {
	return ctx.pathParameters[name]
}

func (ctx *Context) Username() string {
	return ctx.username
}
