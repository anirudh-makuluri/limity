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
- `apps/api` - Go backend with Upstash Redis

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

Both support **memory mode** (default) and **hosted mode** (with API key).

Default: 100 requests per 60 seconds per IP.

## 🔌 Environment Variables

- `RATE_LIMIT_API_KEY` - Optional. Enables hosted rate limiting via Upstash
- `UPSTASH_REDIS_REST_URL` - Redis endpoint (backend only)
- `UPSTASH_REDIS_REST_TOKEN` - Redis token (backend only)

## 📊 Response Format

```typescript
{
  allowed: boolean;      // Whether the request is allowed
  remaining: number;     // Requests remaining in window
  reset: number;         // Unix timestamp when window resets
}
```

## 🎯 Using with Other Frameworks

Limity works with **any** web framework. See our guide:

- **[USING.md](./USING.md)** - Framework-specific examples (Express, Fastify, FastAPI, Flask, Django, etc.)

All use the same API - no framework-specific code needed!


│   ├── fastify-app/          # Fastify example
│   ├── fastapi-app/          # FastAPI example
│   ├── flask-app/            # Flask example
│   └── django-app/           # Django example
├── pnpm-workspace.yaml
└── package.json
```

## 📚 Documentation

**Getting Started:**
- **[USING.md](./USING.md)** - How to use Limity in your projects (npm/pip/local)

**Framework Guides:**
- **[FRAMEWORK_GUIDE.md](./FRAMEWORK_GUIDE.md)** - TypeScript/JavaScript (Express, Fastify, NestJS, etc.)
- **[PYTHON_GUIDE.md](./PYTHON_GUIDE.md)** - Python (FastAPI, Flask, Django, etc.)
- **[PYTHON_QUICK_START.md](./PYTHON_QUICK_START.md)** - Python quick reference

**Publishing & Development:**
- **[PUBLISHING.md](./PUBLISHING.md)** - How to publish to npm and PyPI
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - System design and how it works
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Development guide
- **[PROJECT_COMPLETE.md](./PROJECT_COMPLETE.md)** - Full project inventory

## 📝 License

MIT
