package auth

import "github.com/Lumicrate/gompose/http"

type AuthProvider interface {
	Init() error
	RegisterRoutes(engine http.HTTPEngine)
	Middleware() http.MiddlewareFunc
}
