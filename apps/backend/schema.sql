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

ALTER TABLE users ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "users_select_own_row" ON users;
CREATE POLICY "users_select_own_row"
ON users
FOR SELECT
TO authenticated
USING (external_user_id = auth.uid()::text);
