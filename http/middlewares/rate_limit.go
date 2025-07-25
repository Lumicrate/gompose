package middlewares

import (
	"github.com/Lumicrate/gompose/http"
	"sync"
	"time"
)

var visitors = make(map[string]time.Time)
var mu sync.Mutex

func RateLimitMiddleware(limit time.Duration) http.MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(ctx http.Context) {
			ip := ctx.RemoteIP()
			mu.Lock()
			lastRequest, exists := visitors[ip]
			if exists && time.Since(lastRequest) < limit {
				mu.Unlock()
				ctx.JSON(429, map[string]string{"error": "rate limit exceeded"})
				ctx.Abort()
				return
			}
			visitors[ip] = time.Now()
			mu.Unlock()
			ctx.Next()

			next(ctx)
		}
	}
}
