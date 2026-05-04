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
- `ANALYTICS_ENABLED` - default `true`; set `false` to disable async request event writes
- `ANALYTICS_QUEUE_SIZE` - default `10000`
- `ANALYTICS_BATCH_SIZE` - default `200`
- `ANALYTICS_FLUSH_INTERVAL_MS` - default `1000`
- `ANALYTICS_FLUSH_TIMEOUT_MS` - default `15000`
- `API_KEY_GAUGE_REFRESH_SEC` - default `60`

3. Run:

```bash
go run ./cmd/server
```

Request analytics user attribution is automatic:

- `/api/me`, `/api/me/refresh-key`: from verified Supabase token claims.
- `/check`: resolved asynchronously from `key` by mapping API key owner in the `users` table.

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
- `GET /metrics` (Prometheus)

API key endpoints require:

```http
Authorization: Bearer <supabase_access_token>
```

`GET /api/me` auto-creates an API key for the authenticated user if one does not exist yet.

## Metrics

Prometheus metrics include:

- `http_requests_total{route,method,status}`
- `http_request_duration_seconds_bucket{route,method}`
- `limity_check_total{result}`
- `limity_api_keys_total`
- `limity_redis_errors_total`
- `limity_auth_failures_total`
- `analytics_events_dropped_total`
- `limity_owner_lookup_success_total`
- `limity_owner_lookup_miss_total`
- `limity_owner_lookup_error_total`

Useful PromQL for `/check`:

```promql
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{route="/check"}[5m])) by (le))
```

```promql
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{route="/check"}[5m])) by (le))
```
