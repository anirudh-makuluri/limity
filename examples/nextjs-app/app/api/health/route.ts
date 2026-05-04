import { checkRateLimit } from '@/lib/rate-limit';
import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  const ip = request.headers.get('x-forwarded-for') || 
             request.headers.get('x-real-ip') || 
             'unknown';

  // Rate limit by user ID if provided, otherwise by IP
  let key = ip;
  try {
    const body = await request.json();
    if (body.userId) {
      key = `user:${body.userId}`;
    }
  } catch {
    // Body parsing failed, use IP
  }

  // Check rate limit
  const result = await checkRateLimit({
    key: `create:${key}`,
    limit: 20,
    window: 60,
  });

  if (!result.allowed) {
    return NextResponse.json(
      {
        error: 'Too many requests',
        message: 'Please wait before creating another resource',
      },
      { status: 429 }
    );
  }

  return NextResponse.json(
    {
      success: true,
      message: 'Resource created',
      id: Math.random().toString(36).substr(2, 9),
      rateLimit: {
        remaining: result.remaining,
        reset: result.reset,
      },
    },
    { status: 201 }
  );
}
