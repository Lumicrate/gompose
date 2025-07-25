package middlewares

import (
	"github.com/Lumicrate/gompose/http"
	"log"
	"time"
)

func LoggingMiddleware() http.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) {
			start := time.Now()
			ctx.Next()
			duration := time.Since(start)
			log.Printf("[%s] %s %s %d %s",
				ctx.Method(), ctx.Path(), ctx.RemoteIP(), ctx.Status(), duration)

			next(ctx)
		}
	}
}
