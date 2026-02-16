// rate limit for the middleware

package middleware

import (
	"sync"
	"time"
)

// limit the user/visitor

type visitor struct {
	lastSeen time.Time
	count    int
}

// here we setup the timeframe of rate limiter , and limit on the request a user/visitor can make

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex  // reader,writer mutual exclusion lock , if value is zero than it unlocked
	limit    int           // request per window a single user can make
	window   time.Duration // timeframe
}

// new rate limiter function

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}

	// cleaning all previous old entry every minute

	go rl.cleanup()

	return rl
}

// cleanup function to clear all the old entry every minute

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)

	defer ticker.Stop()

	for range ticker.C {

		rl.mu.Lock()

		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}

		rl.mu.Unlock()

	}

}
