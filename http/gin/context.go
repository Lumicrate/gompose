package ginadapter

import "github.com/gin-gonic/gin"

type GinContext struct {
	ctx    *gin.Context
	values map[string]any
}

func (g *GinContext) JSON(code int, obj any) {
	g.ctx.JSON(code, obj)
}

func (g *GinContext) Bind(obj any) error {
	return g.ctx.ShouldBindJSON(obj)
}

func (g *GinContext) Param(key string) string {
	return g.ctx.Param(key)
}

func (g *GinContext) Query(key string) string {
	return g.ctx.Query(key)
}
func (g *GinContext) SetHeader(key, value string) {
	g.ctx.Writer.Header().Set(key, value)
}

func (g *GinContext) Method() string {
	return g.ctx.Request.Method
}

func (g *GinContext) Path() string {
	return g.ctx.Request.URL.Path
}

func (g *GinContext) Status() int {
	return g.ctx.Writer.Status()
}

func (g *GinContext) SetStatus(code int) {
	g.ctx.Status(code)
}

func (g *GinContext) RemoteIP() string {
	return g.ctx.ClientIP()
}

func (g *GinContext) Abort() {
	g.ctx.Abort()
}

func (g *GinContext) Next() {
	g.ctx.Next()
}

func (g *GinContext) Header(header string) string {
	return g.ctx.GetHeader(header)
}

func (g *GinContext) Set(key string, value any) {
	if g.values == nil {
		g.values = make(map[string]any)
	}
	g.values[key] = value
}

func (g *GinContext) Get(key string) any {
	if g.values == nil {
		return nil
	}
	return g.values[key]
}
