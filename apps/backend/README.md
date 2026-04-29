# Limity API Backend

Simple, minimal rate limiting backend using Go, PostgreSQL, and Upstash Redis.

## Setup

1. Create `.env` from `.env.example`
2. Set your environment variables:
   - `DATABASE_URL` - PostgreSQL connection string (for API key management)
   - `UPSTASH_REDIS_REST_URL` - Upstash Redis REST API URL
   - `UPSTASH_REDIS_REST_TOKEN` - Upstash Redis REST API token
   - `PORT` - Server port (default: 8080)
3. Run: `go run ./cmd/server`

## Database Setup

Create the required tables in your PostgreSQL database:

```sql
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  auth0_id TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  key TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  revoked_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX api_keys_user_id ON api_keys(user_id);
CREATE INDEX api_keys_key ON api_keys(key);
```

## API Endpoints

### Rate Limiting

**POST /check** - Check if a request is within rate limits
```bash
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{
    "key": "user:123",
    "limit": 10,
    "window": 60
  }'
```

Response:
```json
{
  "allowed": true,
  "remaining": 9,
  "reset": 1714396860
}
```

### API Key Management

**POST /api/keys/generate** - Create a new API key
```bash
curl -X POST http://localhost:8080/api/keys/generate \
  -H "Content-Type: application/json" \
  -d '{"user_id": "auth0|123"}'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "auth0|123",
  "key": "limity_abc123def456...",
  "created_at": "2024-04-28T18:13:37Z",
  "revoked_at": null
}
```

**GET /api/keys** - List all API keys for a user
```bash
curl http://localhost:8080/api/keys?user_id=auth0|123
```

Response:
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "auth0|123",
    "key": "limity_abc123def456...",
    "created_at": "2024-04-28T18:13:37Z",
    "revoked_at": null
  }
]
```

**POST /api/keys/:id/revoke** - Revoke an API key
```bash
curl -X POST http://localhost:8080/api/keys/550e8400-e29b-41d4-a716-446655440000/revoke \
  -H "Content-Type: application/json" \
  -d '{"user_id": "auth0|123"}'
```

Response:
```json
{
  "status": "key revoked"
}
```

**GET /health** - Health check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok"
}
