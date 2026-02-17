// rate limit for the middleware

package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/Suthar345Piyush/invoicego/internal/util"
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

// allow user for requesting to the middleware

func (rl *RateLimiter) Allow(ip string) bool {

	rl.mu.Lock()

	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[ip]

	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: now,
			count:    1,
		}
		return true
	}

	//resetting the counter if window has passed

	if now.Sub(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = now
		return true
	}

	// incrementing the  counter

	v.count++
	v.lastSeen = now

	return v.count <= rl.limit
}

// main rate limiter function

func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {

	limiter := NewRateLimiter(limit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := r.RemoteAddr

			// extracting ip from X-Forwarded-For or X-Real-IP if behind proxy

			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				ip = realIP
			}

			if !limiter.Allow(ip) {
				util.WriteJSON(w, http.StatusTooManyRequests, util.Response{
					Success: false,
					Error:   "rate limit exceeded, please try again later",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
