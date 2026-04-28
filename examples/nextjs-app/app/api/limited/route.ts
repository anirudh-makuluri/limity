import { rateLimit } from '@limity/core';
import { NextRequest, NextResponse } from 'next/server';

// More restrictive limit for this endpoint
const LIMIT = 10;
const WINDOW = 60;

export async function GET(request: NextRequest) {
  const ip = request.headers.get('x-forwarded-for') || 
             request.headers.get('x-real-ip') || 
             'unknown';

  // Check rate limit with stricter limits
  const result = await rateLimit({
    key: `limited:${ip}`,
    limit: LIMIT,
    window: WINDOW,
  });

  if (!result.allowed) {
    return NextResponse.json(
      {
        error: 'Rate limit exceeded for this endpoint',
        limit: LIMIT,
        window: WINDOW,
        retryAfter: result.reset - Math.floor(Date.now() / 1000),
      },
      {
        status: 429,
        headers: {
          'X-RateLimit-Limit': LIMIT.toString(),
          'X-RateLimit-Remaining': '0',
          'X-RateLimit-Reset': result.reset.toString(),
          'Retry-After': (result.reset - Math.floor(Date.now() / 1000)).toString(),
        },
      }
    );
  }

  return NextResponse.json(
    {
      message: 'This endpoint has strict rate limits',
      limit: LIMIT,
      window: WINDOW,
      remaining: result.remaining,
      resetAt: new Date(result.reset * 1000).toISOString(),
    },
    {
      headers: {
        'X-RateLimit-Limit': LIMIT.toString(),
        'X-RateLimit-Remaining': result.remaining.toString(),
        'X-RateLimit-Reset': result.reset.toString(),
      },
    }
  );
}
