# Using Limity in Python

Limity's Python SDK works with any Python web framework. Pure Python, zero framework dependencies.

## Installation

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
print(f"Reset: {result.reset}")
```

---

## FastAPI

```python
from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from limity import rate_limit

app = FastAPI()

@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    """Global rate limiting."""
    client_ip = request.client.host
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        return JSONResponse(
            {"error": "Too many requests"},
            status_code=429,
            headers={
                "X-RateLimit-Limit": "100",
                "X-RateLimit-Remaining": "0",
                "X-RateLimit-Reset": str(result.reset),
            },
        )
    
    response = await call_next(request)
    response.headers["X-RateLimit-Remaining"] = str(result.remaining)
    response.headers["X-RateLimit-Reset"] = str(result.reset)
    
    return response

@app.get("/api/data")
async def get_data():
    return {"data": "hello"}
```

**Example:** `examples/fastapi-app/`

---

## Flask

```python
from flask import Flask, jsonify, request
from limity import rate_limit

app = Flask(__name__)

@app.before_request
def check_rate_limit():
    """Global rate limiting."""
    client_ip = request.remote_addr
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        response = jsonify({"error": "Too many requests"})
        response.status_code = 429
        response.headers["X-RateLimit-Remaining"] = "0"
        response.headers["X-RateLimit-Reset"] = str(result.reset)
        return response

@app.get("/api/data")
def get_data():
    return jsonify({"data": "hello"})
```

**Example:** `examples/flask-app/`

---

## Django

```python
# middleware.py
from django.http import JsonResponse
from django.utils.deprecation import MiddlewareMixin
from limity import rate_limit

class RateLimitMiddleware(MiddlewareMixin):
    """Global rate limiting."""
    
    def process_request(self, request):
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

# settings.py
MIDDLEWARE = [
    # ... other middleware
    "myapp.middleware.RateLimitMiddleware",
]

# views.py
from django.http import JsonResponse
from limity import rate_limit

def get_data(request):
    return JsonResponse({"data": "hello"})
```

**Example:** `examples/django-app/`

---

## Starlette

```python
from starlette.applications import Starlette
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.responses import JSONResponse
from limity import rate_limit

class RateLimitMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, call_next):
        client_ip = request.client.host
        
        result = rate_limit(
            key=f"ip:{client_ip}",
            limit=100,
            window=60,
        )
        
        if not result.allowed:
            return JSONResponse(
                {"error": "Too many requests"},
                status_code=429,
            )
        
        response = await call_next(request)
        response.headers["X-RateLimit-Remaining"] = str(result.remaining)
        return response

app = Starlette()
app.add_middleware(RateLimitMiddleware)
```

---

## Quart (Async Flask)

```python
from quart import Quart, jsonify, request
from limity import rate_limit

app = Quart(__name__)

@app.before_request
async def check_rate_limit():
    """Global rate limiting."""
    client_ip = request.remote_addr
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        response = await jsonify({"error": "Too many requests"})
        response.status_code = 429
        return response

@app.get("/api/data")
async def get_data():
    return jsonify({"data": "hello"})
```

---

## Tornado

```python
import tornado.web
from limity import rate_limit

class BaseHandler(tornado.web.RequestHandler):
    async def prepare(self):
        """Check rate limit before handling request."""
        client_ip = self.request.remote_ip
        
        result = rate_limit(
            key=f"ip:{client_ip}",
            limit=100,
            window=60,
        )
        
        if not result.allowed:
            self.set_status(429)
            self.write({"error": "Too many requests"})
            self.finish()

class DataHandler(BaseHandler):
    async def get(self):
        self.write({"data": "hello"})

app = tornado.web.Application([
    (r"/api/data", DataHandler),
])
```

---

## Bottle

```python
from bottle import Bottle, request, response
from limity import rate_limit

app = Bottle()

@app.hook("before_request")
def check_rate_limit():
    """Global rate limiting."""
    client_ip = request.environ.get("REMOTE_ADDR")
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        response.status = 429
        return {"error": "Too many requests"}

@app.get("/api/data")
def get_data():
    return {"data": "hello"}
```

---

## CherryPy

```python
import cherrypy
from limity import rate_limit

class RateLimitPlugin(cherrypy.plugins.SimplePlugin):
    def __init__(self, bus):
        cherrypy.plugins.SimplePlugin.__init__(self, bus)
    
    def start(self):
        self.bus.subscribe("before_handler", self.check)
    
    def check(self):
        client_ip = cherrypy.request.remote.ip
        
        result = rate_limit(
            key=f"ip:{client_ip}",
            limit=100,
            window=60,
        )
        
        if not result.allowed:
            cherrypy.response.status = 429
            return {"error": "Too many requests"}

