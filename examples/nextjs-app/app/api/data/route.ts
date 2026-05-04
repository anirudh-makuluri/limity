import { checkRateLimit } from '@/lib/rate-limit';
import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  const ip = request.headers.get('x-forwarded-for') || 
             request.headers.get('x-real-ip') || 
             'unknown';

  // Check rate limit
  const result = await checkRateLimit({
    key: `data:${ip}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    return NextResponse.json(
      {
        error: 'Too many requests',
        retryAfter: result.reset - Math.floor(Date.now() / 1000),
      },
      {
        status: 429,
        headers: {
          'X-RateLimit-Limit': '100',
          'X-RateLimit-Remaining': result.remaining.toString(),
          'X-RateLimit-Reset': result.reset.toString(),
          'Retry-After': (result.reset - Math.floor(Date.now() / 1000)).toString(),
        },
      }
    );
  }

  return NextResponse.json(
    {
      message: 'Hello from Limity!',
      data: {
        timestamp: new Date().toISOString(),
        ip,
      },
      rateLimit: {
        limit: 100,
        remaining: result.remaining,
        reset: result.reset,
        resetAt: new Date(result.reset * 1000).toISOString(),
      },
    },
    {
      headers: {
        'X-RateLimit-Limit': '100',
        'X-RateLimit-Remaining': result.remaining.toString(),
        'X-RateLimit-Reset': result.reset.toString(),
      },
    }
  );
}
