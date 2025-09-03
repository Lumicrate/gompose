package http

import (
	"net/http"
)

type Context interface {
	JSON(code int, obj any)
	Bind(obj any) error
	Param(key string) string
	Query(key string) string
	QueryParams() map[string][]string
	BindJSON(obj any) error
	SetHeader(key, value string)
	Method() string
	Path() string
	SetStatus(code int)
	Status() int
	RemoteIP() string
	Header(header string) string
	Body(string)

	Abort()
	Next()

	Set(key string, value any)
	Get(key string) any

	Request() *http.Request
}
