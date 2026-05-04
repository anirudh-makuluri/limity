import { hostedLimiter, memoryLimiter, RateLimitOptions, RateLimitResult } from '@limity/core';

export async function checkRateLimit(options: RateLimitOptions): Promise<RateLimitResult> {
  const apiKey = process.env.LIMITY_API_KEY;
  const baseUrl = process.env.LIMITY_BASE_URL;

  if (apiKey) {
    return hostedLimiter(options, apiKey, baseUrl);
  }

  return memoryLimiter(options);
}
