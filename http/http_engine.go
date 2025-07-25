package http

type HandlerFunc func(ctx Context)
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type HTTPEngine interface {
	Init(port int) error
	RegisterRoute(method string, path string, handler HandlerFunc)
	Use(middleware MiddlewareFunc)
	Start() error
}
