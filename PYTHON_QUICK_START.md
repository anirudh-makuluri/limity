# Python SDK Quick Reference

## Installation

```bash
pip install limity
```

## Basic Usage

```python
from limity import rate_limit

result = rate_limit("user:123", limit=100, window=60)

if not result.allowed:
    return error(429, "Too many requests")

print(f"Remaining: {result.remaining}")
```

## Response Object

```python
result.allowed      # bool - whether request is allowed
result.remaining    # int  - requests left in window
result.reset        # int  - unix timestamp when window resets

result.to_dict()    # dict - convert to JSON-serializable dict
```

## Framework Integration

### FastAPI
```python
@app.middleware("http")
async def rate_limit_middleware(request, call_next):
    result = rate_limit(f"ip:{request.client.host}")
    if not result.allowed:
        return JSONResponse({"error": "Too many requests"}, status_code=429)
    return await call_next(request)
```

### Flask
```python
@app.before_request
def check_rate_limit():
    result = rate_limit(f"ip:{request.remote_addr}")
    if not result.allowed:
        return jsonify({"error": "Too many requests"}), 429
```

### Django
```python
class RateLimitMiddleware:
    def __call__(self, request):
        result = rate_limit(f"ip:{request.META.get('REMOTE_ADDR')}")
        if not result.allowed:
            return JsonResponse({"error": "Too many requests"}, status=429)
```

## Modes

### Memory (Default)
```python
result = rate_limit("key")  # Uses in-memory dict
```
- Fast (~1ms)
- No setup needed
- Single instance only

### Hosted
```bash
export LIMITY_API_KEY=your_api_key
```
```python
result = rate_limit("key")  # Uses API with fallback
```
- Scales across instances
- ~100-200ms latency
- Automatic fallback to memory

## Common Patterns

### Rate limit by user ID
```python
result = rate_limit(f"user:{user_id}")
```

### Rate limit by IP
```python
result = rate_limit(f"ip:{request.remote_addr}")
```

### Different limits per endpoint
```python
# Public: 1000 req/min
result = rate_limit(f"public:{ip}", limit=1000, window=60)

# Premium: 100 req/min
result = rate_limit(f"premium:{user_id}", limit=100, window=60)

# Expensive: 10 req/min
result = rate_limit(f"expensive:{user_id}", limit=10, window=60)
```

### Add response headers
```python
response.headers["X-RateLimit-Limit"] = "100"
response.headers["X-RateLimit-Remaining"] = str(result.remaining)
response.headers["X-RateLimit-Reset"] = str(result.reset)

# On 429
response.headers["Retry-After"] = str(result.reset - time.time())
```

## Testing

```bash
# Single request
curl http://localhost:8000/api/data

# Flood (get rate limited)
for i in {1..15}; do curl http://localhost:8000/api/data; done

# Check headers
curl -i http://localhost:8000/api/data
```

## Examples

- **FastAPI:** `examples/fastapi-app/main.py`
- **Flask:** `examples/flask-app/app.py`
- **Django:** `examples/django-app/manage.py`

See [PYTHON_GUIDE.md](./PYTHON_GUIDE.md) for detailed examples with:
- Starlette
- Quart
- Tornado
- Bottle
- CherryPy
- APIFlask
- And more!

## Package Structure

```
packages/python/
├── limity/
│   └── __init__.py       # Core rate limiting logic
├── pyproject.toml        # Modern Python packaging
├── setup.py              # setuptools config
└── README.md             # Full documentation
```

## API Reference

### `rate_limit(key, limit=100, window=60)`

Check if request is allowed.

**Args:**
- `key` (str): Unique identifier
- `limit` (int): Max requests per window (default: 100)
- `window` (int): Window size in seconds (default: 60)

**Returns:** `RateLimitResult`

**Example:**
```python
result = rate_limit("user:123", limit=50, window=60)
if result.allowed:
    process_request()
else:
    return_429_error()
```

## Error Handling

Limity never throws errors:

- On API error → Fail closed (`allowed=False`)
- On network error → Fail closed (`allowed=False`)
- Memory full → Process normally

This ensures your app stays up!

## Performance

- **Memory mode:** ~1ms per check
- **Hosted mode:** ~100-200ms per check

Perfect for:
- Development (memory)
- Production single-instance (memory)
- Distributed systems (hosted with API key)

---

See [PYTHON_GUIDE.md](./PYTHON_GUIDE.md) for framework-specific examples and patterns.
