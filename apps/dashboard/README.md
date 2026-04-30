# Limity Dashboard

API key management dashboard using React, Supabase Auth, and the Limity Go backend.

## Features

- Supabase email/password authentication
- API key management UI
- React Router + TypeScript + Tailwind CSS

## Setup

1. Install dependencies:

```bash
cd apps/dashboard
pnpm install
```

2. Create `.env.local`:

```bash
VITE_SUPABASE_URL=https://your-project-ref.supabase.co
VITE_SUPABASE_ANON_KEY=your-supabase-anon-key
VITE_API_URL=http://localhost:8080
```

3. In Supabase Auth settings, enable Email/Password provider.

4. Run:

```bash
pnpm dev
```

## Backend contract

The dashboard expects:

- `GET /api/me`

The backend auto-creates one API key per user on first authenticated access.
Requests include `Authorization: Bearer {supabase_access_token}`.
