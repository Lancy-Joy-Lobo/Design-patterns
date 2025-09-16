package main

import (
	"fmt"
	"sync"
	"time"
)

type Bucket struct {
	Capacity   int
	Token      int
	RefillRate int
	LastRefill time.Time
}

type RateLimiter struct {
	Userhits map[string]*Bucket
	mu       sync.Mutex
}

func (r *RateLimiter) Allow(api string, userId string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := api + "_" + userId
	bucket := r.Userhits[key]

	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill)

	// Calculate tokens to add based on elapsed time and refill rate
	tokensToAdd := int(elapsed.Seconds()) * bucket.RefillRate

	// Update token count (capped at capacity) and last refill time
	bucket.Token = min(bucket.Capacity, bucket.Token+tokensToAdd)
	bucket.LastRefill = now

	if bucket.Token > 0 {
		fmt.Println("request accepted")
		bucket.Token -= 1
		return true
	}

	fmt.Println("request rejected")
	return false
}

func (r *RateLimiter) NewConfiguration(apiPath string, userId string, capacity int, refillRate int) {
	r.Userhits[apiPath+"_"+userId] = &Bucket{
		Capacity:   capacity,
		Token:      capacity,   // Start with full bucket
		RefillRate: refillRate, // tokens per second
		LastRefill: time.Now(),
	}
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		Userhits: map[string]*Bucket{},
	}
}

func main() {
	userId := "RB162554"
	api := "/submit"

	rateLimiter := NewRateLimiter()
	// Configure: capacity=5 tokens, refill rate=2 tokens per second
	rateLimiter.NewConfiguration(api, userId, 5, 2)

	// Test the rate limiter
	fmt.Println("Testing Token Bucket Rate Limiter")
	fmt.Println("Bucket capacity: 5 tokens, Refill rate: 2 tokens/second")
	fmt.Println()

	// Make rapid requests to exhaust tokens
	for i := 1; i <= 8; i++ {
		fmt.Printf("Request %d: ", i)
		allowed := rateLimiter.Allow(api, userId)
		if !allowed {
			fmt.Println()
		}
		time.Sleep(100 * time.Millisecond) // Small delay
	}

	fmt.Println("\nWaiting 3 seconds for token refill...")
	time.Sleep(3 * time.Second)

	// Test again after refill
	fmt.Println("\nAfter 3 seconds (should have ~6 tokens):")
	for i := 9; i <= 12; i++ {
		fmt.Printf("Request %d: ", i)
		allowed := rateLimiter.Allow(api, userId)
		if !allowed {
			fmt.Println()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
