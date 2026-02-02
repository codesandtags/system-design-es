package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello...")

	// This file runs several clients with different rates to simulate a rate limiter
	capacity := 5.0
	refillRate := 2.0

	highCapacityLowRefillRate := NewTokenBucket(capacity, refillRate) // 5 capacity, 2 tokens/sec refill rate
	fmt.Printf("Running a high capacity of %.0f tokens with a refill rate of %.0f tokens per second...\n", capacity, refillRate)

	for i := 0; i < 10; i++ {
		allowed := highCapacityLowRefillRate.AllowRequest()
		status := "DENIED"
		if allowed {
			status = "ALLOWED"
		}
		fmt.Printf("High Capacity Low Refill Rate - Request %d: %s\n", i+1, status)
	}

	// 2 capacity, 5 tokens/sec refill rate
	capacity = 2.0
	refillRate = 5.0
	lowCapacityHighRefillRate := NewTokenBucket(capacity, refillRate)
	fmt.Printf("\n\nRunning a low capacity of %.0f tokens with a refill rate of %.0f tokens per second...\n", capacity, refillRate)

	for i := 0; i < 10; i++ {
		allowed := lowCapacityHighRefillRate.AllowRequest()
		status := "DENIED"
		if allowed {
			status = "ALLOWED"
		}
		fmt.Printf("Low Capacity High Refill Rate - Request %d: %s\n", i+1, status)
	}
}
