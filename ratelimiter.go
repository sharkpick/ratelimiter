package ratelimiter

import (
	"net/http"
	"strings"
	"sync/atomic"
)

var (
	TheRateLimiter                                                     = &RateLimiter{}
	RateLimiterLimit                                             int64 = 100
	PathPrefixesBypassRateLimiter, PathSuffixesBypassRateLimiter []string
)

func PathShouldBypassRateLimiter(path string) bool {
	for _, prefix := range PathPrefixesBypassRateLimiter {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	for _, suffix := range PathSuffixesBypassRateLimiter {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}

func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shouldArrive := TheRateLimiter.Arrive()
		defer TheRateLimiter.Depart()
		if shouldArrive || PathShouldBypassRateLimiter(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
	})
}

type RateLimiter struct {
	n atomic.Int64
}

func (r *RateLimiter) Arrive() bool {
	n := r.n.Add(1)
	return n <= RateLimiterLimit
}

func (r *RateLimiter) Depart() {
	r.n.Add(-1)
}

func (r *RateLimiter) Count() int64 {
	return r.n.Load()
}
