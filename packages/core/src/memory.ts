import { RateLimitOptions, RateLimitResult, WindowEntry } from './types';

const DEFAULT_LIMIT = 100;
const DEFAULT_WINDOW = 60;

// In-memory store: key -> { count, reset }
const memoryStore = new Map<string, WindowEntry>();

export async function memoryLimiter(
  options: RateLimitOptions
): Promise<RateLimitResult> {
  const { key, limit = DEFAULT_LIMIT, window = DEFAULT_WINDOW } = options;

  const now = Math.floor(Date.now() / 1000);
  const windowStart = now - (now % window);
  const reset = windowStart + window;

  const storeKey = `${key}:${windowStart}`;
  const entry = memoryStore.get(storeKey);

  if (!entry) {
    // First request in window
    memoryStore.set(storeKey, { count: 1, reset });
    return {
      allowed: true,
      remaining: limit - 1,
      reset,
    };
  }

  // Clean up old entries
  if (entry.reset <= now) {
    memoryStore.delete(storeKey);
    memoryStore.set(storeKey, { count: 1, reset });
    return {
      allowed: true,
      remaining: limit - 1,
      reset,
    };
  }

  // Increment counter
  entry.count += 1;

  const allowed = entry.count <= limit;
  const remaining = Math.max(0, limit - entry.count);

  return {
    allowed,
    remaining,
    reset,
  };
}
