package ratelimiter

import (
	"testing"
)

func TestArrive(t *testing.T) {
	for i := int64(0); i < RateLimiterLimit; i++ {
		if shouldArrive := TheRateLimiter.Arrive(); !shouldArrive {
			t.Fatalf("error: expected true but got false")
		}
	}
	if TheRateLimiter.Count() != RateLimiterLimit {
		t.Fatalf("error: expected TheRateLimiter.Count == %d but got %d", RateLimiterLimit, TheRateLimiter.Count())
	}
	shouldArrive := TheRateLimiter.Arrive()
	if shouldArrive {
		t.Fatalf("error: expected false but got true")
	}
	if TheRateLimiter.Count() != RateLimiterLimit+1 {
		t.Fatalf("error: expected TheRateLimiter.Count == %d but got %d", RateLimiterLimit+1, TheRateLimiter.Count())
	}
}

func TestDepart(t *testing.T) {
	target := TheRateLimiter.Count()
	for i := int64(0); i < target; i++ {
		TheRateLimiter.Depart()
	}
	want := int64(0)
	if got := TheRateLimiter.Count(); got != want {
		t.Fatalf("error: expected TheRateLimiter.Count() == 0, got %d", got)
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			TheRateLimiter.Arrive()
			defer TheRateLimiter.Depart()
		}()
	}
}
