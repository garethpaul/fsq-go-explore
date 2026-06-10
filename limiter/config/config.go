package config

import (
	"container/list"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const defaultMaxTrackedKeys = 10000

// NewLimiter is a constructor for Limiter.
func NewLimiter(max int64, ttl time.Duration) *Limiter {
	limiter := &Limiter{Max: max, TTL: ttl}
	limiter.MessageContentType = "text/plain; charset=utf-8"
	limiter.Message = "You have reached the maximum request limit for this tool"
	limiter.StatusCode = 429
	limiter.tokenBuckets = make(map[string]*rate.Limiter)
	limiter.tokenBucketOrder = list.New()
	limiter.tokenBucketEntries = make(map[string]*list.Element)
	limiter.maxTrackedKeys = defaultMaxTrackedKeys
	limiter.IPLookups = []string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}

	return limiter
}

// Limiter is a config struct to limit a particular request handler.
type Limiter struct {
	// HTTP message when limit is reached.
	Message string

	// Content-Type for Message
	MessageContentType string

	// HTTP status code when limit is reached.
	StatusCode int

	// Maximum number of requests to limit per duration.
	Max int64

	// Duration of rate-limiter.
	TTL time.Duration

	// List of places to look up IP address.
	// Default is "RemoteAddr", "X-Forwarded-For", "X-Real-IP".
	// You can rearrange the order as you like.
	IPLookups []string

	// List of HTTP Methods to limit (GET, POST, PUT, etc.).
	// Empty means limit all methods.
	Methods []string

	// List of HTTP headers to limit.
	// Empty means skip headers checking.
	Headers map[string][]string

	// List of basic auth usernames to limit.
	BasicAuthUsers []string

	// Throttler struct
	tokenBuckets map[string]*rate.Limiter

	// LRU bookkeeping bounds request-controlled rate-limiter keys.
	tokenBucketOrder   *list.List
	tokenBucketEntries map[string]*list.Element
	maxTrackedKeys     int

	sync.RWMutex
}

// LimitReached returns a bool indicating if the Bucket identified by key ran out of tokens.
func (l *Limiter) LimitReached(key string) bool {
	l.Lock()
	defer l.Unlock()

	bucket, found := l.tokenBuckets[key]
	if found {
		l.tokenBucketOrder.MoveToFront(l.tokenBucketEntries[key])
	} else {
		if l.maxTrackedKeys > 0 && len(l.tokenBuckets) >= l.maxTrackedKeys {
			oldest := l.tokenBucketOrder.Back()
			oldestKey := oldest.Value.(string)
			delete(l.tokenBuckets, oldestKey)
			delete(l.tokenBucketEntries, oldestKey)
			l.tokenBucketOrder.Remove(oldest)
		}

		bucket = rate.NewLimiter(rate.Every(l.TTL), int(l.Max))
		l.tokenBuckets[key] = bucket
		l.tokenBucketEntries[key] = l.tokenBucketOrder.PushFront(key)
	}

	return !bucket.AllowN(time.Now(), 1)
}
