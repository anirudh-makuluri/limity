import Fastify, { FastifyRequest, FastifyReply } from 'fastify';
import { rateLimit } from '@limity/core';

// Extend Fastify request to include rate limit info
declare module 'fastify' {
  interface FastifyRequest {
    rateLimit?: {
      allowed: boolean;
      remaining: number;
      reset: number;
    };
  }
}

const fastify = Fastify({
  logger: true,
});

// Global rate limit hook (1000 req/min)
fastify.addHook('onRequest', async (request: FastifyRequest, reply: FastifyReply) => {
  // Skip health checks and assets
  if (
    request.url.startsWith('/health') ||
    request.url.startsWith('/metrics')
  ) {
    return;
  }

  const ip = request.headers['x-forwarded-for']?.toString().split(',')[0] ||
             request.ip;

  const result = await rateLimit({
    key: `global:${ip}`,
    limit: 1000,
    window: 60,
  });

  // Set headers on reply
  reply.header('X-RateLimit-Limit', '1000');
  reply.header('X-RateLimit-Remaining', result.remaining.toString());
  reply.header('X-RateLimit-Reset', result.reset.toString());

  if (!result.allowed) {
    reply.code(429).send({
      error: 'Too many requests',
      retryAfter: result.reset - Math.floor(Date.now() / 1000),
    });
    return;
  }

  // Attach to request for downstream access
  request.rateLimit = result;
});

// Health check (not rate limited)
fastify.get('/health', async (request, reply) => {
  reply.send({ status: 'ok' });
});

// GET /api/data - Moderate rate limit (100 req/min)
fastify.get('/api/data', async (request: FastifyRequest, reply: FastifyReply) => {
  const ip = request.headers['x-forwarded-for']?.toString().split(',')[0] ||
             request.ip;

  const result = await rateLimit({
    key: `data:${ip}`,
    limit: 100,
    window: 60,
  });

  if (!result.allowed) {
    reply.code(429).send({
      error: 'Too many requests',
      retryAfter: result.reset - Math.floor(Date.now() / 1000),
    });
    return;
  }

  const { remaining } = request.rateLimit || result;

  reply.send({
    message: 'Hello from Limity!',
    data: {
      timestamp: new Date().toISOString(),
      ip,
    },
    rateLimit: {
      limit: 100,
      remaining,
      reset: result.reset,
      resetAt: new Date(result.reset * 1000).toISOString(),
    },
  });
});

// GET /api/limited - Strict rate limit (10 req/min)
fastify.get('/api/limited', async (request: FastifyRequest, reply: FastifyReply) => {
  const ip = request.headers['x-forwarded-for']?.toString().split(',')[0] ||
             request.ip;

  const result = await rateLimit({
    key: `limited:${ip}`,
    limit: 10,
    window: 60,
  });

  if (!result.allowed) {
    reply.code(429).send({
      error: 'Rate limit exceeded for this endpoint',
      limit: 10,
      window: 60,
      retryAfter: result.reset - Math.floor(Date.now() / 1000),
    });
    return;
  }

  reply.send({
    message: 'This endpoint has strict rate limits',
    limit: 10,
    window: 60,
    remaining: result.remaining,
    resetAt: new Date(result.reset * 1000).toISOString(),
  });
});

// POST /api/create - Custom rate limit based on user
fastify.post('/api/create', async (request: FastifyRequest, reply: FastifyReply) => {
  const ip = request.headers['x-forwarded-for']?.toString().split(',')[0] ||
             request.ip;

  // Try to get user ID from body or headers
  let key = ip;
  try {
    const body = request.body as any;
    if (body?.userId) {
      key = `user:${body.userId}`;
    }
  } catch {
    // Ignore parsing errors
  }

  const result = await rateLimit({
    key: `create:${key}`,
    limit: 20,
    window: 60,
  });

  if (!result.allowed) {
    reply.code(429).send({
      error: 'Too many requests',
      message: 'Please wait before creating another resource',
    });
    return;
  }

  reply.code(201).send({
    success: true,
    message: 'Resource created',
    id: Math.random().toString(36).substr(2, 9),
    rateLimit: {
      remaining: result.remaining,
      reset: result.reset,
    },
  });
});

// Start server
const start = async () => {
  try {
    const port = parseInt(process.env.PORT || '3000', 10);
    await fastify.listen({ port, host: '0.0.0.0' });

    console.log(`
╔════════════════════════════════════════════════════════╗
║         Limity Fastify Example Server                  ║
╚════════════════════════════════════════════════════════╝

Server running on http://localhost:${port}

Endpoints:
  GET  /health           - Health check (not rate limited)
  GET  /api/data         - Get data (100 req/min)
  GET  /api/limited      - Limited endpoint (10 req/min)
  POST /api/create       - Create data (20 req/min)

Global Rate Limit: 1000 requests per 60 seconds per IP

Try these commands:
  
  # Single request
  curl http://localhost:${port}/api/data
  
  # Create resource
  curl -X POST http://localhost:${port}/api/create \\
    -H "Content-Type: application/json" \\
    -d '{"userId":"user123"}'
  
  # Flood with requests
  for i in {1..15}; do curl http://localhost:${port}/api/data; done

Environment Variables:
  - RATE_LIMIT_API_KEY    Use hosted rate limiting
  - PORT                  Server port (default: 3000)

Without RATE_LIMIT_API_KEY, uses fast in-memory rate limiter.
    `);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();
