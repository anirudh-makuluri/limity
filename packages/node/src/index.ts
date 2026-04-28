import { Request, Response, NextFunction } from 'express';
import { rateLimit as coreLimiter, RateLimitOptions, RateLimitResult } from '@limity/core';

export interface RateLimitConfig {
  keyFn?: (req: Request) => string;
  limit?: number;
  window?: number;
  skip?: (req: Request) => boolean;
  onLimitExceeded?: (req: Request, res: Response, result: RateLimitResult) => void;
}

export function rateLimit(config: RateLimitConfig = {}) {
  const {
    keyFn = (req: Request) => req.ip || 'unknown',
    limit = 100,
    window = 60,
    skip = () => false,
    onLimitExceeded,
  } = config;

  return async (req: Request, res: Response, next: NextFunction) => {
    // Skip rate limiting if configured
    if (skip(req)) {
      return next();
    }

    const key = keyFn(req);
    const options: RateLimitOptions = { key, limit, window };

    try {
      const result = await coreLimiter(options);

      // Attach result to request for downstream access
      (req as any).rateLimit = result;

      // Set response headers
      res.set('X-RateLimit-Limit', limit.toString());
      res.set('X-RateLimit-Remaining', result.remaining.toString());
      res.set('X-RateLimit-Reset', result.reset.toString());

      if (!result.allowed) {
        const status = 429;
        res.status(status);

        if (onLimitExceeded) {
          onLimitExceeded(req, res, result);
        } else {
          res.json({
            error: 'Too many requests',
            retryAfter: result.reset - Math.floor(Date.now() / 1000),
          });
        }
        return;
      }

      next();
    } catch (err) {
      console.error('Rate limit error:', err);
      // On error, allow the request to pass
      next();
    }
  };
}

export * from '@limity/core';
