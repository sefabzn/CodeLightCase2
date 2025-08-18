# Environment Setup Instructions

## Frontend Environment Variables

Create `frontend/.env.local` with:

```bash
# Supabase Configuration
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key-here

# Backend API
NEXT_PUBLIC_API_URL=http://localhost:8000
```

## Backend Environment Variables

Create `backend/.env` with:

```bash
# Server Configuration
PORT=8000

# Database Configuration  
DATABASE_URL=postgresql://postgres:postgres@localhost:54322/postgres

# Supabase Configuration
SUPABASE_URL=http://localhost:54321
SUPABASE_ANON_KEY=your-supabase-anon-key-here
SUPABASE_SERVICE_ROLE_KEY=your-supabase-service-role-key-here
```

## Getting Supabase Keys

1. Start Supabase: `cd db && npx supabase start`
2. Run: `npx supabase status` to get the actual keys
3. Copy the keys to your environment files

## Docker Desktop Required

Supabase local development requires Docker Desktop to be installed and running.

Download from: https://www.docker.com/products/docker-desktop/
