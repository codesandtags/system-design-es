/**
 * Token Bucket Rate Limiter using Node.js (Single Threaded)
 * Reference: https://en.wikipedia.org/wiki/Token_bucket
 *
 */

export class TokenBucket {
  private capacity: number; // maximum number of tokens in the bucket
  private tokens: number;
  private refillRate: number; // number of tokens to add per second
  private lastRefillTimestamp: number;

  constructor(capacity: number, refillRate: number) {
    this.capacity = capacity;
    this.tokens = capacity; // start full
    this.refillRate = refillRate;
    this.lastRefillTimestamp = Date.now();
  }

  getTokens(): number {
    return this.tokens;
  }

  /**
   * This method is called for each incoming request to check if it can be allowed.
   *
   * @returns
   */
  allowRequest(): boolean {
    this.refillTokens();

    if (this.tokens >= 1) {
      this.tokens -= 1;
      return true; // request allowed
    }

    return false; // request denied: 429 Too Many Requests
  }

  /**
   * Refill tokens based on the elapsed time since the last refill.
   */
  private refillTokens() {
    const now = Date.now();
    const elapsedTime = (now - this.lastRefillTimestamp) / 1000;

    // Calculate how many tokens to add based on elapsed time and refill rate
    const tokensToAdd = elapsedTime * this.refillRate;

    // Refill tokens, but never exceed capacity
    this.tokens = Math.min(this.capacity, this.tokens + tokensToAdd);
    this.lastRefillTimestamp = now;
  }
}
