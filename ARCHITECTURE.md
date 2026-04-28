# Architecture

## Overview

Limity is a monorepo containing a rate limiting system with:

1. **Hosted backend** - Go API with Redis (optional, via Upstash)
2. **SDK packages** - TypeScript libraries for different environments
3. **Example apps** - Demonstrations of usage

## Design Principles

- **Minimal by default** - Works without any external dependencies
- **Auto-upgrade** - Set an API key to instantly scale to hosted mode
- **Consistent API** - Same behavior everywhere (Node, Edge, Go)
- **Fail gracefully** - Fallback to memory limiter if API fails
- **Production ready** - Clean code, proper error handling

## Architecture Diagram

```
┌─────────────────────────────────────────────┐
│          Limity Monorepo                    │
├─────────────────────────────────────────────┤
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │  packages/core                       │  │
│  │  ├─ memoryLimiter()                  │  │
│  │  ├─ hostedLimiter()                  │  │
│  │  └─ rateLimit() [orchestrator]       │  │
│  └──────────────────────────────────────┘  │
│                  ▲                          │
│          ┌───────┼───────┐                 │
│          │       │       │                 │
│  ┌──────┴──┐ ┌──┴───┐ ┌──┴───────┐       │
│  │packages │ │examples│ │apps/api  │       │
│  │/node    │ │/express│ │(Go)      │       │
│  │         │ │        │ │          │       │
│  │Express  │ │Demo    │ │Upstash   │       │
│  │Middleware│ │Server  │ │Redis API │       │
│  └─────────┘ └────────┘ └──────────┘       │
│                                             │
│  ┌──────────────────────────────────────┐  │
│  │ packages/edge                        │  │
│  │ Fetch/Edge helper (Vercel, CF)       │  │
│  └──────────────────────────────────────┘  │
│                                             │
└─────────────────────────────────────────────┘
```

## Packages

### `packages/core`

Core rate limiting logic shared across all packages.

**Exports:**
- `rateLimit(options)` - Main function
- `memoryLimiter(options)` - In-memory rate limiter
- `hostedLimiter(options, apiKey)` - Hosted API limiter

**Behavior:**
1. If `RATE_LIMIT_API_KEY` set → use hosted API with fallback
2. Otherwise → use memory limiter

### `packages/node`

Express middleware for Node.js servers.

**Exports:**
- `rateLimit(config)` - Middleware factory

**Features:**
- Extracts key from request (IP by default)
- Sets rate limit headers
- Custom error handlers
- Skip function support

### `packages/edge`

Fetch API helper for edge/serverless (Vercel, Cloudflare Workers).

**Exports:**
- `checkRateLimit(request, options)` - Main function

**Features:**
- Works with Fetch API (Request/Response)
- Auto-detects IP from common headers
- Suitable for edge runtimes

### `apps/api`

Go backend with Upstash Redis integration.

**Endpoints:**
- `POST /check` - Rate limit check
- `GET /health` - Health check

**Algorithm:**
Fixed window counter with Redis:
1. Calculate window start: `now - (now % window)`
2. Redis key: `ratelimit:{key}:{windowStart}`
3. INCR to increment counter
4. EXPIRE on first increment
5. Check if count > limit

**Redis:** Uses Upstash REST API (no client library needed)

### `examples/express-app`

Simple demonstration of Limity in Express.

**Features:**
- Basic endpoints with rate limiting
- Custom rate limit per route
- Health check (not limited)
- Error handling
- Response headers

## Data Flow

### In-Memory Mode

```
Request
  ↓
rateLimit(options)
  ↓
Check RATE_LIMIT_API_KEY
  ↓ (not set)
memoryLimiter()
  ↓
Calculate window
  ↓
Check memory store
  ↓
Increment counter
  ↓
Return result
```

### Hosted Mode

```
Request
  ↓
rateLimit(options)
  ↓
Check RATE_LIMIT_API_KEY
  ↓ (set)
hostedLimiter()
  ↓
POST /check to API
  ↓
[Optional fallback to memory on error]
  ↓
Return result
```

## Window Algorithm

Fixed-window counter:

```
Current time: 1714396842
Window: 60 seconds

windowStart = 1714396842 - (1714396842 % 60)
           = 1714396842 - 42
           = 1714396800

windowEnd = windowStart + window
          = 1714396800 + 60
          = 1714396860

Redis key: ratelimit:user:123:1714396800
```

At second 1714396860, counter resets (new window).

## Response Format

All packages return the same shape:

```typescript
{
  allowed: boolean;    // Whether request is allowed
  remaining: number;   // Requests left in window
  reset: number;       // Unix timestamp when window resets
}
```

## Error Handling

**Memory limiter:** Always succeeds (uses Map)

**Hosted limiter:** 
- On API error → return `allowed: false` with default reset time
- On network error → fallback to memory limiter (no throw)

**Middleware:**
- On rate limit exceeded → return 429 with JSON response
- On error → allow request to pass (fail open)

## Configuration

### Environment Variables

**All packages:**
- `RATE_LIMIT_API_KEY` - API key for hosted mode

**Backend (apps/api):**
- `UPSTASH_REDIS_REST_URL` - Redis endpoint
- `UPSTASH_REDIS_REST_TOKEN` - Redis token
- `PORT` - Server port (default: 8080)

### Defaults

- **Limit:** 100 requests
- **Window:** 60 seconds
- **Key (Node):** Request IP address
- **Key (Edge):** Forwarded IP or unknown

## Performance

**Memory limiter:**
- Latency: ~1ms (Map lookup + increment)
- No external dependencies
- Suitable for single-instance deployments

**Hosted mode:**
- Latency: ~100-200ms (network + Redis)
- Scales across multiple instances
- Requires API key + Upstash account

## Security

- API requests use Bearer token authentication
- Keys are function of user/IP + window timestamp
- No sensitive data stored in memory
- Rate limits prevent brute force attacks

## Testing

See individual package READMEs for testing instructions.

Quick test:
```bash
# Terminal 1: Start Go backend
cd apps/api
go run main.go

# Terminal 2: Run example app
cd examples/express-app
pnpm dev

# Terminal 3: Test
curl http://localhost:3000/api/data
```

## Deployment

### Backend (Go)

```bash
# Build
go build -o api main.go

# Run with environment variables
PORT=8080 UPSTASH_REDIS_REST_URL=... UPSTASH_REDIS_REST_TOKEN=... ./api
```

### Node.js (Express)

```bash
# Install
pnpm install

# Run
RATE_LIMIT_API_KEY=... pnpm start
```

### Edge/Serverless

Deploy to Vercel, Cloudflare, etc. The packages are edge-ready (no Node.js dependencies).

## Future Enhancements

- [ ] Dashboard for monitoring
- [ ] Advanced algorithms (sliding window, token bucket)
- [ ] Per-endpoint custom limits
- [ ] Rate limit sharing across distributed systems
- [ ] Metrics/analytics
- [ ] Admin API for key management
