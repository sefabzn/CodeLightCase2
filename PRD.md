PRD – Turkcell Ev+Mobil Paket Danışmanı

Web Application with Golang Backend, Next.js Frontend, Supabase DB

1. System Overview

The web application helps users determine the most cost-effective combination of mobile, home internet, and TV packages based on:

Household profile (lines, usage, TV hours)

Address coverage (fiber, VDSL, FWA availability)

Catalog plans & bundling rules

Automated cost simulation (base price + overages + discounts)

Recommendation of top 3 lowest-cost bundles

Mock checkout with installation slot booking

2. High-Level Architecture

Frontend: Next.js 14 (App Router), TailwindCSS, React Query (for data fetching & caching).

Backend: Golang REST API (Echo/Fiber framework).

Database: Supabase (Postgres + auth).

Data source: Seed data from CSV (coverage, plans, rules, household, etc.) imported into Supabase.

State Management:

Local state (wizard inputs, UI) in React hooks/context.

Remote state (user, catalog, recommendations, checkout) in Supabase + React Query cache.

Integration:

Next.js → calls Golang API.

Golang API → reads/writes Supabase via pgx driver or Supabase client.

3. File & Folder Structure
3.1 Frontend (Next.js)
/frontend
  /app
    /layout.tsx              # Root layout (theme, header, footer)
    /page.tsx                # Entry screen (wizard)
    /recommendations/page.tsx # Show top 3 bundles
    /checkout/page.tsx       # Slot selection + confirmation
  /components
    HouseholdForm.tsx        # Form for household input
    AddressForm.tsx          # Address + coverage selector
    RecommendationCard.tsx   # Each recommended bundle
    SlotPicker.tsx           # Installation slot selector
    SummaryModal.tsx         # Bundle details modal
  /lib
    api.ts                   # API client (fetch wrapper with React Query)
    types.ts                 # Shared TypeScript types
  /context
    WizardContext.tsx        # Stores multi-step wizard state
  /styles
    globals.css
  /utils
    formatters.ts            # currency, labels, etc.
  package.json
  tsconfig.json

3.2 Backend (Golang)
/backend
  /cmd
    /server/main.go          # Entry point
  /internal
    /api
      handlers.go            # Route handlers (users, catalog, recommendation, checkout)
      middleware.go
    /services
      coverage.go            # Coverage matching logic
      recommendation.go      # Core cost + discount simulation
      checkout.go            # Mock order placement
    /models
      user.go
      plans.go
      recommendation.go
    /db
      supabase.go            # Connection pool to Supabase/Postgres
      seed.go                # Data seeding from CSV (if needed)
    /utils
      cost_calculator.go     # Shared functions for overages & discounts
      validator.go           # Input validation
  go.mod
  go.sum

3.3 Database (Supabase / Postgres)

Tables (mirroring dataset in PRD):

users(user_id, name, address_id, current_bundle_label)

coverage(address_id, city, district, fiber, vdsl, fwa)

mobile_plans(plan_id, plan_name, quota_gb, quota_min, monthly_price, overage_gb, overage_min)

home_plans(home_id, name, tech, down_mbps, monthly_price, install_fee)

tv_plans(tv_id, name, hd_hours_included, monthly_price)

bundling_rules(rule_id, rule_type, description, discount_percent, applies_to)

household(user_id, line_id, expected_gb, expected_min, tv_hd_hours)

current_services(user_id, has_home, home_tech, home_speed, has_tv, mobile_plan_ids)

install_slots(slot_id, address_id, slot_start, slot_end, tech)

4. API Design
Endpoint	Method	Description
/api/users/{id}	GET	Fetch user + household profile
/api/catalog	GET	Fetch coverage, plans, bundling rules, slots
/api/recommendation	POST	Calculate & return top 3 bundles
/api/checkout	POST	Mock booking of installation slot + order creation
5. Data Flow

Household & Address Input (Frontend)

User enters household members, expected GB/min, TV hours.

Selects address → backend queries coverage.

Catalog Fetch (Frontend → Backend → Supabase)

On first load, frontend fetches catalog data.

Recommendation Calculation (Backend)

Core logic in recommendation.go:

Match household needs → best fit mobile/home/tv plans.

Apply overage cost.

Apply line discounts + bundle discounts.

Sort by lowest total, return top 3.

Recommendation Display (Frontend)

Show 3 cards with total cost, savings, breakdown.

“View Details” expands modal with full plan details.

Checkout Flow

User selects bundle → picks available slot.

POST /api/checkout → stores mock order in Supabase.

Returns confirmation ID.

6. State Management

Frontend Local State:

Household + address input lives in WizardContext.

Temp UI state (modal open, form inputs) in component hooks.

Frontend Remote State:

React Query caches users, catalog, recommendation, checkout.

Backend State:

Stateless REST API; Supabase DB is source of truth.

Database State:

Stores catalog, rules, users, slots, household.

7. Services & Connections

Next.js → Golang API: via REST fetch.

Golang API → Supabase: via pgx or Supabase Go client.

Supabase → Frontend (optional): direct Supabase SDK for auth/user management if needed.

8. Non-Functional Requirements

Performance: Recommendation calc < 500ms for household ≤ 5 lines.

Scalability: Backend stateless; can run behind load balancer.

Reliability: Mock checkout always succeeds (for MVP).

Security: JWT-based auth (Supabase Auth).

Testing: Unit tests for cost calculation, integration tests for API.

9. Demo Flow (Expected)

Enter household info + address.

Coverage check → “Fiber available” badge.

Show top 3 recommended bundles (with savings).

Select bundle → choose installation slot.

Confirm order → mock success page with order_id.

✅ This PRD sets up a clear Golang + Next.js + Supabase implementation path, including architecture, state, services, and file structures.