import { TokenBucket } from "./token-bucket";

console.log("Hello...");

// This file runs several clients with different rates to simulate a rate limiter
let capacity = 5;
let refillRate = 2;

const highCapacityLowRefillRate = new TokenBucket(capacity, refillRate); // 5 capacity, 2 tokens/sec refill rate
console.log(`Running a high capacity of ${capacity} tokens with a refill rate of ${refillRate} tokens per second...`);

for (let i = 0; i < 10; i++) {
    const allowed = highCapacityLowRefillRate.allowRequest();
    console.log(`High Capacity Low Refill Rate - Request ${i + 1}: ${allowed ? "ALLOWED" : "DENIED"}`);
}

// 2 capacity, 5 tokens/sec refill rate
capacity = 2;
refillRate = 5;
const lowCapacityHighRefillRate = new TokenBucket(capacity, refillRate);
console.log(`\n\nRunning a low capacity of ${capacity} tokens with a refill rate of ${refillRate} tokens per second...`);

for (let i = 0; i < 10; i++) {
    const allowed = lowCapacityHighRefillRate.allowRequest();
    console.log(`Low Capacity High Refill Rate - Request ${i + 1}: ${allowed ? "ALLOWED" : "DENIED"}`);
}

