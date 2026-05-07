# Limity

Reliable, developer-first rate limiting for Node, Edge, and Python.

[![CI](https://github.com/anirudh-makuluri/limity/actions/workflows/ci.yml/badge.svg)](https://github.com/anirudh-makuluri/limity/actions/workflows/ci.yml)

Limity is focused on one thing first: predictable behavior in production.

## Reliability Signals

- Versioned releases with changelog and upgrade guidance.
- CI on every push and pull request.
- Core tests include boundary and concurrency scenarios.
- Explicit algorithm and tradeoff documentation.

## Packages

| Package | Purpose | Use When |
| --- | --- | --- |
| `@limity/core` | Rate limiting engine | You want direct control over keys/limits/windows |
| `@limity/node` | Express middleware | You need route-level protection in Node APIs |
| `@limity/edge` | Edge helper | You run on Vercel/Cloudflare/Workers-like runtimes |
| `limity` (Python) | Python SDK | You build Flask/FastAPI/Django services |

## Quick Start

```ts
import { rateLimit } from '@limity/core';

const result = await rateLimit({ key: 'user:123', limit: 100, window: 60 });
if (!result.allowed) {
  // 429
}
```

## Express Example with Retry Header

```ts
import express from 'express';
import { rateLimit } from '@limity/node';

const app = express();

app.use(rateLimit({ limit: 100, window: 60 }));
app.get('/api', (_req, res) => res.json({ ok: true }));

app.listen(3000);
```

`@limity/node` sets:
- `X-RateLimit-Limit`
- `X-RateLimit-Remaining`
- `X-RateLimit-Reset`
- `Retry-After` when blocked

## Algorithm

Current default algorithm is **fixed-window counter**.

See detailed behavior and tradeoffs in [docs/ALGORITHMS.md](./docs/ALGORITHMS.md).

## Production Examples

- [examples/express-app](./examples/express-app)
- [examples/fastify-app](./examples/fastify-app)
- [examples/nextjs-app](./examples/nextjs-app)
- [examples/fastapi-app](./examples/fastapi-app)
- [examples/flask-app](./examples/flask-app)
- [examples/django-app](./examples/django-app)

## Project Docs

- [CHANGELOG.md](./CHANGELOG.md)
- [ROADMAP.md](./ROADMAP.md)
- [CONTRIBUTING.md](./CONTRIBUTING.md)
- [PUBLISHING.md](./PUBLISHING.md)
- [USING.md](./USING.md)

## License

MIT
