# Limity Express Example

Simple Express server demonstrating Limity rate limiting.

## Setup

1. Install dependencies:
   ```bash
   pnpm install
   ```

2. Start the server:
   ```bash
   pnpm dev
   ```

Server runs on `http://localhost:3000`

## Testing

### Single request
```bash
curl http://localhost:3000/api/data
```

Response:
```json
{
  "message": "Hello from Limity!",
  "data": {
    "timestamp": "2024-04-27T12:34:56.789Z",
    "requestCount": 1
  },
  "rateLimit": {
    "limit": 10,
    "remaining": 9,
    "reset": 1714239296,
    "resetAt": "2024-04-27T12:34:56.789Z"
  }
}
```

### Rate limit exceeded
```bash
# Make 11 requests in quick succession
for i in {1..11}; do curl http://localhost:3000/api/data; done
```

11th request returns 429:
```json
{
  "error": "Too many requests",
  "retryAfter": 5
}
```

### Health check (not rate limited)
```bash
curl http://localhost:3000/health
```

## Configuration

Edit `src/index.ts` to customize:

- `limit` - Requests allowed per window
- `window` - Window duration in seconds
- `keyFn` - Function to extract rate limit key (default: IP address)
- `skip` - Function to skip rate limiting for certain routes

## Production Usage

Set `RATE_LIMIT_API_KEY` to use hosted rate limiting:

```bash
export RATE_LIMIT_API_KEY=your_api_key
pnpm dev
```

Without it, uses in-memory limiter (suitable for single-instance deployments).

## Endpoints

- `GET /health` - Health check (always allowed)
- `GET /api/data` - Fetch data
- `POST /api/data` - Create data
- `GET /api/limited` - Limited endpoint

All routes except `/health` are rate limited.
