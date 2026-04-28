import { RateLimitOptions, RateLimitResult } from './types';

const DEFAULT_LIMIT = 100;
const DEFAULT_WINDOW = 60;

export async function hostedLimiter(
  options: RateLimitOptions,
  apiKey: string,
  apiUrl: string = 'https://api.limity.dev'
): Promise<RateLimitResult> {
  const { key, limit = DEFAULT_LIMIT, window = DEFAULT_WINDOW } = options;

  try {
    const response = await fetch(`${apiUrl}/check`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${apiKey}`,
      },
      body: JSON.stringify({
        key,
        limit,
        window,
      }),
    });

    if (!response.ok) {
      // On API error, we don't throw - just return false
      console.error(`Rate limit API error: ${response.status}`);
      return {
        allowed: false,
        remaining: 0,
        reset: Math.floor(Date.now() / 1000) + window,
      };
    }

    const data = await response.json() as RateLimitResult;
    return data;
  } catch (err) {
    // Network error - don't throw
    console.error('Rate limit API request failed:', err);
    return {
      allowed: false,
      remaining: 0,
      reset: Math.floor(Date.now() / 1000) + window,
    };
  }
}
