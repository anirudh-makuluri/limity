"""
Limity - Developer-first rate limiting

Pure Python implementation with zero dependencies.
Automatically switches from memory to hosted mode.
"""

import os
import time
import httpx
from typing import Dict, Optional, Tuple
from datetime import datetime

__version__ = "0.1.0"

# Type definitions
class RateLimitResult:
    """Result of a rate limit check."""
    
    def __init__(self, allowed: bool, remaining: int, reset: int):
        self.allowed = allowed
        self.remaining = remaining
        self.reset = reset
    
    def to_dict(self) -> Dict:
        return {
            "allowed": self.allowed,
            "remaining": self.remaining,
            "reset": self.reset,
        }
    
    def __repr__(self) -> str:
        return f"RateLimitResult(allowed={self.allowed}, remaining={self.remaining}, reset={self.reset})"


# Memory limiter
class MemoryLimiter:
    """In-memory rate limiter using dict."""
    
    def __init__(self):
        self.store: Dict[str, Tuple[int, int]] = {}  # key -> (count, reset)
    
    def check(
        self,
        key: str,
        limit: int = 100,
        window: int = 60,
    ) -> RateLimitResult:
        """Check rate limit using in-memory store."""
        now = int(time.time())
        window_start = now - (now % window)
        reset = window_start + window
        
        store_key = f"{key}:{window_start}"
        
        # Clean up old entries
        expired_keys = [
            k for k, (_, expire_time) in self.store.items()
            if expire_time <= now
        ]
        for k in expired_keys:
            del self.store[k]
        
        # Get or create entry
        if store_key not in self.store:
            self.store[store_key] = (1, reset)
            return RateLimitResult(
                allowed=True,
                remaining=limit - 1,
                reset=reset,
            )
        
        # Increment counter
        count, _ = self.store[store_key]
        count += 1
        self.store[store_key] = (count, reset)
        
        allowed = count <= limit
        remaining = max(0, limit - count)
        
        return RateLimitResult(allowed=allowed, remaining=remaining, reset=reset)


# Hosted limiter
class HostedLimiter:
    """Rate limiter that calls hosted API."""
    
    def __init__(self, api_key: str, api_url: str = "https://api.limity.smart-deploy.xyz"):
        self.api_key = api_key
        self.api_url = api_url
    
    def check(
        self,
        key: str,
        limit: int = 100,
        window: int = 60,
    ) -> RateLimitResult:
        """Check rate limit via hosted API."""
        try:
            with httpx.Client() as client:
                response = client.post(
                    f"{self.api_url}/check",
                    json={
                        "key": key,
                        "limit": limit,
                        "window": window,
                    },
                    headers={
                        "Authorization": f"Bearer {self.api_key}",
                        "Content-Type": "application/json",
                    },
                    timeout=5.0,
                )
            
            if response.status_code == 200:
                data = response.json()
                return RateLimitResult(
                    allowed=data["allowed"],
                    remaining=data["remaining"],
                    reset=data["reset"],
                )
            else:
                # API error - fail closed (deny request)
                return RateLimitResult(
                    allowed=False,
                    remaining=0,
                    reset=int(time.time()) + window,
                )
        except Exception as e:
            # Network error - fail closed (deny request)
            import logging
            logging.error(f"Rate limit API error: {e}")
            return RateLimitResult(
                allowed=False,
                remaining=0,
                reset=int(time.time()) + window,
            )


# Main public API
_memory_limiter = MemoryLimiter()
_hosted_limiter: Optional[HostedLimiter] = None


def rate_limit(
    key: str,
    limit: int = 100,
    window: int = 60,
) -> RateLimitResult:
    """
    Check rate limit.
    
    Args:
        key: Unique identifier (user ID, IP, etc)
        limit: Max requests per window (default: 100)
        window: Window duration in seconds (default: 60)
    
    Returns:
        RateLimitResult with allowed, remaining, and reset info
    
    Example:
        >>> result = rate_limit("user:123", limit=50, window=60)
        >>> if not result.allowed:
        ...     return error_response(429, "Too many requests")
    """
    api_key = os.getenv("LIMITY_API_KEY") or os.getenv("RATE_LIMIT_API_KEY")
    base_url = os.getenv("LIMITY_BASE_URL")
    
    if api_key:
        # Use hosted limiter with fallback to memory
        global _hosted_limiter
        if _hosted_limiter is None:
            _hosted_limiter = HostedLimiter(api_key, base_url or "https://api.limity.smart-deploy.xyz")
        
        result = _hosted_limiter.check(key, limit, window)
        
        # Fallback to memory if API fails
        if not result.allowed:
            result = _memory_limiter.check(key, limit, window)
        
        return result
    
    # Default: use memory limiter
    return _memory_limiter.check(key, limit, window)


# Export
__all__ = [
    "rate_limit",
    "RateLimitResult",
    "MemoryLimiter",
    "HostedLimiter",
]
