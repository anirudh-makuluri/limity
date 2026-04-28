# @limity/edge

**Rate limiting for edge functions and serverless environments.** Works with Vercel Edge Functions, Cloudflare Workers, and other edge runtimes.

## Installation

```bash
npm install @limity/edge
# or
pnpm add @limity/edge
# or
yarn add @limity/edge
```

## Quick Start

### Vercel Edge Functions

```typescript
import { checkRateLimit } from '@limity/edge';

export default async function handler(req: Request) {
  const result = await checkRateLimit(req, {
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    return new Response('Too many requests', { status: 429 });
  }

  return new Response('Hello world!');
}

export const config = {
  runtime: 'edge',
};
```

### Cloudflare Workers

```typescript
import { checkRateLimit } from '@limity/edge';

export default {
  async fetch(request: Request): Promise<Response> {
    const result = await checkRateLimit(request, {
      limit: 100,
      window: 60,
    });

    if (!result.allowed) {
      return new Response('Too many requests', { status: 429 });
    }

    return new Response('Hello world!');
  },
};
```

## How It Works

By default, the middleware:
- Rates limits by **IP address** (auto-detected from headers)
- Allows **100 requests** per **60 seconds**
- Returns **429 Too Many Requests** when limit exceeded
- Works with standard Fetch API (Request/Response)

## API Reference

### `checkRateLimit(request, options)`

Checks if a request should be allowed based on rate limits.

**Parameters:**

```typescript
checkRateLimit(
  request: Request,
  options?: {
    limit?: number;              // Max requests per window (default: 100)
    window?: number;             // Window size in seconds (default: 60)
    keyFn?: (req: Request) => string;  // Extract rate limit key
  }
): Promise<RateLimitResult>
```

**Returns:**

```typescript
interface RateLimitResult {
  allowed: boolean;  // True if request is allowed
  remaining: number; // Requests remaining in window
  reset: number;     // Unix timestamp when window resets
}
```

## Examples

### Custom IP Detection (Cloudflare)

```typescript
import { checkRateLimit } from '@limity/edge';

export default {
  async fetch(request: Request) {
    const result = await checkRateLimit(request, {
      limit: 100,
      window: 60,
      keyFn: (req) => req.headers.get('cf-connecting-ip') || 'unknown',
    });

    if (!result.allowed) {
      return new Response('Too many requests', { status: 429 });
    }

    return new Response('Hello!');
  },
};
```

### Custom IP Detection (Vercel)

```typescript
import { checkRateLimit } from '@limity/edge';

export default async function handler(req: Request) {
  const result = await checkRateLimit(req, {
    limit: 50,
    window: 60,
    keyFn: (r) => r.headers.get('x-forwarded-for')?.split(',')[0] || 'unknown',
  });

  if (!result.allowed) {
    return new Response('Too many requests', { status: 429 });
  }

  return new Response('Hello!');
}

export const config = { runtime: 'edge' };
```

### Rate Limit by Path

```typescript
import { checkRateLimit } from '@limity/edge';

export default async function handler(req: Request) {
  const url = new URL(req.url);
  const isLogin = url.pathname === '/auth/login';
  
  const result = await checkRateLimit(req, {
    limit: isLogin ? 5 : 100,      // Strict for login
    window: isLogin ? 300 : 60,    // 5 per 5min vs 100 per min
  });

  if (!result.allowed) {
    return new Response('Too many requests', { status: 429 });
  }

  return new Response('OK');
}

export const config = { runtime: 'edge' };
```

### Return Rate Limit Headers

```typescript
import { checkRateLimit } from '@limity/edge';

export default async function handler(req: Request) {
  const result = await checkRateLimit(req, {
    limit: 100,
    window: 60,
  });

  const headers = new Headers({
    'X-RateLimit-Limit': '100',
    'X-RateLimit-Remaining': result.remaining.toString(),
    'X-RateLimit-Reset': result.reset.toString(),
  });

  if (!result.allowed) {
    return new Response('Too many requests', {
      status: 429,
      headers,
    });
  }

  return new Response('Hello!', { headers });
}

export const config = { runtime: 'edge' };
```

## Supported Platforms

✅ **Vercel Edge Functions** - Automatic IP detection from `x-forwarded-for`
✅ **Cloudflare Workers** - Use `cf-connecting-ip` header
✅ **Deno** - Works with Deno Deploy
✅ **AWS Lambda@Edge** - Compatible with Lambda@Edge
✅ **Netlify Edge Functions** - Works with Netlify
✅ **Any Fetch API runtime** - Standard Request/Response API

## Configuration

### Environment Variables

- **`RATE_LIMIT_API_KEY`** (optional) - API key for hosted mode. If not set, uses memory mode.

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

## Memory vs Hosted Mode

### Memory Mode (Default)

Stores rate limits in the edge function runtime. Perfect for:
- Most edge function use cases
- Single-region deployments
- Development and testing

### Hosted Mode (With API Key)

Uses a remote hosted service. Perfect for:
- Multi-region deployments
- Global rate limiting
- Consistent limits across all edge locations

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

## When to Use

✅ Use `@limity/edge` when:
- You're deploying to edge functions (Vercel, Cloudflare, etc.)
- You need rate limiting in serverless environments
- You want zero cold start overhead
- You need sub-millisecond latency

🔗 Use `@limity/core` instead if:
- You're building a library or SDK
- You need maximum control over rate limiting
- You're not on edge/serverless

📦 Use `@limity/node` instead if:
- You're building an Express.js application
- You need a traditional Node.js server

## Performance

- **Memory mode**: ~0.1-0.5ms per check (in-process)
- **Hosted mode**: ~50-200ms per check (network + API)
- **No cold start penalty** - Lightweight import
