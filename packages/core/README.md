# @limity/core

**Framework-agnostic rate limiting core.** Use this directly if you need rate limiting in any JavaScript/TypeScript application.

## Installation

```bash
npm install @limity/core
# or
pnpm add @limity/core
# or
yarn add @limity/core
```

## Quick Start

```typescript
import { rateLimit } from '@limity/core';

const result = await rateLimit({
  key: 'user:123',
  limit: 100,
  window: 60,
});

if (!result.allowed) {
  throw new Error('Too many requests');
}

console.log(`Requests remaining: ${result.remaining}`);
console.log(`Resets at: ${new Date(result.reset * 1000)}`);
```

## How It Works

### Memory Mode (Default)

Stores rate limit counters in memory. Perfect for:
- Single-server applications
- Development and testing
- Low-traffic applications

```typescript
const result = await rateLimit({
  key: 'user:123',
  limit: 100,
  window: 60,
});
```

### Hosted Mode (With API Key)

Uses a remote hosted service. Perfect for:
- Multi-server deployments
- Serverless/edge functions
- High-traffic applications
- Shared rate limiting across services

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

Then use the same code - it automatically switches to hosted mode:

```typescript
const result = await rateLimit({
  key: 'user:123',
  limit: 100,
  window: 60,
});
// Uses hosted API, falls back to memory on failure
```

## API Reference

### `rateLimit(options)`

Checks if a request should be allowed based on rate limits.

**Parameters:**

```typescript
interface RateLimitOptions {
  key: string;       // Unique identifier (user ID, IP, etc.)
  limit?: number;    // Max requests per window (default: 100)
  window?: number;   // Window size in seconds (default: 60)
}
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

### Rate Limit by User ID

```typescript
const userId = 'user:456';
const result = await rateLimit({
  key: userId,
  limit: 50,
  window: 60,
});

if (!result.allowed) {
  console.log(`Try again in ${result.reset - Math.floor(Date.now() / 1000)}s`);
}
```

### Rate Limit by IP Address

```typescript
const ipAddress = '192.168.1.1';
const result = await rateLimit({
  key: `ip:${ipAddress}`,
  limit: 100,
  window: 60,
});
```

### Different Limits for Different Endpoints

```typescript
// Strict limit for login attempts
const loginResult = await rateLimit({
  key: `login:${email}`,
  limit: 5,
  window: 300, // 5 per 5 minutes
});

// Relaxed limit for API calls
const apiResult = await rateLimit({
  key: `api:${userId}`,
  limit: 1000,
  window: 60,
});
```

## Configuration

### Environment Variables

- **`RATE_LIMIT_API_KEY`** (optional) - API key for hosted mode. If not set, uses memory mode.

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

### Default Values

- **Limit**: 100 requests
- **Window**: 60 seconds
- **Mode**: Memory (no API key needed)

## When to Use

✅ Use `@limity/core` when:
- You need simple rate limiting logic
- You're building a Node.js library or SDK
- You want full control over rate limiting
- You're building for a single environment

🔗 Use `@limity/node` instead if:
- You're building an Express.js application
- You want automatic IP-based rate limiting

🌐 Use `@limity/edge` instead if:
- You're deploying to edge functions (Vercel, Cloudflare, etc.)
- You need rate limiting in serverless environments

## Performance

- **Memory mode**: ~0.1ms per check (in-process)
- **Hosted mode**: ~50-200ms per check (network + API)
