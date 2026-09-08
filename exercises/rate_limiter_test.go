package exercises

import (
	"testing"
	"time"
)

type RateLimiter struct {
	Requests map[string][]Request
	Limit    int
	Window   time.Duration
}

type Request struct {
	UserID    string
	CreatedAt time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		Limit:    limit,
		Window:   window,
		Requests: make(map[string][]Request),
	}
}

func (r *RateLimiter) Allow(key string) bool {
	return r.AllowAt(key, time.Now())
}

func (r *RateLimiter) AllowAt(key string, now time.Time) bool {
	if key == "" {
		return false
	}

	requests := r.Requests[key]

	requestsOnWindow := numberOfRequestsInWindow(requests, now.Add(-r.Window), now)

	if len(requestsOnWindow) >= r.Limit {
		return false
	}

	requestsOnWindow = append(requestsOnWindow, Request{UserID: key, CreatedAt: now})
	r.Requests[key] = requestsOnWindow

	return true
}

func numberOfRequestsInWindow(requests []Request, timeLimit time.Time, timeTo time.Time) []Request {
	requestsOnWindow := make([]Request, 0)
	for _, request := range requests {
		if request.CreatedAt.After(timeLimit) && request.CreatedAt.Before(timeTo) {
			requestsOnWindow = append(requestsOnWindow, request)
		}
	}
	return requestsOnWindow
}

func TestRateLimiterAllowsUpToLimitWithinWindow(t *testing.T) {
	r := NewRateLimiter(3, time.Minute)
	base := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)

	if !r.AllowAt("u1", base) {
		t.Fatalf("expected first request to be allowed")
	}
	if !r.AllowAt("u1", base.Add(10*time.Second)) {
		t.Fatalf("expected second request to be allowed")
	}
	if !r.AllowAt("u1", base.Add(20*time.Second)) {
		t.Fatalf("expected third request to be allowed")
	}
}

func TestRateLimiterBlocksWhenLimitExceededWithinWindow(t *testing.T) {
	r := NewRateLimiter(2, time.Minute)
	base := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)

	if !r.AllowAt("u1", base) {
		t.Fatalf("expected first request to be allowed")
	}
	if !r.AllowAt("u1", base.Add(15*time.Second)) {
		t.Fatalf("expected second request to be allowed")
	}
	if r.AllowAt("u1", base.Add(30*time.Second)) {
		t.Fatalf("expected third request in same window to be blocked")
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	r := NewRateLimiter(2, time.Minute)
	base := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)

	if !r.AllowAt("u1", base) {
		t.Fatalf("expected first request to be allowed")
	}
	if !r.AllowAt("u1", base.Add(10*time.Second)) {
		t.Fatalf("expected second request to be allowed")
	}
	if r.AllowAt("u1", base.Add(20*time.Second)) {
		t.Fatalf("expected request inside same window to be blocked")
	}

	if !r.AllowAt("u1", base.Add(time.Minute+time.Second)) {
		t.Fatalf("expected request after window reset to be allowed")
	}
}

func TestRateLimiterIsolatedPerKey(t *testing.T) {
	r := NewRateLimiter(1, time.Minute)
	base := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)

	if !r.AllowAt("u1", base) {
		t.Fatalf("expected u1 first request to be allowed")
	}
	if r.AllowAt("u1", base.Add(5*time.Second)) {
		t.Fatalf("expected u1 second request in same window to be blocked")
	}

	if !r.AllowAt("u2", base.Add(5*time.Second)) {
		t.Fatalf("expected u2 first request to be allowed independently")
	}
}

func TestRateLimiterBoundaryAtExactWindow(t *testing.T) {
	r := NewRateLimiter(1, time.Minute)
	base := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)

	if !r.AllowAt("u1", base) {
		t.Fatalf("expected first request to be allowed")
	}
	if !r.AllowAt("u1", base.Add(time.Minute)) {
		t.Fatalf("expected request at exact window boundary to be allowed")
	}
}
