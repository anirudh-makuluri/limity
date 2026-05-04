# Limity

Developer-first rate limiting tool. Built for simplicity and performance.

**Works everywhere:** TypeScript/JavaScript, Python, Go.

## 🎯 Features

- **In-memory by default** - Zero setup, works immediately
- **Auto-upgrade to hosted** - Set an API key, instantly scales
- **Consistent API** - Same logic everywhere (TypeScript, Python, Go)
- **Minimal overhead** - ~1ms latency with memory limiter
- **Production-ready** - Clean code, proper error handling
- **Framework agnostic** - Use with any framework

## 📦 Packages

### TypeScript/JavaScript

**Choose the package that fits your use case:**

| Package | Use Case | Best For |
|---------|----------|----------|
| **[@limity/core](./packages/core)** | Direct rate limiting logic | Libraries, SDKs, custom implementations |
| **[@limity/node](./packages/node)** | Express middleware | Express.js applications |
| **[@limity/edge](./packages/edge)** | Edge functions & serverless | Vercel, Cloudflare, Deno Deploy |

**Examples:**
- `examples/express-app` - Express.js usage
- `examples/nextjs-app` - Next.js usage
- `examples/fastify-app` - Fastify usage

### Python
- `packages/python` - Pure Python SDK (zero dependencies)
- `examples/fastapi-app` - FastAPI example
- `examples/flask-app` - Flask example
- `examples/django-app` - Django example

### Go
- `apps/backend` - Go API backend

### Apps
- `apps/dashboard` - React dashboard for auth + API key management

## 🚀 Quick Start

### TypeScript/JavaScript

**For basic rate limiting:**

```typescript
import { rateLimit } from '@limity/core';

const result = await rateLimit({
  key: 'user:123',
  limit: 100,
  window: 60,
});

if (!result.allowed) {
  return error(429, 'Too many requests');
}
```

`window` is in **seconds** (for example, `window: 60` = 1 minute).

**For Express.js (automatic IP rate limiting):**

```typescript
import express from 'express';
import { rateLimit } from '@limity/node';

const app = express();
app.use(rateLimit()); // Automatic IP-based rate limiting

app.listen(3000);
```

**For edge functions (Vercel, Cloudflare, etc):**

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

  return new Response('Hello!');
}
```

### Python

```python
from limity import rate_limit

result = rate_limit(
  key='user:123',
  limit=100,
  window=60,
)

if not result.allowed:
  return error(429, 'Too many requests')
```

`window` is in **seconds** here as well.

Both support **memory mode** (default) and **hosted mode** (with API key).

Default: 100 requests per 60 seconds per IP.
All `window` values across SDKs are in **seconds**.

## 🔌 Environment Variables

- `LIMITY_API_KEY` - Optional. Enables hosted rate limiting

## 📊 Response Format

```typescript
{
  allowed: boolean;      // Whether the request is allowed
  remaining: number;     // Requests remaining in window
  reset: number;         // Unix timestamp when window resets
}
```

## 📚 Documentation

- **[USING.md](./USING.md)** - How to use Limity in your projects
- **[PUBLISHING.md](./PUBLISHING.md)** - How to publish packages

## 📝 License

MIT
