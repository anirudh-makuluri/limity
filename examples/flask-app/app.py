#!/usr/bin/env python3
"""Flask server with Limity rate limiting."""

import time
from flask import Flask, jsonify, request
from limity import rate_limit

app = Flask(__name__)


@app.before_request
def check_rate_limit():
    """Global rate limiting."""
    # Skip health checks and assets
    if request.path in ["/health"]:
        return
    
    client_ip = request.remote_addr or "unknown"
    
    result = rate_limit(
        key=f"global:{client_ip}",
        limit=1000,
        window=60,
    )
    
    if not result.allowed:
        response = jsonify({
            "error": "Too many requests",
            "retry_after": result.reset - int(time.time()),
        })
        response.status_code = 429
        response.headers["X-RateLimit-Limit"] = "1000"
        response.headers["X-RateLimit-Remaining"] = "0"
        response.headers["X-RateLimit-Reset"] = str(result.reset)
        response.headers["Retry-After"] = str(result.reset - int(time.time()))
        return response


@app.after_request
def add_rate_limit_headers(response):
    """Add rate limit headers to response."""
    # Store in Flask globals if we checked in before_request
    if hasattr(request, "environ"):
        # Headers already added in check_rate_limit
        pass
    
    return response


@app.get("/health")
def health():
    """Health check (not rate limited)."""
    return jsonify({"status": "ok"})


@app.get("/api/data")
def get_data():
    """Get data endpoint (100 req/min per IP)."""
    client_ip = request.remote_addr or "unknown"
    
    result = rate_limit(
        key=f"data:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        response = jsonify({"error": "Too many requests"})
        response.status_code = 429
        return response
    
    response = jsonify({
        "message": "Hello from Limity!",
        "timestamp": time.time(),
        "ip": client_ip,
        "rate_limit": {
            "limit": 100,
            "remaining": result.remaining,
            "reset": result.reset,
        },
    })
    
    response.headers["X-RateLimit-Limit"] = "100"
    response.headers["X-RateLimit-Remaining"] = str(result.remaining)
    response.headers["X-RateLimit-Reset"] = str(result.reset)
    
    return response


@app.get("/api/limited")
def limited_endpoint():
    """Limited endpoint (10 req/min per IP)."""
    client_ip = request.remote_addr or "unknown"
    
    result = rate_limit(
        key=f"limited:{client_ip}",
        limit=10,
        window=60,
    )
    
    if not result.allowed:
        response = jsonify({"error": "Rate limit exceeded for this endpoint"})
        response.status_code = 429
        return response
    
    return jsonify({
        "message": "This endpoint has strict limits",
        "limit": 10,
        "window": 60,
        "remaining": result.remaining,
        "reset": result.reset,
    })


@app.post("/api/create")
def create_resource():
    """Create resource (20 req/min per user)."""
    client_ip = request.remote_addr or "unknown"
    
    # Try to get user ID from JSON body
    user_id = None
    try:
        user_id = request.get_json().get("userId")
    except:
        pass
    
    # Rate limit by user ID if provided, otherwise by IP
    key = f"user:{user_id}" if user_id else f"ip:{client_ip}"
    
    result = rate_limit(
        key=f"create:{key}",
        limit=20,
        window=60,
    )
    
    if not result.allowed:
        response = jsonify({"error": "Too many requests"})
        response.status_code = 429
        return response
    
    return jsonify({
        "success": True,
        "id": f"resource_{int(time.time() * 1000)}",
        "created_at": time.time(),
        "rate_limit": {
            "remaining": result.remaining,
            "reset": result.reset,
        },
    }), 201


if __name__ == "__main__":
    print("""
╔════════════════════════════════════════════════════════╗
║       Limity Flask Example Server                      ║
╚════════════════════════════════════════════════════════╝

Starting on http://localhost:5000

Endpoints:
  GET  /health           - Health check (not rate limited)
  GET  /api/data         - Get data (100 req/min)
  GET  /api/limited      - Limited endpoint (10 req/min)
  POST /api/create       - Create data (20 req/min)

Global Rate Limit: 1000 requests per 60 seconds per IP

Try:
  curl http://localhost:5000/api/data
  
  # Flood requests
  for i in {1..15}; do curl http://localhost:5000/api/data; done

Environment:
  LIMITY_API_KEY - Use hosted rate limiting (optional)
  LIMITY_BASE_URL - Override hosted API URL (optional)
    """)
    
    app.run(debug=True, host="0.0.0.0", port=5000)
