CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  external_user_id TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL,
  api_key TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS users_external_user_id_idx ON users(external_user_id);
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);

CREATE TABLE IF NOT EXISTS request_events (
  id BIGSERIAL PRIMARY KEY,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
  method TEXT NOT NULL,
  route TEXT NOT NULL,
  path TEXT NOT NULL,
  status_code INT NOT NULL,
  duration_ms BIGINT NOT NULL,
  client_ip TEXT,
  user_agent TEXT,
  owner_user_id TEXT
);

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND table_name = 'request_events'
      AND column_name = 'external_user_id'
  ) AND NOT EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND table_name = 'request_events'
      AND column_name = 'owner_user_id'
  ) THEN
    ALTER TABLE request_events RENAME COLUMN external_user_id TO owner_user_id;
  END IF;
END$$;

CREATE INDEX IF NOT EXISTS request_events_timestamp_idx ON request_events(timestamp DESC);
CREATE INDEX IF NOT EXISTS request_events_route_timestamp_idx ON request_events(route, timestamp DESC);
CREATE INDEX IF NOT EXISTS request_events_status_timestamp_idx ON request_events(status_code, timestamp DESC);
CREATE INDEX IF NOT EXISTS request_events_owner_user_id_timestamp_idx ON request_events(owner_user_id, timestamp DESC);

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE request_events ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "users_select_own_row" ON users;
CREATE POLICY "users_select_own_row"
ON users
FOR SELECT
TO authenticated
USING (external_user_id = auth.uid()::text);

DROP POLICY IF EXISTS "request_events_select_own_rows" ON request_events;
CREATE POLICY "request_events_select_own_rows"
ON request_events
FOR SELECT
TO authenticated
USING (owner_user_id = auth.uid()::text);
