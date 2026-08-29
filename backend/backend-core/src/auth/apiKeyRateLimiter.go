package auth

import (
	"sync"
	"time"
)

type apiKeyRateLimitBucket struct {
	windowStart time.Time
	count       uint32
	lastSeen    time.Time
}

var apiKeyRateLimitState = struct {
	mu      sync.Mutex
	buckets map[string]apiKeyRateLimitBucket
}{
	buckets: make(map[string]apiKeyRateLimitBucket),
}

func allowAPIKeyRequest(apiKeyUID string, limit *uint32) bool {
	if limit == nil || *limit == 0 {
		return true
	}
	now := time.Now()
	apiKeyRateLimitState.mu.Lock()
	defer apiKeyRateLimitState.mu.Unlock()

	for uid, bucket := range apiKeyRateLimitState.buckets {
		if now.Sub(bucket.lastSeen) > 10*time.Minute {
			delete(apiKeyRateLimitState.buckets, uid)
		}
	}

	bucket := apiKeyRateLimitState.buckets[apiKeyUID]
	if bucket.windowStart.IsZero() || now.Sub(bucket.windowStart) >= time.Minute {
		bucket.windowStart = now
		bucket.count = 0
	}
	bucket.lastSeen = now
	if bucket.count >= *limit {
		apiKeyRateLimitState.buckets[apiKeyUID] = bucket
		return false
	}
	bucket.count++
	apiKeyRateLimitState.buckets[apiKeyUID] = bucket
	return true
}
