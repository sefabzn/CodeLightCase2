# Turkcell Ev+Mobil Paket Danışmanı - Development Commands

.PHONY: help db-up db-down api web test-api test-web test seed clean

# Default target
help:
	@echo "Available targets:"
	@echo "  db-up      - Start Supabase local development"
	@echo "  db-down    - Stop Supabase local development"
	@echo "  api        - Start Go backend server"
	@echo "  web        - Start Next.js frontend"
	@echo "  test-api   - Run backend tests"
	@echo "  test-web   - Run frontend tests"
	@echo "  test       - Run all tests"
	@echo "  seed       - Load sample data into database"
	@echo "  clean      - Clean build artifacts"

# Database
db-up:
	cd db && npx supabase start

db-down:
	cd db && npx supabase stop

# Backend
api:
	cd backend && go run ./cmd/server

# Frontend
web:
	cd frontend && npm run dev

# Testing
test-api:
	cd backend && go test ./...

test-web:
	cd frontend && npm test

test: test-api test-web

# Database seeding
seed:
	cd db && npx supabase db reset

# Cleanup
clean:
	cd backend && go clean
	cd frontend && npm run build || true
	cd frontend && rm -rf .next || true
