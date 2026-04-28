# Limity Dashboard

API key management dashboard using React, Auth0, and Supabase.

## Features

- **Auth0 Integration** - Secure OAuth2 authentication
- **Supabase** - PostgreSQL database for API keys
- **React Router** - Client-side routing
- **Tailwind CSS** - Utility-first styling
- **Portless** - Dev server with stable .localhost URLs (no ports!)
- **Full TypeScript** - Type-safe throughout

## Setup

### 1. Install Dependencies

```bash
cd apps/dashboard
pnpm install
```

### 2. Configure Environment Variables

Create a `.env.local` file:

```bash
# Auth0
VITE_AUTH0_DOMAIN=your-auth0-domain.auth0.com
VITE_AUTH0_CLIENT_ID=your-auth0-client-id

# API Backend
VITE_API_URL=http://api.localhost

# Supabase (optional - for direct access)
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
```

### 3. Database Schema (Supabase)

```sql
-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT NOT NULL,
  auth0_id TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- API Keys table
CREATE TABLE api_keys (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  key TEXT NOT NULL UNIQUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  revoked_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX api_keys_user_id ON api_keys(user_id);
CREATE INDEX api_keys_key ON api_keys(key);
```

### 4. Backend API Setup

The dashboard calls a backend API for key management. You have two options:

**Option A: Update Go API** (apps/api)
- Add Auth0 JWT validation middleware
- Add Supabase integration
- Add `/api/keys` endpoints for CRUD operations

**Option B: Create Node.js Backend**
- Create a new backend package
- Setup Auth0 JWT verification
- Connect to Supabase

### 5. Run Development Server

```bash
pnpm dev
```

The dashboard runs with **Portless**, which assigns it a stable `.localhost` URL instead of a port number:

```
http://dashboard.localhost
```

For Auth0 setup, add this redirect URI to your Auth0 app settings:
- `http://dashboard.localhost`

If your backend API also uses portless (e.g., `http://api.localhost`), update `.env.local`:
```bash
VITE_API_URL=http://api.localhost
```

## Project Structure

```
src/
├── components/
│   └── Layout.tsx          # Root layout with navigation
├── pages/
│   ├── Home.tsx            # Landing page
│   └── Dashboard.tsx       # API keys dashboard
├── api/
│   └── keys.ts             # API client for backend
├── app.tsx                 # Router setup
├── entry-client.tsx        # App entry point
└── index.css               # Tailwind styles
```

## API Endpoints

The dashboard expects these endpoints from a backend server:

```
POST   /api/keys/generate    - Create new API key
GET    /api/keys            - List user's API keys
POST   /api/keys/:id/revoke - Revoke an API key
```

Each request must include an `Authorization: Bearer {auth0_token}` header.

In the project directory, you can run:

### `npm run dev`

Runs the app in the development mode.<br>
Open [http://localhost:5173](http://localhost:5173) to view it in the browser.

### `npm run build`

Builds the app for production to the `dist` folder.<br>
It correctly bundles Solid in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br>
Your app is ready to be deployed!

## Deployment

Learn more about deploying your application with the [documentations](https://vite.dev/guide/static-deploy.html)
