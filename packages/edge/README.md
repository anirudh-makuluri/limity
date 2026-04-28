# Limity Edge

Rate limiting for edge/serverless environments (Vercel, Cloudflare, etc).

## Usage

### Vercel Edge Functions

```typescript
import { checkRateLimit } from '@limity/edge';

export default async function handler(req: Request) {
  const result = await checkRateLimit(req, {
    limit: 50,
    window: 60,
  });

  if (!result.allowed) {
    return new Response('Too many requests', { status: 429 });
  }

  return new Response('Hello!');
}
```

### Cloudflare Workers

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

## Custom Key Function

```typescript
const result = await checkRateLimit(request, {
  keyFn: (req) => {
    if (req instanceof Request) {
      // Prefer Cloudflare header
      return req.headers.get('cf-connecting-ip') || 'unknown';
    }
    return req.ip || 'unknown';
  },
});
```

## Features

- Works with Fetch API (Request/Response)
- Auto-detects IP from common headers
- Same fallback behavior as core package
- Suitable for edge runtimes (no Node.js dependencies)

## Environment

Set `RATE_LIMIT_API_KEY` for hosted rate limiting:

```bash
RATE_LIMIT_API_KEY=your_api_key
```
