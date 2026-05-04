# @limity/node

**Express.js middleware for rate limiting.** Drop-in middleware that adds automatic rate limiting to your Express application.

## Installation

```bash
npm install @limity/node
# or
pnpm add @limity/node
# or
yarn add @limity/node
```

## Quick Start

```typescript
import express from 'express';
import { rateLimit } from '@limity/node';

const app = express();

// Add rate limiting to all routes
app.use(rateLimit());

app.get('/api/data', (req, res) => {
  res.json({ data: 'hello' });
});

app.listen(3000);
```

## How It Works

By default, the middleware:
- Rates limits by **IP address**
- Allows **100 requests** per **60 seconds**
- Returns **429 Too Many Requests** when limit exceeded
- Sets standard rate limit headers

## Configuration

Customize the middleware with options:

```typescript
app.use(rateLimit({
  limit: 50,              // Max requests per window
  window: 60,             // Window size in seconds
  keyFn: (req) => req.user?.id || req.ip,  // Extract rate limit key
  skip: (req) => req.path === '/health',   // Skip rate limiting
  onLimitExceeded: (req, res, result) => { // Custom error handler
    res.status(429).json({
      error: 'Too many requests',
      retryAfter: result.reset - Math.floor(Date.now() / 1000),
    });
  },
}));
```

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `limit` | `number` | 100 | Max requests per window |
| `window` | `number` | 60 | Window size in seconds |
| `keyFn` | `(req) => string` | IP address | Function to extract rate limit key |
| `skip` | `(req) => boolean` | undefined | Skip rate limiting for matching requests |
| `onLimitExceeded` | `(req, res, result) => void` | 429 response | Custom handler when limit exceeded |

## Examples

### Rate Limit by User ID

```typescript
app.use(rateLimit({
  keyFn: (req) => req.user?.id || 'anonymous',
  limit: 100,
  window: 60,
}));
```

### Different Limits for Different Routes

```typescript
// Strict limit for login
app.post('/auth/login', 
  rateLimit({
    limit: 5,
    window: 300, // 5 per 5 minutes
  }),
  (req, res) => {
    // Handle login
  }
);

// Relaxed limit for API
app.get('/api/data',
  rateLimit({
    limit: 1000,
    window: 60,
  }),
  (req, res) => {
    // Handle API request
  }
);
```

### Skip Rate Limiting for Specific Routes

```typescript
app.use(rateLimit({
  skip: (req) => req.path === '/health' || req.path === '/status',
}));
```

### Access Rate Limit Info in Handlers

```typescript
app.get('/api/data', (req, res) => {
  const { allowed, remaining, reset } = (req as any).rateLimit;
  
  res.set('X-RateLimit-Remaining', remaining.toString());
  res.set('X-RateLimit-Reset', reset.toString());
  
  res.json({
    data: 'hello',
    rateLimit: { remaining, reset },
  });
});
```

### Custom Error Response

```typescript
app.use(rateLimit({
  onLimitExceeded: (req, res, result) => {
    const retryAfter = result.reset - Math.floor(Date.now() / 1000);
    
    res.status(429)
      .set('Retry-After', retryAfter.toString())
      .json({
        error: 'Too many requests',
        message: `Try again in ${retryAfter} seconds`,
        retryAfter,
      });
  },
}));
```

## Response Headers

The middleware automatically sets these headers:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1682500000
```

## Rate Limit Data

The rate limit result is attached to the request object:

```typescript
app.get('/api/example', (req, res) => {
  const { allowed, remaining, reset } = (req as any).rateLimit;
  // {
  //   allowed: true,
  //   remaining: 45,
  //   reset: 1682500000
  // }
});
```

## Configuration

### Environment Variables

- **`LIMITY_API_KEY`** (optional) - API key for hosted mode. If not set, uses memory mode.

```bash
export LIMITY_API_KEY=your_api_key
```

## Memory vs Hosted Mode

### Memory Mode (Default)

```bash
# No API key needed
app.use(rateLimit());
```

Perfect for:
- Single server deployments
- Development and testing
- Small to medium traffic

### Hosted Mode (With API Key)

```bash
export LIMITY_API_KEY=your_api_key
```

Perfect for:
- Multi-server deployments
- Load-balanced applications
- High-traffic applications
- Distributed rate limiting

## When to Use

✅ Use `@limity/node` when:
- You're building an Express.js application
- You want automatic IP-based rate limiting
- You need middleware that's easy to set up

🔗 Use `@limity/core` instead if:
- You're building a custom HTTP framework
- You need maximum control over rate limiting
- You're not using Express.js

🌐 Use `@limity/edge` instead if:
- You're deploying to edge functions (Vercel, Cloudflare, etc.)
- You need rate limiting in serverless environments

## Performance

- **Memory mode**: ~1-2ms per request (in-process)
- **Hosted mode**: ~50-200ms per request (network + API)
