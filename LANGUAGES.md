# Limity Across Languages

Limity works the same everywhere - TypeScript/JavaScript, Python, and Go.

## 🎯 Quick Comparison

| Feature | TypeScript | Python | Go |
|---------|-----------|--------|-----|
| **Core** | `@limity/core` | `limity` | Built-in |
| **Setup** | `npm install` | `pip install` | `go run` |
| **Memory mode** | ✅ | ✅ | HTTP endpoint |
| **Hosted mode** | ✅ | ✅ | HTTP endpoint |
| **Frameworks** | 10+ | 8+ | Standalone |

---

## TypeScript/JavaScript

### Install

```bash
npm install @limity/core
# or
pnpm add @limity/core
```

### Basic Usage

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
```

### Express Middleware

```typescript
import { rateLimit } from '@limity/node';

app.use(rateLimit());
```

### Next.js API Route

```typescript
import { rateLimit } from '@limity/core';
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  const result = await rateLimit({
    key: request.headers.get('x-forwarded-for'),
    limit: 100,
    window: 60,
  });
  
  if (!result.allowed) {
    return NextResponse.json({}, { status: 429 });
  }
  
  return NextResponse.json({ data: 'hello' });
}
```

### Edge Functions

```typescript
import { checkRateLimit } from '@limity/edge';

export default {
  async fetch(request: Request) {
    const result = await checkRateLimit(request);
    
    if (!result.allowed) {
      return new Response('429', { status: 429 });
    }
    
    return new Response('ok');
  },
};
```

### Supported Frameworks

- **Express** - `@limity/node` middleware
- **Next.js** - Direct API routes
- **Fastify** - Via hooks
- **NestJS** - Via custom middleware
- **Koa** - Via middleware
- **Hapi** - Via plugins
- **Remix** - Via loaders
- **SvelteKit** - Via hooks
- **AstroJS** - Via API routes
- **Elysia** - Via derive

See **[FRAMEWORK_GUIDE.md](./FRAMEWORK_GUIDE.md)** for detailed examples.

---

## Python

### Install

```bash
pip install limity
```

### Basic Usage

```python
from limity import rate_limit

result = rate_limit(
    key='user:123',
    limit=100,
    window=60,
)

if not result.allowed:
    raise Exception('Too many requests')
```

### FastAPI

```python
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from limity import rate_limit

app = FastAPI()

@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    result = rate_limit(f"ip:{request.client.host}")
    
    if not result.allowed:
        return JSONResponse(
            {"error": "Too many requests"},
            status_code=429,
        )
    
    return await call_next(request)
```

### Flask

```python
from flask import Flask, jsonify, request
from limity import rate_limit

app = Flask(__name__)

@app.before_request
def check_rate_limit():
    result = rate_limit(f"ip:{request.remote_addr}")
    
    if not result.allowed:
        return jsonify({"error": "Too many requests"}), 429
```

### Django

```python
# middleware.py
from django.http import JsonResponse
from limity import rate_limit

class RateLimitMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response
    
    def __call__(self, request):
        result = rate_limit(f"ip:{request.META.get('REMOTE_ADDR')}")
        
        if not result.allowed:
            return JsonResponse(
                {"error": "Too many requests"},
                status=429,
            )
        
        return self.get_response(request)
```

### Supported Frameworks

- **FastAPI** - Via middleware
- **Flask** - Via `@app.before_request`
- **Django** - Via middleware
- **Starlette** - Via middleware
- **Quart** - Via `@app.before_request`
- **Tornado** - Via `prepare()`
- **Bottle** - Via `@app.hook`
- **CherryPy** - Via plugins

See **[PYTHON_GUIDE.md](./PYTHON_GUIDE.md)** for detailed examples.

---

## Go

### Backend

The Go backend is a standalone HTTP service that handles rate limiting for all languages.

### Setup

```bash
cd apps/api
go run main.go
```

Runs on `http://localhost:8080`

### Usage

Call the API from any language:

```bash
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{
    "key": "user:123",
    "limit": 100,
    "window": 60
  }'
```

Response:

```json
{
  "allowed": true,
  "remaining": 99,
  "reset": 1714396860
}
```

### Environment Variables

```bash
PORT=8080
UPSTASH_REDIS_REST_URL=https://...
UPSTASH_REDIS_REST_TOKEN=...
```

See **[apps/api/README.md](./apps/api/README.md)** for details.

---

## Response Format (All Languages)

All languages return the same structure:

### TypeScript

```typescript
interface RateLimitResult {
  allowed: boolean;      // Request allowed?
  remaining: number;     // Requests left in window
  reset: number;         // Unix timestamp of reset
}
```

### Python

