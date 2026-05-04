#!/usr/bin/env python3
"""Django server with Limity rate limiting."""

import os
import django
from django.conf import settings

# Configure Django settings
if not settings.configured:
    settings.configure(
        DEBUG=True,
        DATABASES={
            'default': {
                'ENGINE': 'django.db.backends.sqlite3',
                'NAME': ':memory:',
            }
        },
        INSTALLED_APPS=[
            'django.contrib.contenttypes',
            'django.contrib.auth',
        ],
        MIDDLEWARE=[
            'limity_django.RateLimitMiddleware',
        ],
        ROOT_URLCONF='__main__',
        SECRET_KEY='dev-secret-key',
    )
    django.setup()

import json
import time
from django.http import JsonResponse
from django.views.decorators.http import require_http_methods
from limity import rate_limit


# Rate limiting middleware
class RateLimitMiddleware:
    """Global rate limiting middleware."""
    
    def __init__(self, get_response):
        self.get_response = get_response
    
    def __call__(self, request):
        # Skip health checks
        if request.path in ["/health"]:
            return self.get_response(request)
        
        client_ip = request.META.get("REMOTE_ADDR", "unknown")
        
        result = rate_limit(
            key=f"global:{client_ip}",
            limit=1000,
            window=60,
        )
        
        if not result.allowed:
            response = JsonResponse({
                "error": "Too many requests",
                "retry_after": result.reset - int(time.time()),
            }, status=429)
            response["X-RateLimit-Limit"] = "1000"
            response["X-RateLimit-Remaining"] = "0"
            response["X-RateLimit-Reset"] = str(result.reset)
            response["Retry-After"] = str(result.reset - int(time.time()))
            return response
        
        response = self.get_response(request)
        response["X-RateLimit-Limit"] = "1000"
        response["X-RateLimit-Remaining"] = str(result.remaining)
        response["X-RateLimit-Reset"] = str(result.reset)
        
        return response


# Views
def health(request):
    """Health check (not rate limited)."""
    return JsonResponse({"status": "ok"})


def get_data(request):
    """Get data endpoint (100 req/min per IP)."""
    client_ip = request.META.get("REMOTE_ADDR", "unknown")
    
    result = rate_limit(
        key=f"data:{client_ip}",
        limit=100,
        window=60,
    )
    
    if not result.allowed:
        response = JsonResponse({"error": "Too many requests"}, status=429)
        response["X-RateLimit-Limit"] = "100"
        response["X-RateLimit-Remaining"] = "0"
        response["X-RateLimit-Reset"] = str(result.reset)
        return response
    
    response = JsonResponse({
        "message": "Hello from Limity!",
        "timestamp": time.time(),
        "ip": client_ip,
        "rate_limit": {
            "limit": 100,
            "remaining": result.remaining,
            "reset": result.reset,
        },
    })
    
    response["X-RateLimit-Limit"] = "100"
    response["X-RateLimit-Remaining"] = str(result.remaining)
    response["X-RateLimit-Reset"] = str(result.reset)
    
    return response


def limited_endpoint(request):
    """Limited endpoint (10 req/min per IP)."""
    client_ip = request.META.get("REMOTE_ADDR", "unknown")
    
    result = rate_limit(
        key=f"limited:{client_ip}",
        limit=10,
        window=60,
    )
    
    if not result.allowed:
        return JsonResponse(
            {"error": "Rate limit exceeded for this endpoint"},
            status=429,
        )
    
    return JsonResponse({
        "message": "This endpoint has strict limits",
        "limit": 10,
        "window": 60,
        "remaining": result.remaining,
        "reset": result.reset,
    })


@require_http_methods(["POST"])
def create_resource(request):
    """Create resource (20 req/min per user)."""
    client_ip = request.META.get("REMOTE_ADDR", "unknown")
    
    # Try to get user ID from JSON body
    user_id = None
    try:
        user_id = json.loads(request.body).get("userId")
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
        return JsonResponse({"error": "Too many requests"}, status=429)
    
    return JsonResponse({
        "success": True,
        "id": f"resource_{int(time.time() * 1000)}",
        "created_at": time.time(),
        "rate_limit": {
            "remaining": result.remaining,
            "reset": result.reset,
        },
    }, status=201)


# URL routing
from django.urls import path

urlpatterns = [
    path("health", health),
    path("api/data", get_data),
    path("api/limited", limited_endpoint),
    path("api/create", create_resource),
]


if __name__ == "__main__":
    from django.core.management import execute_from_command_line
    
    print("""
╔════════════════════════════════════════════════════════╗
║       Limity Django Example Server                     ║
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
  LIMITY_API_KEY - Use hosted rate limiting (optional)
  LIMITY_BASE_URL - Override hosted API URL (optional)
    """)
    
    # Run development server
    execute_from_command_line(["manage.py", "runserver", "0.0.0.0:8000"])
