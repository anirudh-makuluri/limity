# Limity API Backend

Go backend for rate-limiting checks and user profile/API key retrieval.

## Setup

1. Create `.env`.
2. Configure:

- `DATABASE_URL` - Postgres connection string
- `UPSTASH_REDIS_REST_URL` - Upstash Redis URL
- `UPSTASH_REDIS_REST_TOKEN` - Upstash Redis token
- `SUPABASE_URL` - e.g. `https://your-project-ref.supabase.co`
- `SUPABASE_ANON_KEY` - Supabase anon key (used to validate user access tokens)
- `PORT` - default `8080`

3. Run:

```bash
go run ./cmd/server
```

## Database schema

```sql
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  external_user_id TEXT NOT NULL UNIQUE,
  api_key TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
```

## API

- `POST /check`
- `GET /health`
- `GET /api/me`

API key endpoints require:

```http
Authorization: Bearer <supabase_access_token>
```

`GET /api/me` auto-creates an API key for the authenticated user if one does not exist yet.
