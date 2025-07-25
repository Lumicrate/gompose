package ginadapter

import (
	"fmt"
	"github.com/Lumicrate/gompose/http"
	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	engine *gin.Engine
	port   int
}

func New(port int) *GinEngine {
	return &GinEngine{
		engine: gin.Default(),
		port:   port,
	}
}

func (g *GinEngine) Init(_ int) error {
	return nil
}

func (g *GinEngine) RegisterRoute(method string, path string, handler http.HandlerFunc) {
	ginHandler := func(c *gin.Context) {
		handler(&GinContext{ctx: c})
	}

	switch method {
	case "GET":
		g.engine.GET(path, ginHandler)
	case "POST":
		g.engine.POST(path, ginHandler)
	case "PUT":
		g.engine.PUT(path, ginHandler)
	case "PATCH":
		g.engine.PATCH(path, ginHandler)
	case "DELETE":
		g.engine.DELETE(path, ginHandler)
	default:
		panic(fmt.Sprintf("Unsupported method: %s", method))
	}
}

func (g *GinEngine) Use(middleware http.MiddlewareFunc) {
	g.engine.Use(func(c *gin.Context) {
		final := middleware(func(ctx http.Context) {})
		final(&GinContext{ctx: c})
	})
}

func (g *GinContext) QueryParams() map[string][]string {
	return g.ctx.Request.URL.Query()
}

func (g *GinContext) BindJSON(obj any) error {
	return g.ctx.ShouldBindJSON(obj)
}

func (g *GinEngine) Start() error {
	return g.engine.Run(fmt.Sprintf(":%d", g.port))
}
