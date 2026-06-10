package config

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterCapsTrackedKeys(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.maxTrackedKeys = 3

	for i := 0; i < 20; i++ {
		limiter.LimitReached(fmt.Sprintf("key-%d", i))

		if got := len(limiter.tokenBuckets); got > limiter.maxTrackedKeys {
			t.Fatalf("tracked %d keys, want at most %d", got, limiter.maxTrackedKeys)
		}
	}

	if got := limiter.tokenBucketOrder.Len(); got != limiter.maxTrackedKeys {
		t.Fatalf("LRU order contains %d keys, want %d", got, limiter.maxTrackedKeys)
	}
	if got := len(limiter.tokenBucketEntries); got != limiter.maxTrackedKeys {
		t.Fatalf("LRU index contains %d keys, want %d", got, limiter.maxTrackedKeys)
	}
}

func TestLimiterEvictsLeastRecentlyUsedKey(t *testing.T) {
	limiter := NewLimiter(1, time.Hour)
	limiter.maxTrackedKeys = 2

	limiter.LimitReached("oldest")
	limiter.LimitReached("recent")
	limiter.LimitReached("oldest")
	limiter.LimitReached("new")

	if _, found := limiter.tokenBuckets["recent"]; found {
		t.Fatal("least recently used key was retained")
	}
	if _, found := limiter.tokenBuckets["oldest"]; !found {
		t.Fatal("recently accessed key was evicted")
	}
	if reached := limiter.LimitReached("recent"); reached {
		t.Fatal("evicted key did not receive a fresh token bucket")
	}
}
