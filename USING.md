# Using Limity in Your Projects

Quick guide for consuming Limity packages in your applications.

## 📦 Using Packages Now (Local Development)

### Option A: GitHub (Recommended for Now)

Clone the repo into your project:

```bash
git clone https://github.com/anirudh-makuluri/limity.git
cd your-project
npm install ../path/to/limity/packages/core
npm install ../path/to/limity/packages/node  # or edge
```

Or add to `package.json`:
```json
{
  "dependencies": {
    "@limity/core": "file:../limity/packages/core",
    "@limity/node": "file:../limity/packages/node"
  }
}
```

Then:
```bash
npm install
```

### Option B: Monorepo Workspace (Best for Internal Projects)

If you have your own monorepo:

1. Copy `limity` folder into your workspace
2. Update your `pnpm-workspace.yaml`:
```yaml
packages:
  - "packages/**"
  - "limity/packages/**"
```

3. Install once:
```bash
pnpm install
```

4. Use in your projects:
```typescript
import { rateLimit } from '@limity/core';
```

---

## 🎯 Usage Examples

### Express App

```bash
npm install @limity/node express
```

```typescript
import express from 'express';
import { rateLimit } from '@limity/node';

const app = express();

// Global rate limiting: 100 req/min per IP
app.use(rateLimit({
  limit: 100,
  window: 60,
}));

// Custom rate limiting
app.use(rateLimit({
  keyFn: (req) => req.user?.id || req.ip,  // Rate limit by user ID
  limit: 50,
  window: 60,
  onLimitExceeded: (req, res) => {
    res.status(429).json({ error: 'Too many requests' });
  },
}));

app.get('/api/data', (req, res) => {
  res.json({ data: 'hello' });
});

app.listen(3000);
```

### FastAPI App

```bash
pip install limity
```

```python
from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from limity import rate_limit

app = FastAPI()

@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    result = rate_limit(
        key=f"ip:{request.client.host}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        return JSONResponse(
            {"error": "Too many requests"},
            status_code=429,
        )
    
    return await call_next(request)

@app.get("/api/data")
async def get_data():
    return {"data": "hello"}
```

### Raw Core Usage

```typescript
import { rateLimit } from '@limity/core';

// Simple check
const result = await rateLimit({
  key: 'user:123',
  limit: 100,
  window: 60,
});

if (!result.allowed) {
  console.log(`Rate limited. Try again in ${result.reset - Date.now() / 1000}s`);
} else {
  console.log(`Allowed. ${result.remaining} requests remaining`);
}
```

### Edge/Serverless

```typescript
import { checkRateLimit } from '@limity/edge';

export default {
  async fetch(request: Request) {
    const result = await checkRateLimit(request, {
      limit: 100,
      window: 60,
    });
    
    if (!result.allowed) {
      return new Response('429 Too Many Requests', { status: 429 });
    }
    
    return new Response('ok');
  },
};
```

---

## 🔧 Configuration

### Environment Variables

```bash
# Enable hosted mode (connect to backend API)
export RATE_LIMIT_API_KEY=your_api_key

# Optional: custom API URL (if self-hosted)
export RATE_LIMIT_API_URL=http://localhost:8080
```

### Memory vs Hosted Mode

**Memory Mode (Default)**
```typescript
const result = await rateLimit({ key: 'user:123' });
// Uses in-memory map, ~1ms latency
```

**Hosted Mode**
```bash
export RATE_LIMIT_API_KEY=abc123
```
```typescript
const result = await rateLimit({ key: 'user:123' });
// Uses backend API, ~100-200ms latency, multi-instance safe
```

---

## 🚀 Common Patterns

### Rate Limit Different Endpoints Differently

```typescript
// Express
app.get('/api/public', rateLimit({ limit: 1000 }), publicHandler);
app.post('/api/expensive', rateLimit({ limit: 10 }), expensiveHandler);
```

```python
# FastAPI
@router.get("/api/public")
async def public():
    result = rate_limit("key", limit=1000)
    # ...

@router.post("/api/expensive")
async def expensive():
    result = rate_limit("key", limit=10)
    # ...
```

### Rate Limit by User ID

```typescript
// Express middleware
app.use(rateLimit({
  keyFn: (req) => req.user?.id || req.ip,
}));
```

```python
# FastAPI
from limity import rate_limit

@app.get("/api/data")
async def get_data(request: Request, user_id: str):
    result = rate_limit(f"user:{user_id}", limit=100)
    if not result.allowed:
        raise HTTPException(status_code=429)
```

### Skip Rate Limiting for Specific Routes

```typescript
// Express
app.use(rateLimit({
  skip: (req) => req.path === '/health' || req.path === '/metrics',
}));
```

```python
# FastAPI
@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    if request.url.path in ['/health', '/metrics']:
        return await call_next(request)
    
    result = rate_limit(f"ip:{request.client.host}")
    if not result.allowed:
        return JSONResponse({"error": "Too many requests"}, status_code=429)
    
    return await call_next(request)
```

### Add Rate Limit Headers

```typescript
// Express
app.use((req, res, next) => {
  const originalSend = res.send;
  
  res.send = function(data) {
    if (req.rateLimit) {
      res.set('X-RateLimit-Limit', '100');
      res.set('X-RateLimit-Remaining', req.rateLimit.remaining.toString());
      res.set('X-RateLimit-Reset', req.rateLimit.reset.toString());
    }
    return originalSend.call(this, data);
  };
  
  next();
});
```

```python
# FastAPI
@app.middleware("http")
async def add_rate_limit_headers(request: Request, call_next):
    result = rate_limit(f"ip:{request.client.host}")
    response = await call_next(request)
    
    response.headers['X-RateLimit-Limit'] = '100'
    response.headers['X-RateLimit-Remaining'] = str(result.remaining)
    response.headers['X-RateLimit-Reset'] = str(result.reset)
    
    return response
```

---

## 🔌 Integration Checklist

- [ ] Install package (`npm install @limity/core` or `pip install limity`)
- [ ] Import in your app
- [ ] Add rate limiting middleware
- [ ] Set rate limits per endpoint
- [ ] (Optional) Set `RATE_LIMIT_API_KEY` for hosted mode
- [ ] Test: flood endpoint and verify 429 responses
- [ ] Deploy with proper limits

---

## 📚 Further Reading

- [PYTHON_GUIDE.md](./PYTHON_GUIDE.md) - Python framework examples
- [FRAMEWORK_GUIDE.md](./FRAMEWORK_GUIDE.md) - TypeScript/JavaScript framework examples
- [ARCHITECTURE.md](./ARCHITECTURE.md) - How rate limiting works
- [PUBLISHING.md](./PUBLISHING.md) - How to publish updates

---

## Support

For issues, bugs, or questions:

1. Check [ARCHITECTURE.md](./ARCHITECTURE.md) for how rate limiting works
2. See [FRAMEWORK_GUIDE.md](./FRAMEWORK_GUIDE.md) for your framework
3. Review examples: `examples/[framework]-app/`
4. Open a GitHub issue: https://github.com/anirudh-makuluri/limity/issues

---

## Summary

**You can use Limity in your projects right now:**

```bash
# TypeScript
npm install ../path/to/limity/packages/core

# Python
pip install -e ../path/to/limity/packages/python

# Or wait for npm/PyPI publishing (coming soon!)
npm install @limity/core
pip install limity
```

Start with a simple middleware, adjust limits as needed, and enjoy worry-free rate limiting! 🚀
