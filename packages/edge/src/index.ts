import { rateLimit as coreLimiter, RateLimitResult } from '@limity/core';

export interface EdgeRequest {
  ip?: string;
  headers?: Record<string, string>;
  url?: string;
}

export async function checkRateLimit(
  request: Request | EdgeRequest,
  options?: {
    keyFn?: (req: Request | EdgeRequest) => string;
    limit?: number;
    window?: number;
  }
): Promise<RateLimitResult> {
  const {
    keyFn = (req: Request | EdgeRequest) => {
      if (req instanceof Request) {
        return req.headers.get('x-forwarded-for') || 
               req.headers.get('cf-connecting-ip') ||
               'unknown';
      }
      return req.ip || 'unknown';
    },
    limit = 100,
    window = 60,
  } = options || {};

  const key = keyFn(request);

  return coreLimiter({
    key,
    limit,
    window,
  });
}

export * from '@limity/core';
