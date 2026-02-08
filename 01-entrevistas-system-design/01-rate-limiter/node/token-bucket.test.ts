import { TokenBucket } from "./token-bucket";

describe("TokenBucket Rate Limiter", () => {
  let tokenBucket: TokenBucket;

  describe("Initialization", () => {
    it("should initialize with full capacity", () => {
      tokenBucket = new TokenBucket(10, 1);
      expect(tokenBucket.getTokens()).toBe(10);
    });
  });

  describe("allowRequest", () => {
    beforeEach(() => {
      tokenBucket = new TokenBucket(5, 2); // 5 capacity, 2 tokens/sec refill rate
    });

    it("should allow requests while tokens are available", () => {
      expect(tokenBucket.allowRequest()).toBe(true);
      expect(tokenBucket.allowRequest()).toBe(true);
      expect(tokenBucket.allowRequest()).toBe(true);
    });

    it("should deny requests when no tokens are available", () => {
      // Consume all tokens
      for (let i = 0; i < 5; i++) {
        tokenBucket.allowRequest();
      }

      // Next request should be denied
      expect(tokenBucket.allowRequest()).toBe(false);
    });

    it("should refill tokens over time", async () => {
      // Consume all tokens
      for (let i = 0; i < 5; i++) {
        tokenBucket.allowRequest();
      }

      expect(tokenBucket.allowRequest()).toBe(false);

      // Wait 1 second (refill rate is 2 tokens/sec = 2 tokens available)
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Should have at least 1 token to allow a request
      expect(tokenBucket.allowRequest()).toBe(true);
    });

    it("should not exceed capacity after refilling", async () => {
      // Start with 5 tokens (full capacity)
      expect(tokenBucket.getTokens()).toBe(5);

      // Wait 2 seconds (would add 4 tokens, but limited by capacity)
      await new Promise((resolve) => setTimeout(resolve, 2000));

      // Should still be at capacity (5)
      expect(tokenBucket.getTokens()).toBeLessThanOrEqual(5);
    });
  });

  describe("Rate limiting scenarios", () => {
    it("should handle burst traffic correctly", () => {
      tokenBucket = new TokenBucket(10, 5); // 10 capacity, 5 tokens/sec

      // Allow 10 requests (burst)
      for (let i = 0; i < 10; i++) {
        expect(tokenBucket.allowRequest()).toBe(true);
      }

      // Next 5 requests should fail
      for (let i = 0; i < 5; i++) {
        expect(tokenBucket.allowRequest()).toBe(false);
      }
    });

    it("should gradually allow requests after refill", async () => {
      tokenBucket = new TokenBucket(3, 1); // 3 capacity, 1 token/sec

      // Consume all 3 tokens
      for (let i = 0; i < 3; i++) {
        expect(tokenBucket.allowRequest()).toBe(true);
      }

      // Request is denied
      expect(tokenBucket.allowRequest()).toBe(false);

      // Wait 500ms (0.5 tokens refilled)
      await new Promise((resolve) => setTimeout(resolve, 500));

      // Still should be denied (need at least 1 full token)
      expect(tokenBucket.allowRequest()).toBe(false);

      // Wait another 600ms (total 1.1 tokens refilled)
      await new Promise((resolve) => setTimeout(resolve, 600));

      // Now should allow request
      expect(tokenBucket.allowRequest()).toBe(true);
      expect(tokenBucket.allowRequest()).toBe(false);
    });
  });
});
