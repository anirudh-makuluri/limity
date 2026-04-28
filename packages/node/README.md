# Limity Node

Express middleware for rate limiting.

## Usage

```typescript
import express from 'express';
import { rateLimit } from '@limity/node';

const app = express();

// Basic usage (100 requests per 60 seconds per IP)
app.use(rateLimit());

// Custom configuration
app.use(rateLimit({
  limit: 50,
  window: 60,
  keyFn: (req) => req.user?.id || req.ip,
  skip: (req) => req.path === '/health',
  onLimitExceeded: (req, res, result) => {
    res.json({
      error: 'Rate limit exceeded',
      retryAfter: result.reset - Math.floor(Date.now() / 1000),
    });
  },
}));

app.get('/api/data', (req, res) => {
  const { remaining, reset } = (req as any).rateLimit;
  res.json({ data: 'hello', remaining, reset });
});

app.listen(3000);
```

## Configuration

- `keyFn` - Function to extract rate limit key from request (default: IP address)
- `limit` - Max requests per window (default: 100)
- `window` - Window size in seconds (default: 60)
- `skip` - Function to skip rate limiting for certain requests
- `onLimitExceeded` - Custom handler when limit is exceeded

## Headers

The middleware sets standard rate limit headers:

- `X-RateLimit-Limit` - Request limit per window
- `X-RateLimit-Remaining` - Remaining requests in window
- `X-RateLimit-Reset` - Unix timestamp when window resets

## Environment

Set `RATE_LIMIT_API_KEY` to use hosted rate limiting:

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

Without it, uses fast in-memory limiter.
