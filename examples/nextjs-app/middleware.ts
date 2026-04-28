import { rateLimit } from '@limity/core';
import { NextRequest, NextResponse } from 'next/server';

export async function middleware(request: NextRequest) {
  // Skip rate limiting for health checks and assets
  if (
    request.nextUrl.pathname.startsWith('/health') ||
    request.nextUrl.pathname.startsWith('/_next') ||
    request.nextUrl.pathname.includes('.')
  ) {
    return NextResponse.next();
  }

  const ip = request.headers.get('x-forwarded-for') || 
             request.headers.get('x-real-ip') || 
             'unknown';

  // Global rate limit (optional, adjust as needed)
  const result = await rateLimit({
    key: `global:${ip}`,
    limit: 1000,
    window: 60,
  });

  if (!result.allowed) {
    return NextResponse.json(
      { error: 'Too many requests' },
      { status: 429 }
    );
  }

  // Attach rate limit info to response headers
  const response = NextResponse.next();
  response.headers.set('X-RateLimit-Limit', '1000');
  response.headers.set('X-RateLimit-Remaining', result.remaining.toString());
  response.headers.set('X-RateLimit-Reset', result.reset.toString());

  return response;
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!_next/static|_next/image|favicon.ico).*)',
  ],
};
