package main

import (
	"testing"
	"time"
)

func TestInitialization(t *testing.T) {
	tb := NewTokenBucket(10, 1)
	if tb.GetTokens() != 10 {
		t.Errorf("Expected 10 tokens, got %f", tb.GetTokens())
	}
}

func TestAllowRequest(t *testing.T) {
	// 5 capacity, 2 tokens/sec refill rate
	tb := NewTokenBucket(5, 2)

	// Should allow requests while tokens are available
	if !tb.AllowRequest() {
		t.Error("Request 1 should be allowed")
	}
	if !tb.AllowRequest() {
		t.Error("Request 2 should be allowed")
	}
	if !tb.AllowRequest() {
		t.Error("Request 3 should be allowed")
	}
}

func TestDenyRequest(t *testing.T) {
	// 5 capacity, 2 tokens/sec refill rate
	tb := NewTokenBucket(5, 2)

	// Consume all tokens
	for i := 0; i < 5; i++ {
		tb.AllowRequest()
	}

	// Next request should be denied
	if tb.AllowRequest() {
		t.Error("Request should be denied when no tokens available")
	}
}

func TestRefillTokens(t *testing.T) {
	// 5 capacity, 2 tokens/sec refill rate
	tb := NewTokenBucket(5, 2)

	// Consume all tokens
	for i := 0; i < 5; i++ {
		tb.AllowRequest()
	}

	if tb.AllowRequest() {
		t.Error("Request should be denied immediately after empty")
	}

	// Wait 1 second (refill rate is 2 tokens/sec = 2 tokens available)
	time.Sleep(1 * time.Second)

	// Should have at least 1 token to allow a request
	if !tb.AllowRequest() {
		t.Error("Request should be allowed after refill")
	}
}

func TestNotExceedCapacity(t *testing.T) {
	// Start with 5 tokens (full capacity)
	tb := NewTokenBucket(5, 2)

	if tb.GetTokens() != 5 {
		t.Errorf("Expected 5 tokens, got %f", tb.GetTokens())
	}

	// Wait 2 seconds (would add 4 tokens, but limited by capacity)
	time.Sleep(2 * time.Second)

	// Should still be at capacity (5)
	if tb.GetTokens() > 5 {
		t.Errorf("Tokens should not exceed capacity. Got %f", tb.GetTokens())
	}
}

func TestBurstTraffic(t *testing.T) {
	// 10 capacity, 5 tokens/sec
	tb := NewTokenBucket(10, 5)

	// Allow 10 requests (burst)
	for i := 0; i < 10; i++ {
		if !tb.AllowRequest() {
			t.Errorf("Request %d in burst should be allowed", i+1)
		}
	}

	// Next 5 requests should fail
	for i := 0; i < 5; i++ {
		if tb.AllowRequest() {
			t.Errorf("Request %d in burst overflow should be denied", i+1)
		}
	}
}

func TestGradualRefill(t *testing.T) {
	// 3 capacity, 1 token/sec
	tb := NewTokenBucket(3, 1)

	// Consume all 3 tokens
	for i := 0; i < 3; i++ {
		tb.AllowRequest()
	}

	// Request is denied
	if tb.AllowRequest() {
		t.Error("Request should be denied")
	}

	// Wait 500ms (0.5 tokens refilled)
	time.Sleep(500 * time.Millisecond)

	// Still should be denied (need at least 1 full token)
	if tb.AllowRequest() {
		t.Error("Request should still be denied after partial refill")
	}

	// Wait another 600ms (total 1.1 tokens refilled)
	time.Sleep(600 * time.Millisecond)

	// Now should allow request
	if !tb.AllowRequest() {
		t.Error("Request should be allowed after full token refill")
	}
}
