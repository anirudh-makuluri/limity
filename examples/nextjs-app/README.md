# Limity Next.js Example

Full-featured Next.js example using Limity for rate limiting.

## Features

- **App Router** - Uses Next.js 14+ App Router
- **API Routes** - Multiple endpoints with different limits
- **Middleware** - Global rate limiting middleware
- **Interactive UI** - Test rate limiting with buttons
- **Rate Limit Headers** - Standard HTTP rate limit headers
- **TypeScript** - Fully typed

## Setup

```bash
# Install dependencies
pnpm install

# Copy environment (optional, for hosted mode)
cp .env.example .env.local

# Start development server
pnpm dev
```

Server runs on `http://localhost:3000`

## Endpoints

### GET /api/data
- **Limit:** 100 requests per 60 seconds
- **Key:** IP address
- **Response:** JSON with rate limit info

```bash
curl http://localhost:3000/api/data
```

### GET /api/limited
- **Limit:** 10 requests per 60 seconds (stricter!)
- **Key:** IP address
- **Response:** JSON with strict limit warning

```bash
curl http://localhost:3000/api/limited
```

### POST /api/health
- **Limit:** 20 requests per 60 seconds
- **Key:** User ID if provided, otherwise IP
- **Body:** `{ "userId": "user123" }` (optional)
- **Response:** Created resource with ID

```bash
curl -X POST http://localhost:3000/api/health \
  -H "Content-Type: application/json" \
  -d '{"userId":"user123"}'
```

## Testing Rate Limiting

### Via UI
1. Open `http://localhost:3000`
2. Click buttons to test endpoints
3. Watch the rate limit info in responses

### Via curl - Flood requests
```bash
# GET endpoint - succeeds 100 times, then fails
for i in {1..105}; do 
  curl http://localhost:3000/api/data
  echo "Request $i"
done

# Limited endpoint - succeeds 10 times, then fails
for i in {1..15}; do 
  curl http://localhost:3000/api/limited
  echo "Request $i"
done
```

### Check response headers
```bash
curl -i http://localhost:3000/api/data
```

Look for:
- `X-RateLimit-Limit` - Max requests per window
- `X-RateLimit-Remaining` - Requests left in window
- `X-RateLimit-Reset` - Unix timestamp when window resets
- `Retry-After` - Seconds to wait before retrying (on 429)

## How It Works

### 1. Extract IP
```typescript
const ip = request.headers.get('x-forwarded-for') || 
           request.headers.get('x-real-ip') || 
           'unknown';
```

### 2. Create unique key per endpoint
```typescript
const result = await rateLimit({
  key: `data:${ip}`,
  limit: 100,
  window: 60,
});
```

### 3. Check result
```typescript
if (!result.allowed) {
  return NextResponse.json(
    { error: 'Too many requests' },
    { status: 429 }
  );
}
```

### 4. Return with headers
```typescript
return NextResponse.json(data, {
  headers: {
    'X-RateLimit-Remaining': result.remaining.toString(),
    'X-RateLimit-Reset': result.reset.toString(),
  },
});
```

## Middleware

`middleware.ts` provides global rate limiting:

- Applies to all routes (except static assets)
- Limit: 1000 requests per 60 seconds
- Falls through to rate-limited endpoints

You can customize or disable it:

```typescript
// In middleware.ts, adjust:
const result = await rateLimit({
  key: `global:${ip}`,
  limit: 1000,  // <-- Change this
  window: 60,   // <-- Or this
});
```

## Per-User Rate Limiting

For authenticated users, rate limit by user ID instead of IP:

```typescript
const userId = request.headers.get('x-user-id');

const result = await rateLimit({
  key: `user:${userId}`,
  limit: 100,
  window: 60,
});
```

## Environment Variables

### Optional: Use Hosted Rate Limiting

```bash
# .env.local
NEXT_PUBLIC_LIMITY_API_KEY=your_api_key
```

Then the app will:
1. Use hosted API if key is set
2. Fall back to memory limiter if API fails
3. Work seamlessly without any changes

## Deployment

### Vercel (Native Support)
```bash
vercel env add LIMITY_API_KEY
vercel deploy
```

The middleware and API routes run on Vercel's Edge Runtime (fast!).

### Docker
```bash
docker build -t limity-next .
docker run -p 3000:3000 limity-next
```

### Self-hosted
```bash
pnpm build
pnpm start
```

## Customization

### Change limits per endpoint

In `app/api/[endpoint]/route.ts`:

```typescript
const result = await rateLimit({
  key: `endpoint:${ip}`,
  limit: 50,    // <-- Adjust
  window: 120,  // <-- Adjust (in seconds)
});
```

### Rate limit by user ID

```typescript
// Extract from JWT, session, headers, etc.
const userId = await getCurrentUserId(request);

const result = await rateLimit({
  key: `api:${userId}`,
  limit: 100,
  window: 60,
});
```

### Skip rate limiting

```typescript
// Skip for certain paths
if (request.nextUrl.pathname.startsWith('/api/public')) {
  return NextResponse.next();
}
```

## Architecture

```
Request
  ↓
middleware.ts (global limit)
  ↓
API route handler
  ↓
rateLimit() from @limity/core
  ↓
Memory limiter (or hosted API if key set)
  ↓
Return 429 if limited, otherwise continue
```

## Performance

- **Memory mode:** ~1ms per request
- **Hosted mode:** ~100-200ms per request
- Limits are per IP/user, not global

## Files

```
app/
├── page.tsx              # Home page with UI
├── layout.tsx            # Root layout
├── globals.css           # Global styles
├── api/
│   ├── data/
│   │   └── route.ts      # GET /api/data (100/min)
│   ├── limited/
│   │   └── route.ts      # GET /api/limited (10/min)
│   └── health/
│       └── route.ts      # POST /api/health (20/min)
middleware.ts             # Global rate limiting
package.json
tsconfig.json
next.config.js
```

## Next Steps

1. Try the interactive UI at `http://localhost:3000`
2. Trigger rate limits by clicking rapidly
3. Check response headers with `curl -i`
4. Add your own endpoints
5. Deploy to Vercel or your host
6. Set `LIMITY_API_KEY` for hosted mode

---

That's it! You now have a fully rate-limited Next.js app. 🚀
