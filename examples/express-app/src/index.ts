import express, { Request } from 'express';
import { rateLimit } from '@limity/node';

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(express.json());

// Apply rate limiting to all routes
// Default: 10 requests per 10 seconds per IP (for demo purposes)
app.use(
  rateLimit({
    limit: 10,
    window: 10,
    keyFn: (req) => req.ip || 'unknown',
    skip: (req) => req.path === '/health',
  })
);

// Health check (not rate limited)
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

// Example API endpoint
app.get('/api/data', (req: any, res) => {
  const { remaining, reset } = req.rateLimit;
  const resetTime = new Date(reset * 1000).toISOString();

  res.json({
    message: 'Hello from Limity!',
    data: {
      timestamp: new Date().toISOString(),
      requestCount: 1,
    },
    rateLimit: {
      limit: 10,
      remaining,
      reset,
      resetAt: resetTime,
    },
  });
});

// Example endpoint with custom rate limit
app.get('/api/limited', (req: any, res) => {
  const { remaining, reset } = req.rateLimit;

  res.json({
    message: 'This endpoint is rate limited!',
    remaining,
    reset,
  });
});

// Example endpoint that creates a resource
app.post('/api/data', (req: any, res) => {
  const { remaining, reset } = req.rateLimit;

  res.status(201).json({
    message: 'Resource created',
    id: Math.random().toString(36).substr(2, 9),
    rateLimit: {
      remaining,
      reset,
    },
  });
});

// Error handler for rate limit exceeded
app.use((err: any, req: Request, res: any, next: any) => {
  if (err.status === 429) {
    return res.status(429).json({
      error: 'Too many requests',
      retryAfter: err.retryAfter,
    });
  }
  next(err);
});

app.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════════════════════╗
║           Limity Express Example Server                ║
╚════════════════════════════════════════════════════════╝

Server running on http://localhost:${PORT}

Endpoints:
  GET  /health           - Health check (not rate limited)
  GET  /api/data         - Get data (rate limited)
  POST /api/data         - Create data (rate limited)
  GET  /api/limited      - Limited endpoint (rate limited)

Rate Limit: 10 requests per 10 seconds per IP

Try these commands:
  
  # Single request
  curl http://localhost:${PORT}/api/data
  
  # Flood with requests
  for i in {1..15}; do curl http://localhost:${PORT}/api/data; done

Environment Variables:
  - LIMITY_API_KEY    Use hosted rate limiting
  - LIMITY_BASE_URL       Override hosted API URL (optional)
  - PORT                  Server port (default: 3000)

Without LIMITY_API_KEY, uses fast in-memory rate limiter.
  `);
});