```python
class RateLimitResult:
    allowed: bool       # Request allowed?
    remaining: int      # Requests left in window
    reset: int          # Unix timestamp of reset
    
    def to_dict(self):  # Convert to JSON
        return {...}
```

### Go (from API)

```go
type CheckResponse struct {
    Allowed   bool   `json:"allowed"`
    Remaining int    `json:"remaining"`
    Reset     int64  `json:"reset"`
}
```

---

## Environment Variables (All Languages)

### Enable Hosted Mode

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

Then any language automatically uses hosted API with fallback to memory.

### Backend Configuration (Go)

```bash
export UPSTASH_REDIS_REST_URL=https://...
export UPSTASH_REDIS_REST_TOKEN=...
```

---

## Performance (All Languages)

| Mode | Latency | Scalability |
|------|---------|-------------|
| **Memory** | ~1ms | Single instance |
| **Hosted** | ~100-200ms | Multiple instances |

Choose based on your needs:
- **Development:** Use memory
- **Single instance:** Use memory
- **Distributed:** Use hosted with API key

---

## Common Patterns

### Rate Limit by User ID

**TypeScript:**
```typescript
const result = await rateLimit({
  key: `user:${userId}`,
});
```

**Python:**
```python
result = rate_limit(f"user:{user_id}")
```

**Go (API call):**
```
POST /check
{"key": "user:123"}
```

### Rate Limit by IP

**TypeScript:**
```typescript
const result = await rateLimit({
  key: `ip:${req.ip}`,
});
```

**Python:**
```python
result = rate_limit(f"ip:{request.remote_addr}")
```

### Different Limits Per Endpoint

**TypeScript:**
```typescript
// Public: 1000 req/min
const result = await rateLimit({
  key: `public:${ip}`,
  limit: 1000,
  window: 60,
});
```

**Python:**
```python
# Public: 1000 req/min
result = rate_limit(f"public:{ip}", limit=1000, window=60)
```

### Add Response Headers

**TypeScript:**
```typescript
res.set('X-RateLimit-Limit', '100');
res.set('X-RateLimit-Remaining', result.remaining.toString());
res.set('X-RateLimit-Reset', result.reset.toString());
```

**Python:**
```python
response.headers['X-RateLimit-Limit'] = '100'
response.headers['X-RateLimit-Remaining'] = str(result.remaining)
response.headers['X-RateLimit-Reset'] = str(result.reset)
```

---

## Examples

### TypeScript/JavaScript Examples

- **Express:** `examples/express-app/`
- **Next.js:** `examples/nextjs-app/`
- **Fastify:** `examples/fastify-app/`

Run:
```bash
cd examples/[framework]-app
npm install
npm run dev
```

### Python Examples

- **FastAPI:** `examples/fastapi-app/`
- **Flask:** `examples/flask-app/`
- **Django:** `examples/django-app/`

Run:
```bash
cd examples/[framework]-app
pip install -r requirements.txt
python [main.py|app.py|manage.py]
```

### Go Backend

```bash
cd apps/api
go run main.go
```

---

## Testing (All Languages)

```bash
# Single request
curl http://localhost:8000/api/data

# Flood requests (get rate limited)
for i in {1..105}; do curl http://localhost:8000/api/data; done

# Check headers
curl -i http://localhost:8000/api/data
```

---

## Architecture

All languages use the same architecture:

```
┌─ Core Logic ─────────────────────┐
│  • Memory limiter (default)      │
│  • Hosted limiter (with API key) │
│  • Fixed-window counter          │
└──────────────────────────────────┘
         ↓
┌─ Framework Integration ──────────┐
│  • Express middleware (TS)       │
│  • FastAPI middleware (Python)   │
│  • Flask hooks (Python)          │
│  • Django middleware (Python)    │
│  • HTTP API (Go)                 │
└──────────────────────────────────┘
         ↓
┌─ Application ────────────────────┐
│  Return 429 if rate limited      │
│  Continue if allowed             │
└──────────────────────────────────┘
```

---

## Deployment

### TypeScript/JavaScript

```bash
npm install --prod
npm start
# or
vercel deploy  # For Vercel
```

### Python

```bash
pip install limity
python app.py
# or
gunicorn app:app  # For production
```

### Go

```bash
go build -o limity-api
./limity-api
# Set environment variables before running
```

---

## Summary

✅ **Works everywhere** - TypeScript, Python, Go  
✅ **Same API** - Consistent behavior across languages  
✅ **Same response format** - Identical structure  
✅ **Same modes** - Memory + hosted  
✅ **Same defaults** - 100 req/60 sec  
✅ **Easy to switch** - Just set `RATE_LIMIT_API_KEY`  

Pick your language and framework, and you're rate-limiting! 🚀
