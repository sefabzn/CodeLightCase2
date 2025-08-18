# Turkcell Ev+Mobil Paket Danışmanı

Web application that helps users find the most cost-effective combination of mobile, home internet, and TV packages based on household profile and address coverage.

## Architecture

- **Frontend**: Next.js 14 (App Router) + TypeScript + TailwindCSS
- **Backend**: Go REST API with Echo framework
- **Database**: Supabase (PostgreSQL)

## Setup

### Prerequisites

- Node.js 18+ with pnpm
- Go 1.21+
- Supabase CLI

### Quick Start

1. **Clone and install dependencies**:
```bash
git clone <repo-url>
cd CodeLightCase2
```

2. **Setup database**:
```bash
make db-up
```

3. **Start development servers**:
```bash
# Terminal 1 - Backend
make api

# Terminal 2 - Frontend  
make web
```

4. **Access the application**:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8000

## Run

Individual services:

**Unix/Linux/macOS (using Makefile):**
```bash
# Database
make db-up      # Start Supabase
make db-down    # Stop Supabase
make seed       # Load sample data

# Backend
make api        # Start Go server on :8000

# Frontend
make web        # Start Next.js on :3000
```

**Windows (using PowerShell script):**
```powershell
# Database
.\dev.ps1 db-up      # Start Supabase
.\dev.ps1 db-down    # Stop Supabase
.\dev.ps1 seed       # Load sample data

# Backend
.\dev.ps1 api        # Start Go server on :8000

# Frontend
.\dev.ps1 web        # Start Next.js on :3000
```

## Test

**Unix/Linux/macOS:**
```bash
# Backend tests
make test-api

# Frontend tests
make test-web

# All tests
make test
```

**Windows:**
```powershell
# Backend tests
.\dev.ps1 test-api

# Frontend tests
.\dev.ps1 test-web

# All tests
.\dev.ps1 test
```

## Project Structure

```
├── frontend/          # Next.js application
├── backend/           # Go REST API
├── db/               # Database migrations & seeds
├── scripts/          # Utility scripts
└── README.md
```

## Development Flow

1. Enter household info + address
2. Coverage check → display available technologies
3. Show top 3 recommended bundles with costs & savings
4. Select bundle → choose installation slot
5. Confirm order → display confirmation with order ID
