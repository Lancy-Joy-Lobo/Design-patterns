package main

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	apiRequests map[string]int
	limiter     map[string][]time.Time
	window      time.Duration
}

func NewRateLimiter(window time.Duration) *RateLimiter {
	return &RateLimiter{
		apiRequests: make(map[string]int),
		limiter:     make(map[string][]time.Time),
		window:      window,
	}
}

func (r *RateLimiter) NewConfiguration(apiName string, limit int) {
	r.apiRequests[apiName] = limit
}

func (r *RateLimiter) Allow(apiPath, userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := apiPath + "_" + userID
	now := time.Now()
	timeStamps := r.limiter[key]
	valid := make([]time.Time, 0, len(timeStamps))

	// Keep only timestamps inside the window
	for _, t := range timeStamps {
		if now.Sub(t) <= r.window {
			valid = append(valid, t) // keep original timestamps
		}
	}

	// Reject if limit reached
	if len(valid) >= r.apiRequests[apiPath] {
		fmt.Println("rejected")
		r.limiter[key] = valid
		return false
	}

	// Accept and add new timestamp
	valid = append(valid, now)
	r.limiter[key] = valid
	fmt.Println("accepted")
	return true
}

func main() {
	// Sliding window = 60 seconds
	rateLimiter := NewRateLimiter(60 * time.Second)
	rateLimiter.NewConfiguration("/login", 10)
	rateLimiter.NewConfiguration("/submit", 5)

	apiPath := "/login"

	for i := 0; i < 15; i++ {
		allowed := rateLimiter.Allow(apiPath, "112")
		fmt.Println("Request", i, "->", allowed)
		if i == 12 {
			time.Sleep(61 * time.Second) // wait for window reset
		}
	}
}
