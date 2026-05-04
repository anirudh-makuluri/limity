#!/usr/bin/env python3
"""FastAPI server with Limity rate limiting."""

import time
from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from limity import rate_limit

app = FastAPI(title="Limity FastAPI Example")


@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    """Global rate limiting middleware."""
    # Skip health checks
    if request.url.path in ["/health", "/docs", "/openapi.json"]:
        return await call_next(request)
    
    client_ip = request.client.host if request.client else "unknown"
    
    result = rate_limit(
        key=f"global:{client_ip}",
        limit=1000,
        window=60,
    )
    
    if not result.allowed:
        return JSONResponse(
            {"error": "Too many requests", "retry_after": result.reset - int(time.time())},
            status_code=429,
            headers={
                "X-RateLimit-Limit": "1000",
                "X-RateLimit-Remaining": "0",
                "X-RateLimit-Reset": str(result.reset),
                "Retry-After": str(result.reset - int(time.time())),
            },
        )
    
    response = await call_next(request)
    response.headers["X-RateLimit-Limit"] = "1000"
    response.headers["X-RateLimit-Remaining"] = str(result.remaining)
    response.headers["X-RateLimit-Reset"] = str(result.reset)
    
    return response


@app.get("/health")
async def health():
    """Health check (not rate limited)."""
    return {"status": "ok"}


@app.get("/api/data")
async def get_data(request: Request):
    """Get data endpoint (100 req/min per IP)."""
    client_ip = request.client.host if request.client else "unknown"
    
    result = rate_limit(
        key=f"data:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        raise HTTPException(
            status_code=429,
            detail="Too many requests",
        )
    
    return {
        "message": "Hello from Limity!",
        "timestamp": time.time(),
        "ip": client_ip,
        "rate_limit": {
            "limit": 100,
            "remaining": result.remaining,
            "reset": result.reset,
        },
    }


@app.get("/api/limited")
async def limited_endpoint(request: Request):
    """Limited endpoint (10 req/min per IP)."""
    client_ip = request.client.host if request.client else "unknown"
    
    result = rate_limit(
        key=f"limited:{client_ip}",
        limit=10,
        window=60,
    )
    
    if not result.allowed:
        raise HTTPException(
            status_code=429,
            detail="Rate limit exceeded for this endpoint",
        )
    
    return {
        "message": "This endpoint has strict limits",
        "limit": 10,
        "window": 60,
        "remaining": result.remaining,
        "reset": result.reset,
    }


@app.post("/api/create")
async def create_resource(request: Request, user_id: str = None):
    """Create resource (20 req/min per user)."""
    # Rate limit by user ID if provided, otherwise by IP
    key = f"user:{user_id}" if user_id else f"ip:{request.client.host}"
    
    result = rate_limit(
        key=f"create:{key}",
        limit=20,
        window=60,
    )
    
    if not result.allowed:
        raise HTTPException(
            status_code=429,
            detail="Too many requests",
        )
    
    return {
        "success": True,
        "id": f"resource_{int(time.time() * 1000)}",
        "created_at": time.time(),
        "rate_limit": {
            "remaining": result.remaining,
            "reset": result.reset,
        },
    }


if __name__ == "__main__":
    import uvicorn
    
    print("""
╔════════════════════════════════════════════════════════╗
║      Limity FastAPI Example Server                     ║
╚════════════════════════════════════════════════════════╝

Starting on http://localhost:8000

Endpoints:
  GET  /health           - Health check (not rate limited)
  GET  /api/data         - Get data (100 req/min)
  GET  /api/limited      - Limited endpoint (10 req/min)
  POST /api/create       - Create data (20 req/min)

Global Rate Limit: 1000 requests per 60 seconds per IP

Try:
  curl http://localhost:8000/api/data
  
  # Flood requests
  for i in {1..15}; do curl http://localhost:8000/api/data; done

Environment:
  RATE_LIMIT_API_KEY - Use hosted rate limiting (optional)
  LIMITY_BASE_URL - Override hosted API URL (optional)
    """)
    
    uvicorn.run(app, host="0.0.0.0", port=8000)