class Root:
    @cherrypy.expose
    @cherrypy.tools.json_out()
    def api_data(self):
        return {"data": "hello"}

# Setup
cherrypy.engine.subscribe("start", RateLimitPlugin(cherrypy.engine).start)
cherrypy.quickstart(Root())
```

---

## APIFlask

```python
from apiflask import APIFlask
from limity import rate_limit

app = APIFlask(__name__)

@app.before_request
def check_rate_limit():
    """Global rate limiting."""
    from flask import request
    client_ip = request.remote_addr
    
    result = rate_limit(
        key=f"ip:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        from flask import jsonify
        response = jsonify({"error": "Too many requests"})
        response.status_code = 429
        return response

@app.get("/api/data")
def get_data():
    return {"data": "hello"}
```

---

## Common Patterns

### Extract Client IP

```python
# FastAPI
request.client.host

# Flask
request.remote_addr

# Django
request.META.get("REMOTE_ADDR")

# Starlette
request.client.host

# Tornado
self.request.remote_ip

# Bottle
request.environ.get("REMOTE_ADDR")
```

### Rate Limit by User ID

```python
# Instead of IP, use user ID
user_id = request.user.id  # or from token, header, etc.

result = rate_limit(
    key=f"user:{user_id}",
    limit=100,
    window=60,
)
```

### Different Limits Per Endpoint

```python
@app.get("/api/public")
def public_endpoint():
    result = rate_limit("key", limit=1000, window=60)
    # ...

@app.get("/api/premium")
def premium_endpoint():
    result = rate_limit("key", limit=100, window=60)
    # ...

@app.post("/api/expensive")
def expensive_endpoint():
    result = rate_limit("key", limit=10, window=60)
    # ...
```

### Skip Rate Limiting

```python
# FastAPI middleware
if request.url.path in ["/health", "/metrics"]:
    return await call_next(request)

# Flask before_request
if request.path in ["/health", "/metrics"]:
    return

# Django middleware
if request.path in ["/health", "/metrics"]:
    return None
```

### Add Response Headers

```python
# Standard rate limit headers
response.headers["X-RateLimit-Limit"] = str(limit)
response.headers["X-RateLimit-Remaining"] = str(result.remaining)
response.headers["X-RateLimit-Reset"] = str(result.reset)

# On 429 response
response.headers["Retry-After"] = str(result.reset - time.time())
```

---

## Environment Variables

### Enable Hosted Rate Limiting

```bash
export LIMITY_API_KEY=your_api_key
python app.py
```

Without it, uses fast in-memory limiter.

---

## Testing

### Single request

```bash
curl http://localhost:8000/api/data
```

### Flood requests

```bash
for i in {1..105}; do curl http://localhost:8000/api/data; done
```

### Check headers

```bash
curl -i http://localhost:8000/api/data
```

Look for:
- `X-RateLimit-Limit` - Max requests
- `X-RateLimit-Remaining` - Requests left
- `X-RateLimit-Reset` - Unix timestamp
- `Retry-After` - Seconds to wait (on 429)

---

## API Reference

### `rate_limit(key, limit=100, window=60)`

Check if request is allowed.

**Parameters:**
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
    wait_time = result.reset - int(time.time())
    print(f"Try again in {wait_time} seconds")
```

---

## Performance

- **Memory mode:** ~1ms per check
- **Hosted mode:** ~100-200ms per check

Memory mode is suitable for single-instance deployments. Use hosted mode for scaling.

---

## Examples

Complete working examples:

- **FastAPI:** `examples/fastapi-app/main.py`
- **Flask:** `examples/flask-app/app.py`
- **Django:** `examples/django-app/manage.py`

Run them:

```bash
# FastAPI
cd examples/fastapi-app
pip install -r requirements.txt
python main.py

# Flask
cd examples/flask-app
pip install -r requirements.txt
python app.py

# Django
cd examples/django-app
pip install -r requirements.txt
python manage.py runserver
```

---

## Summary

✅ Works with any Python web framework  
✅ Same API everywhere (FastAPI, Flask, Django, etc.)  
✅ Zero framework dependencies  
✅ Auto-switches from memory to hosted mode  
✅ Never throws errors  

Choose your framework above, copy the pattern, and you're rate-limiting! 🚀
