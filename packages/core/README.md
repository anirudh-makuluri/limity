# Limity Core

Core rate limiting logic shared across all Limity packages.

## Usage

```typescript
import { rateLimit } from '@limity/core';

const result = await rateLimit({
  key: 'user:123',
  limit: 50,      // optional, default 100
  window: 60,     // optional, default 60
});

if (!result.allowed) {
  console.log('Rate limited');
}

console.log(`Remaining: ${result.remaining}`);
console.log(`Reset at: ${new Date(result.reset * 1000)}`);
```

## Behavior

- **Without API key**: Uses in-memory limiter (fast, zero setup)
- **With API key**: Uses hosted API with automatic fallback to memory on failure
- **Fixed window**: Resets every `window` seconds

## Environment

Set `RATE_LIMIT_API_KEY` to enable hosted mode:

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

## Defaults

- Limit: 100 requests
- Window: 60 seconds
