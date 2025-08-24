package http

type HandlerFunc func(ctx Context)
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Route struct {
	Method    string
	Path      string
	Entity    any
	Protected bool
}

type HTTPEngine interface {
	Init(port int) error
	RegisterRoute(method string, path string, handler HandlerFunc, entity any, isProtected bool)
	Use(middleware MiddlewareFunc)
	Start() error
	Routes() []Route
}
