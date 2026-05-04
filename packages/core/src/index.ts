import { RateLimitOptions, RateLimitResult } from './types.js';
import { memoryLimiter } from './memory.js';
import { hostedLimiter } from './hosted.js';

export async function rateLimit(options: RateLimitOptions): Promise<RateLimitResult> {
  // Safely access process (Node.js environment)
  const env = (globalThis as any).process?.env;
  const apiKey = env?.LIMITY_API_KEY ?? env?.RATE_LIMIT_API_KEY;
  const baseUrl = env?.LIMITY_BASE_URL;

  if (apiKey) {
    // Use hosted limiter with fallback to memory
    try {
      return await hostedLimiter(options, apiKey, baseUrl);
    } catch (err) {
      (globalThis as any).console?.error?.('Hosted limiter failed, falling back to memory:', err);
      return memoryLimiter(options);
    }
  }

  // Default: use memory limiter
  return memoryLimiter(options);
}

export * from './types.js';
export { memoryLimiter } from './memory.js';
export { hostedLimiter } from './hosted.js';
