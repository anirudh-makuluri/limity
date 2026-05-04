# Limity Fastify Example

Simple Fastify server demonstrating Limity rate limiting.

## Features

- **Global hook** - Rate limits all requests
- **Per-endpoint limits** - Different limits for different routes
- **User-based limiting** - Rate limit by user ID or IP
- **Standard headers** - X-RateLimit headers
- **TypeScript** - Fully typed

## Setup

```bash
pnpm install
pnpm dev
```

Server runs on `http://localhost:3000`

## Endpoints

### GET /health
- Always allowed (not rate limited)

```bash
curl http://localhost:3000/health
```

### GET /api/data
- Limit: 100 requests per 60 seconds
- Key: IP address

```bash
curl http://localhost:3000/api/data
```

### GET /api/limited
- Limit: 10 requests per 60 seconds (strict!)
- Key: IP address

```bash
curl http://localhost:3000/api/limited
```

### POST /api/create
- Limit: 20 requests per 60 seconds
- Key: User ID if provided, otherwise IP

```bash
curl -X POST http://localhost:3000/api/create \
  -H "Content-Type: application/json" \
  -d '{"userId":"user123"}'
```

## Testing Rate Limiting

### Flood requests
```bash
# Succeeds 100 times, then fails
for i in {1..105}; do 
  curl http://localhost:3000/api/data
  echo "Request $i"
done
```

### Check headers
```bash
curl -i http://localhost:3000/api/data
```

Look for:
- `X-RateLimit-Limit` - Max requests per window
- `X-RateLimit-Remaining` - Requests left in window
- `X-RateLimit-Reset` - Unix timestamp when window resets
- `Retry-After` - Seconds to wait before retrying (on 429)

## Architecture

```
Request
  ↓
Global hook (onRequest)
  ↓
Check global rate limit (1000 req/min)
  ↓
Route handler
  ↓
Check endpoint-specific rate limit
  ↓
Return 429 if limited, otherwise process request
```

## How It Works

### 1. Global Hook

```typescript
fastify.addHook('onRequest', async (request, reply) => {
  const ip = request.ip;
  
  const result = await rateLimit({
    key: `global:${ip}`,
    limit: 1000,
    window: 60,
  });

  if (!result.allowed) {
    reply.code(429).send({ error: 'Too many requests' });
    return;
  }

  request.rateLimit = result;
});
```

### 2. Route Handler

```typescript
fastify.get('/api/data', async (request, reply) => {
  const result = await rateLimit({
    key: `data:${ip}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    reply.code(429).send({ error: 'Too many requests' });
    return;
  }

  reply.send({ data: 'hello', remaining: result.remaining });
});
```

## Customization

### Change limits

Edit `src/index.ts`:

```typescript
// Global limit
limit: 1000,

// Per-endpoint limit
limit: 100,
```

### Rate limit by user ID

```typescript
let key = ip;
const userId = request.headers['x-user-id'];
if (userId) {
  key = `user:${userId}`;
}

const result = await rateLimit({ key, limit, window });
```

### Skip rate limiting

```typescript
fastify.addHook('onRequest', async (request, reply) => {
  if (request.url.startsWith('/webhook')) {
    return; // Skip
  }

  // Check rate limit
});
```

## Environment

### Optional: Use Hosted Rate Limiting

```bash
export LIMITY_API_KEY=your_api_key
pnpm dev
```

Without it, uses fast in-memory limiter.

### Port Configuration

```bash
export PORT=8000
pnpm dev
```

## Performance

- **Memory mode:** ~1ms per request
- **Hosted mode:** ~100-200ms per request

Memory mode is suitable for single-instance deployments.

## Files

```
src/
└── index.ts         # Fastify server with rate limiting
package.json
tsconfig.json
```

## Next Steps

1. Run the example: `pnpm dev`
2. Test endpoints: `curl http://localhost:3000/api/data`
3. Trigger rate limits: spam requests
4. Check headers: `curl -i http://localhost:3000/api/data`
5. Deploy and scale with `LIMITY_API_KEY`

---

Simple, fast, and production-ready! 🚀
