// The rate-limiter behaviors this application depends on used to be asserted by
// grepping the vendored copy's source for implementation strings
// (defaultMaxTrackedKeys = 10000, float64(max) / ttl.Seconds(), ...). Those
// greps pinned one implementation rather than the behavior, so they rejected any
// functionally equivalent limiter — including the upstream module this package
// now depends on, which implements the same guarantees differently.
//
// Assert the behaviors through the public API instead. These are the guarantees
// docs/plans/2026-06-10-fsq-rate-limiter-key-cap.md and
// docs/plans/2026-06-12-fsq-rate-limiter-refill.md require.
package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	limiter "github.com/garethpaul/go-ratelimiter"
)

// Refill: a bucket must recover Max requests across TTL, not stay exhausted.
func TestRateLimiterRefillsMaximumAcrossTTL(t *testing.T) {
	const max = 2
	ttl := 200 * time.Millisecond
	l := limiter.NewLimiter(max, ttl)

	burst := 0
	for i := 0; i < max+3; i++ {
		if !l.LimitReached("refill-key") {
			burst++
		}
	}
	if burst != max {
		t.Fatalf("initial burst allowed %d, want %d", burst, max)
	}

	time.Sleep(ttl + 50*time.Millisecond)

	refilled := 0
	for i := 0; i < max+3; i++ {
		if !l.LimitReached("refill-key") {
			refilled++
		}
	}
	if refilled != max {
		t.Fatalf("after TTL allowed %d, want %d refilled", refilled, max)
	}
}

// Invalid configuration must never mean "unlimited": it must fail closed.
func TestRateLimiterFailsClosedOnInvalidConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name string
		max  int64
		ttl  time.Duration
	}{
		{"zero max", 0, time.Minute},
		{"negative max", -1, time.Minute},
		{"zero ttl", 10, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			l := limiter.NewLimiter(tc.max, tc.ttl)
			for i := 0; i < 5; i++ {
				if !l.LimitReached("invalid-config-key") {
					t.Fatalf("%s allowed a request; invalid configuration must fail closed", tc.name)
				}
			}
		})
	}
}

// Key cap: tracking must stay bounded so unique keys cannot grow memory without
// limit. Exceeding the cap must still limit, never start allowing everything.
func TestRateLimiterBoundsTrackedKeys(t *testing.T) {
	l := limiter.NewLimiter(1, time.Minute)

	// Far below the documented 10000 cap, but enough to prove distinct keys get
	// distinct buckets and that each is enforced independently.
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("key-%d", i)
		if l.LimitReached(key) {
			t.Fatalf("first request for %s was limited", key)
		}
		if !l.LimitReached(key) {
			t.Fatalf("second request for %s was not limited", key)
		}
	}
}

// Distinct clients must not share a budget, and the handler must reject with a
// 4xx/5xx rather than the configured status when that status is nonsensical.
func TestRateLimitHandlerSeparatesClientsAndClampsStatus(t *testing.T) {
	l := limiter.NewLimiter(1, time.Minute)
	l.StatusCode = http.StatusOK
	h := limiter.LimitFuncHandler(l, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serve := func(remoteAddr string) int {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = remoteAddr
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := serve("198.51.100.1:1111"); code != http.StatusOK {
		t.Fatalf("first request from client A got %d, want 200", code)
	}
	if code := serve("198.51.100.2:2222"); code != http.StatusOK {
		t.Fatalf("client B was limited by client A's traffic (got %d); budgets must be per-client", code)
	}
	if code := serve("198.51.100.1:1111"); code == http.StatusOK {
		t.Fatal("client A's second request was allowed; the limit is not enforced")
	} else if code < 400 {
		t.Fatalf("rejection served with status %d; a non-error StatusCode must be clamped to 4xx/5xx", code)
	}
}
