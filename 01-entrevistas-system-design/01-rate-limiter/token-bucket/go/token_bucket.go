package main

import (
	"math"
	"sync"
	"time"
)

// TokenBucket Rate Limiter
// Reference: https://en.wikipedia.org/wiki/Token_bucket
type TokenBucket struct {
	capacity            float64                   // Maximum number of tokens in the bucket
	tokens              float64                   // Current number of tokens
	refillRate          float64                   // Number of tokens to add per second
	lastRefillTimestamp time.Time                 // Timestamp of the last refill
	mu                  sync.Mutex                // Mutex for thread-safe access
}

// NewTokenBucket creates a new TokenBucket with the given capacity and refill rate.
func NewTokenBucket(capacity float64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:            capacity,
		tokens:              capacity, // start full
		refillRate:          refillRate,
		lastRefillTimestamp: time.Now(),
	}
}

// GetTokens returns the current number of tokens.
func (tb *TokenBucket) GetTokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return tb.tokens
}

// AllowRequest checks if a request can be allowed.
// It returns true if allowed, false otherwise.
func (tb *TokenBucket) AllowRequest() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refillTokens()

	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true // request allowed
	}

	return false // request denied: 429 Too Many Requests
}

// refillTokens updates the token count based on elapsed time.
func (tb *TokenBucket) refillTokens() {
	now := time.Now()
	elapsedTime := now.Sub(tb.lastRefillTimestamp).Seconds()

	// Calculate how many tokens to add based on elapsed time and refill rate
	tokensToAdd := elapsedTime * tb.refillRate

	// Refill tokens, but never exceed capacity
	tb.tokens = math.Min(tb.capacity, tb.tokens+tokensToAdd)
	tb.lastRefillTimestamp = now
}
