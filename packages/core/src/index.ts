import { RateLimitOptions, RateLimitResult } from './types';
import { memoryLimiter } from './memory';
import { hostedLimiter } from './hosted';

export async function rateLimit(options: RateLimitOptions): Promise<RateLimitResult> {
  // Safely access process (Node.js environment)
  const apiKey = (globalThis as any).process?.env?.RATE_LIMIT_API_KEY;

  if (apiKey) {
    // Use hosted limiter with fallback to memory
    try {
      return await hostedLimiter(options, apiKey);
    } catch (err) {
      (globalThis as any).console?.error?.('Hosted limiter failed, falling back to memory:', err);
      return memoryLimiter(options);
    }
  }

  // Default: use memory limiter
  return memoryLimiter(options);
}

export * from './types';
export { memoryLimiter } from './memory';
export { hostedLimiter } from './hosted';
