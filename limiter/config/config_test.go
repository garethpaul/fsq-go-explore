package config

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestLimiterRefillsConfiguredMaximumAcrossTTL(t *testing.T) {
	limiter := NewLimiter(10, time.Minute)

	for i := 0; i < 10; i++ {
		if reached := limiter.LimitReached("client"); reached {
			t.Fatalf("request %d reached limit during initial burst", i+1)
		}
	}
	if reached := limiter.LimitReached("client"); !reached {
		t.Fatal("request after initial burst did not reach limit")
	}

	got := float64(limiter.tokenBuckets["client"].Limit())
	want := float64(10) / time.Minute.Seconds()
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("refill rate = %v requests/second, want %v", got, want)
	}
}

func TestLimiterRejectsInvalidRateConfiguration(t *testing.T) {
	tests := []struct {
		name string
		max  int64
		ttl  time.Duration
	}{
		{name: "zero maximum", max: 0, ttl: time.Minute},
		{name: "negative maximum", max: -1, ttl: time.Minute},
		{name: "zero duration", max: 10, ttl: 0},
		{name: "negative duration", max: 10, ttl: -time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewLimiter(tt.max, tt.ttl)
			if reached := limiter.LimitReached("client"); !reached {
				t.Fatal("invalid limiter configuration allowed a request")
			}
		})
	}
}

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
