package handlers

import (
	"sync"
	"time"
)

type slidingWindowLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	attempts map[string][]time.Time
	now      func() time.Time
	lastGC   time.Time
}

func newSlidingWindowLimiter(limit int, window time.Duration) *slidingWindowLimiter {
	return &slidingWindowLimiter{
		limit:    limit,
		window:   window,
		attempts: make(map[string][]time.Time),
		now:      time.Now,
	}
}

func (limiter *slidingWindowLimiter) Allow(key string) bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := limiter.now()
	cutoff := now.Add(-limiter.window)
	if limiter.lastGC.IsZero() || now.Sub(limiter.lastGC) >= limiter.window {
		for attemptKey, timestamps := range limiter.attempts {
			if len(timestamps) == 0 || !timestamps[len(timestamps)-1].After(cutoff) {
				delete(limiter.attempts, attemptKey)
			}
		}
		limiter.lastGC = now
	}
	existing := limiter.attempts[key]
	firstActive := 0
	for firstActive < len(existing) && !existing[firstActive].After(cutoff) {
		firstActive++
	}
	active := existing[firstActive:]
	if len(active) >= limiter.limit {
		limiter.attempts[key] = active
		return false
	}

	limiter.attempts[key] = append(active, now)
	return true
}
