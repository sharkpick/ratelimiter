package ratelimiter

import (
	"context"
	"net/http"
	"sync/atomic"
)

var (
	connections    atomic.Int64
	MaxConnections int64 = 300
)

func ArriveWithContext(ctx context.Context) bool {
	n := connections.Add(1)
	go func() {
		<-ctx.Done()
		connections.Add(-1)
	}()
	return n <= MaxConnections
}

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if shouldArrive := ArriveWithContext(ctx); shouldArrive {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
}
