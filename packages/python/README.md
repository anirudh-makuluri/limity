# Limity Python

Pure Python rate limiting library. Zero framework dependencies.

## Install

```bash
pip install limity
```

## Quick Start

```python
from limity import rate_limit

# Check rate limit
result = rate_limit(
    key="user:123",
    limit=100,
    window=60,
)

if not result.allowed:
    return error(429, "Too many requests")

print(f"Remaining: {result.remaining}")
print(f"Reset at: {result.reset}")
```

## Usage with Frameworks

### FastAPI

```python
from fastapi import FastAPI, HTTPException, Request
from limity import rate_limit

app = FastAPI()

@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    client_ip = request.client.host
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        raise HTTPException(status_code=429, detail="Too many requests")
    
    response = await call_next(request)
    response.headers["X-RateLimit-Limit"] = "100"
    response.headers["X-RateLimit-Remaining"] = str(result.remaining)
    response.headers["X-RateLimit-Reset"] = str(result.reset)
    
    return response

@app.get("/api/data")
async def get_data():
    return {"data": "hello"}
```

### Flask

```python
from flask import Flask, jsonify, request
from limity import rate_limit

app = Flask(__name__)

@app.before_request
def check_rate_limit():
    client_ip = request.remote_addr
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        return jsonify({"error": "Too many requests"}), 429

@app.get("/api/data")
def get_data():
    return {"data": "hello"}
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
        client_ip = request.META.get("REMOTE_ADDR")
        
        result = rate_limit(
            key=f"ip:{client_ip}",
            limit=100,
            window=60,
        )
        
        if not result.allowed:
            return JsonResponse(
                {"error": "Too many requests"},
                status=429,
            )
        
        response = self.get_response(request)
        response["X-RateLimit-Limit"] = "100"
        response["X-RateLimit-Remaining"] = str(result.remaining)
        response["X-RateLimit-Reset"] = str(result.reset)
        
        return response
```

## API

### `rate_limit(key, limit=100, window=60)`

Check if a request is allowed.

**Arguments:**
- `key` (str): Unique identifier (user ID, IP, etc)
- `limit` (int): Max requests per window (default: 100)
- `window` (int): Window duration in seconds (default: 60)

**Returns:**
- `RateLimitResult` with `.allowed`, `.remaining`, `.reset`

**Example:**
```python
result = rate_limit("user:123", limit=50, window=60)

if result.allowed:
    # Process request
    pass
else:
    # Rate limited
    print(f"Retry after {result.reset - time.time()} seconds")
```

## RateLimitResult

```python
class RateLimitResult:
    allowed: bool      # Whether the request is allowed
    remaining: int     # Requests left in window
    reset: int         # Unix timestamp when window resets

# Convert to dict
result.to_dict()
# {"allowed": True, "remaining": 99, "reset": 1714396860}
```

## Modes

### Memory Limiter (Default)

Fast, in-process rate limiting using dict.

```python
result = rate_limit("user:123")
# Uses memory (no RATE_LIMIT_API_KEY set)
```

**Pros:**
- ~1ms latency
- Zero external dependencies
- Works everywhere

**Cons:**
- Limited to single instance
- No persistence across restarts

**Best for:** Development, single-instance deployments, local testing

### Hosted Mode

Scale across multiple instances using hosted API.

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

Then:
```python
result = rate_limit("user:123")
# Automatically uses hosted API
# Falls back to memory on API failure
```

**Pros:**
- Works across multiple instances
- Scales infinitely
- Persistent across restarts

**Cons:**
- ~100-200ms latency
- Requires API key

**Best for:** Production, distributed systems, scaling

## Error Handling

The library never throws. Instead:

- On API error → Fail closed (return `allowed=False`)
- On network error → Fail closed (return `allowed=False`)
- On memory limit → Process normally (memory limits are checked)

This ensures your app stays up even if the rate limiter fails.

## Environment

### `RATE_LIMIT_API_KEY`

API key for hosted rate limiting. If set, enables hosted mode with automatic fallback.

```bash
export RATE_LIMIT_API_KEY=your_api_key
```

Without it, uses memory limiter.

## Performance

- **Memory mode:** ~1ms per check
- **Hosted mode:** ~100-200ms per check

Memory mode is suitable for single-instance deployments. Use hosted mode for scaling.

## Examples

See `examples/` for complete applications:

- `flask-app/` - Flask with rate limiting
- `fastapi-app/` - FastAPI with rate limiting
- `django-app/` - Django with rate limiting

## Development

```bash
# Install with dev dependencies
pip install -e ".[dev]"

# Run tests
pytest

# Run with coverage
pytest --cov

# Format code
black limity/

# Type checking
mypy limity/
```

## License

MIT
