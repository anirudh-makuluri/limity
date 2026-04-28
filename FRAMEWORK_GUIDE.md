# Using Limity in Different Frameworks

Limity's core is framework-agnostic. Here are examples for popular Node.js frameworks.

## Quick Reference

All examples follow this pattern:

```typescript
import { rateLimit } from '@limity/core';

// Extract key (IP, user ID, etc)
const key = getKey(request);

// Check rate limit
const result = await rateLimit({
  key,
  limit: 100,
  window: 60,
});

// Respond based on result
if (!result.allowed) {
  return respond(429, 'Too many requests');
}

return respond(200, data);
```

---

## Express

Use the built-in middleware from `@limity/node`:

```typescript
import express from 'express';
import { rateLimit } from '@limity/node';

const app = express();

// Simple: use defaults
app.use(rateLimit());

// Or customize
app.use(rateLimit({
  limit: 50,
  window: 60,
  keyFn: (req) => req.user?.id || req.ip,
  skip: (req) => req.path === '/health',
}));

app.listen(3000);
```

**File:** `examples/express-app/`

---

## Next.js

Use `@limity/core` in API routes and middleware:

```typescript
// app/api/data/route.ts
import { rateLimit } from '@limity/core';
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  const ip = request.headers.get('x-forwarded-for') || 'unknown';
  
  const result = await rateLimit({
    key: `api:${ip}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    return NextResponse.json(
      { error: 'Too many requests' },
      { status: 429 }
    );
  }

  return NextResponse.json({ data: 'hello' });
}
```

**File:** `examples/nextjs-app/`

---

## Fastify

```typescript
import Fastify from 'fastify';
import { rateLimit } from '@limity/core';

const fastify = Fastify();

// Global hook
fastify.addHook('onRequest', async (request, reply) => {
  const result = await rateLimit({
    key: request.ip,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    reply.status(429).send({ error: 'Too many requests' });
  }

  // Attach to request for downstream use
  (request as any).rateLimit = result;
});

fastify.get('/data', async (request, reply) => {
  const { remaining, reset } = (request as any).rateLimit;
  return { data: 'hello', remaining, reset };
});

fastify.listen({ port: 3000 });
```

---

## Hapi

```typescript
import Hapi from '@hapi/hapi';
import { rateLimit } from '@limity/core';

const server = Hapi.server({
  port: 3000,
  host: 'localhost',
});

// Register as plugin
await server.register({
  plugin: {
    name: 'limity',
    version: '1.0.0',
    register: async (server) => {
      server.ext('onRequest', async (request, h) => {
        const result = await rateLimit({
          key: request.info.remoteAddress,
          limit: 100,
          window: 60,
        });

        if (!result.allowed) {
          return h.response({ error: 'Too many requests' }).code(429);
        }

        request.app.rateLimit = result;
        return h.continue;
      });
    },
  },
});

server.route({
  method: 'GET',
  path: '/data',
  handler: (request, h) => {
    const { remaining } = request.app.rateLimit;
    return { data: 'hello', remaining };
  },
});

await server.start();
```

---

## Koa

```typescript
import Koa from 'koa';
import { rateLimit } from '@limity/core';

const app = new Koa();

// Middleware
app.use(async (ctx, next) => {
  const result = await rateLimit({
    key: ctx.request.ip,
    limit: 100,
    window: 60,
  });

  ctx.set('X-RateLimit-Limit', '100');
  ctx.set('X-RateLimit-Remaining', result.remaining.toString());
  ctx.set('X-RateLimit-Reset', result.reset.toString());

  if (!result.allowed) {
    ctx.status = 429;
    ctx.body = { error: 'Too many requests' };
    return;
  }

  ctx.state.rateLimit = result;
  await next();
});

app.use(async (ctx) => {
  const { remaining } = ctx.state.rateLimit;
  ctx.body = { data: 'hello', remaining };
});

app.listen(3000);
```

---

## NestJS

```typescript
import { NestFactory } from '@nestjs/core';
import { Injectable, NestMiddleware } from '@nestjs/common';
import { rateLimit } from '@limity/core';

@Injectable()
export class RateLimitMiddleware implements NestMiddleware {
  async use(req: any, res: any, next: () => void) {
    const result = await rateLimit({
      key: req.ip,
      limit: 100,
      window: 60,
    });

    if (!result.allowed) {
      res.status(429).json({ error: 'Too many requests' });
      return;
    }

    req.rateLimit = result;
    next();
  }
}

// In your module
@Module({
  controllers: [DataController],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(RateLimitMiddleware).forRoutes('*');
  }
}

@Controller('api')
export class DataController {
  @Get('data')
  getData(@Req() req: any) {
    const { remaining } = req.rateLimit;
    return { data: 'hello', remaining };
  }
}

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  await app.listen(3000);
}

bootstrap();
```

---

## Remix

```typescript
// app/routes/api/data.tsx
import { json, type LoaderFunction } from '@remix-run/node';
import { rateLimit } from '@limity/core';

export const loader: LoaderFunction = async ({ request }) => {
  const clientIp = request.headers.get('x-forwarded-for') || 'unknown';

  const result = await rateLimit({
    key: `api:${clientIp}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    return json(
      { error: 'Too many requests' },
      { status: 429 }
    );
  }

  return json({
    data: 'hello',
    rateLimit: {
      remaining: result.remaining,
      reset: result.reset,
    },
  });
};
```

---

## AstroJS

```typescript
// src/pages/api/data.ts
import type { APIRoute } from 'astro';
import { rateLimit } from '@limity/core';

export const GET: APIRoute = async ({ request }) => {
  const clientIp = request.headers.get('x-forwarded-for') || 'unknown';

  const result = await rateLimit({
    key: `api:${clientIp}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    return new Response(
      JSON.stringify({ error: 'Too many requests' }),
      { status: 429 }
    );
  }

  return new Response(
    JSON.stringify({
      data: 'hello',
      remaining: result.remaining,
    })
  );
};
```

---

## SvelteKit

```typescript
// src/routes/api/data/+server.ts
import type { RequestHandler } from '@sveltejs/kit';
import { rateLimit } from '@limity/core';

export const GET: RequestHandler = async ({ request }) => {
  const clientIp = request.headers.get('x-forwarded-for') || 'unknown';

  const result = await rateLimit({
    key: `api:${clientIp}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    return new Response(
      JSON.stringify({ error: 'Too many requests' }),
      { status: 429 }
    );
  }

  return new Response(
    JSON.stringify({
      data: 'hello',
      remaining: result.remaining,
    })
  );
};
```

---

## Elysia (Bun)

```typescript
import { Elysia } from 'elysia';
import { rateLimit } from '@limity/core';

new Elysia()
  .derive(async ({ request }) => {
    const result = await rateLimit({
      key: request.headers.get('x-forwarded-for') || 'unknown',
      limit: 100,
      window: 60,
    });

    return { rateLimit: result };
  })
  .get('/api/data', ({ rateLimit }) => {
    if (!rateLimit.allowed) {
      return { error: 'Too many requests' };
    }

    return { data: 'hello', remaining: rateLimit.remaining };
  })
  .listen(3000);
```

---

## Deno (Fresh)

```typescript
// routes/api/data.ts
import { rateLimit } from '@limity/core';
import type { Handlers } from '$fresh/server.ts';

export const handler: Handlers = {
  async GET(req: Request) {
    const url = new URL(req.url);
    const clientIp = req.headers.get('x-forwarded-for') || 'unknown';

    const result = await rateLimit({
      key: `api:${clientIp}`,
      limit: 100,
      window: 60,
    });

    if (!result.allowed) {
      return new Response(
        JSON.stringify({ error: 'Too many requests' }),
        { status: 429 }
      );
    }

    return new Response(
      JSON.stringify({
        data: 'hello',
        remaining: result.remaining,
      })
    );
  },
};
```

---

## Common Patterns

### Extract Client IP

```typescript
// Express
req.ip

// Next.js / Fetch
request.headers.get('x-forwarded-for') || 
request.headers.get('x-real-ip')

// Fastify
request.ip

// Koa
ctx.request.ip

// NestJS
req.ip
```

### Set Response Headers

```typescript
// Standard rate limit headers
response.set('X-RateLimit-Limit', limit.toString());
response.set('X-RateLimit-Remaining', result.remaining.toString());
response.set('X-RateLimit-Reset', result.reset.toString());
response.set('Retry-After', retryAfter.toString());
```

### Rate Limit by User ID

```typescript
// Instead of IP, use user ID
const key = user?.id || request.ip;

const result = await rateLimit({
  key: `user:${key}`,
  limit: 100,
  window: 60,
});
```

### Different Limits Per Endpoint

```typescript
// Public endpoint - loose limit
const result1 = await rateLimit({
  key: `public:${ip}`,
  limit: 1000,
  window: 60,
});

// Authenticated endpoint - tighter limit
const result2 = await rateLimit({
  key: `premium:${userId}`,
  limit: 50,
  window: 60,
});

// Admin endpoint - very tight limit
const result3 = await rateLimit({
  key: `admin:${userId}`,
  limit: 10,
  window: 60,
});
```

### Skip Rate Limiting

```typescript
// For health checks, webhooks, etc.
const skipPaths = ['/health', '/webhook', '/.well-known'];

if (skipPaths.includes(path)) {
  return next();
}

const result = await rateLimit({ key, limit, window });
```

---

## Environment Variables

All frameworks:

```bash
# Enable hosted rate limiting (optional)
export RATE_LIMIT_API_KEY=your_api_key

# Backend (if running Go API)
export UPSTASH_REDIS_REST_URL=https://...
export UPSTASH_REDIS_REST_TOKEN=...
```

Without `RATE_LIMIT_API_KEY`, Limity uses the fast in-memory limiter.

---

## Testing

Every framework above can be tested the same way:

```bash
# Single request
curl http://localhost:3000/api/data

# Flood (get rate limited)
for i in {1..105}; do curl http://localhost:3000/api/data; done

# Check headers
curl -i http://localhost:3000/api/data
```

---

## Summary

✅ Works with any Node.js framework  
✅ Works with Deno, Bun, edge runtimes  
✅ Same API everywhere  
✅ Auto-switches from memory to hosted mode  
✅ Zero breaking changes when scaling  

Choose your framework above, copy the pattern, and you're rate-limiting! 🚀
