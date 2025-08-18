MVP Build Plan — Turkcell Ev+Mobil Paket Danışmanı

(Backend: Go • Frontend: Next.js • DB: Supabase/Postgres)

Each task below is tiny, testable, and has a clear start/end with acceptance criteria.
Execute tasks in order (top → bottom). You can pause after any task and manually verify.

0) Conventions & Repo Layout

Repo: monorepo with frontend/, backend/, db/, scripts/.

Package managers: pnpm (frontend), go (backend).

Run targets: make or npm scripts equivalents.

Testing: Go testing & httptest; frontend vitest + @testing-library/react.

A) Bootstrap & Tooling
A1. Create monorepo skeleton

Start: Empty directory
End: frontend/, backend/, db/, scripts/, root README
Steps

Create folders and .gitignore for Node, Go, macOS.

Add root README.md with quick-start.
Test / DoD

ls shows expected folders.

README renders with sections: Setup, Run, Test.

A2. Initialize Supabase project (local or cloud)

Start: Supabase CLI installed
End: Supabase project created; .env.local placeholders
Steps

supabase init in /db.

Add .env variables in frontend/.env.local and backend/.env (SUPABASE_URL, SUPABASE_ANON_KEY, DATABASE_URL).
Test / DoD

supabase start works; supabase status is healthy.

A3. Initialize Go backend module

Start: /backend empty
End: go.mod, minimal main
Steps

go mod init app

Create /cmd/server/main.go with “hello” HTTP server on :8000.
Test / DoD

go run ./cmd/server starts.

curl localhost:8000/health returns 200 OK and {"status":"ok"} (add stub).

A4. Initialize Next.js frontend (App Router + TS)

Start: /frontend empty
End: Next.js app bootstrapped
Steps

pnpm create next-app@latest → frontend

Enable app/ router, TypeScript, ESLint, Tailwind.
Test / DoD

pnpm dev shows Next.js starter at /.

A5. Add Makefile (root)

Start: none
End: Makefile with common targets
Steps

Targets: db-up, db-down, api, web, test-api, test-web, seed.
Test / DoD

make web runs Next.

make api runs Go server.

B) Database (Schema & Seed)
B1. Define SQL schema

Start: /db initialized
End: db/migrations/001_init.sql with all tables
Steps

Create tables: users, coverage, mobile_plans, home_plans, tv_plans, bundling_rules, household, current_services, install_slots.

Proper types (boolean for tech flags, numeric(10,2) for prices, timestamps for slots).
Test / DoD

supabase db reset applies migration without errors.

Tables exist via supabase db dump or psql \dt.

B2. Add minimal indexes

Start: schema exists
End: db/migrations/002_indexes.sql
Steps

Index on coverage(address_id), install_slots(address_id), PK/unique constraints on ids.
Test / DoD

Migration applies cleanly (supabase db reset).

B3. Create CSV seed files

Start: none
End: /db/seed/*.csv with provided sample rows
Steps

Add CSVs: coverage.csv, mobile_plans.csv, home_plans.csv, tv_plans.csv, bundling_rules.csv, users.csv, household.csv, current_services.csv, install_slots.csv.
Test / DoD

Files exist and open; header rows present.

B4. Write SQL seed script

Start: CSVs present
End: db/migrations/003_seed.sql
Steps

Use COPY ... FROM PROGRAM 'cat /.../file.csv' CSV HEADER (or Supabase import).
Test / DoD

supabase db reset loads data; simple SELECT count(*) returns >0 for each table.

B5. Add DB smoke script

Start: seed loaded
End: /scripts/db_smoke.sql
Steps

Add queries: 1 row from each table; ensure referential coherence (e.g., install_slots for address_id).
Test / DoD

psql -f scripts/db_smoke.sql prints expected rows.

C) Backend — Foundations
C1. HTTP framework & router

Start: hello server
End: Echo/Fiber app with base middlewares
Steps

Add router, JSON middleware, recover, CORS (frontend origin).

Add /health handler.
Test / DoD

GET /health returns {status:"ok"} with CORS headers.

C2. Config loader

Start: none
End: /internal/utils/config.go
Steps

Load env vars (PORT, DATABASE_URL, Supabase URL/keys).
Test / DoD

Misconfigured env returns clear error on boot; with valid env, boots.

C3. PG connection pool

Start: none
End: /internal/db/supabase.go using pgxpool
Steps

Create pool, ping on startup; graceful shutdown.
Test / DoD

Startup logs “DB connected”; SIGINT closes pool.

C4. Domain models (structs)

Start: none
End: /internal/models/*.go
Steps

Define Go structs mirroring SQL tables with JSON tags.
Test / DoD

go vet & go build succeed.

C5. Repository functions (read-only)

Start: none
End: /internal/db/repo_read.go
Steps

Functions: GetUser(id), GetHousehold(userID), GetCoverage(addressID), GetCatalog(), GetInstallSlots(addressID, tech).
Test / DoD

Each function has a unit test with pgxpool to local DB and fixtures, returning non-empty rows.

D) Backend — Cost Engine (Unit-Test First)

Implement with table-driven tests; keep each function pure.

D1. Overage for a single mobile line

Start: none
End: utils/cost_calculator.go::CalcMobileOverage(lineUsage, plan)
Test / DoD

Given usage > quota, returns over_gb*overage_gb + over_min*overage_min.

4 test cases: under/over GB, under/over mins.

D2. Mobile line total

Start: D1 done
End: CalcMobileLineCost(usage, plan)
Test / DoD

Returns monthly_price + overage.

3 tests: exact quota, over GB, over mins.

D3. Aggregate mobile cost for household

Start: D2 done
End: CalcMobileTotal(lines[])
Test / DoD

Sums costs for multiple lines; 2 tests (1 line, 3 lines).

D4. Extra line discount

Start: D3 done
End: ApplyExtraLineDiscount(mobileTotal, lineCount)
Rules: 2nd line -5%, 3+ lines -10% on mobile component only
Test / DoD

Cases: 1 line (0%), 2 lines (5%), 3 lines (10%).

D5. Home plan selector by tech/speed

Start: none
End: SelectHomePlan(techAvailable[], neededMbps, homePlans)
Rules: choose lowest plan with down_mbps >= needed; fallback to next tech in priority
Test / DoD

3 tests: fiber available, fiber missing → vdsl, none → fwa.

D6. TV plan selector by HD hours

Start: none
End: SelectTvPlan(hours, tvPlans)
Test / DoD

Picks minimal plan covering hours; 3 tests (20/60/120 hours).

D7. Bundle discount

Start: none
End: CalcBundleDiscount(hasMobile, hasHome, hasTv)
Rules: mobile+home=10%; mobile+home+tv=15%; else 0
Test / DoD

3 tests: pair, triple, none.

D8. Assemble combo total

Start: D1–D7 done
End: CalcGrandTotal(mobileAfterLineDisc, home, tv, bundleDiscount)
Test / DoD

Calculates final total with discount; 2 tests.

D9. Reasoning generator

Start: none
End: BuildReasoning(struct)
Test / DoD

Returns concise string with chosen plans and discounts applied.

E) Backend — Recommendation Service
E1. Input DTOs & validation

Start: none
End: /internal/api/dto.go, /internal/utils/validator.go
Steps

Define DTOs for POST /api/recommendation.

Validate presence/types for user_id, address_id, household[].
Test / DoD

Invalid payload returns 400 with error list; valid passes.

E2. Coverage check service

Start: repos ready
End: /internal/services/coverage.go::ComputeCoverage(addressID)
Steps

Read coverage row; produce []tech ordered by preference (fiber, vdsl, fwa).
Test / DoD

Known address returns expected tech array; unknown returns empty.

E3. Compute needed home Mbps

Start: none
End: EstimateNeededMbps(household)
Rule (simple): map total HD hours → Mbps tiers (e.g., 50/100/200).
Test / DoD

3 cases yield expected Mbps.

E4. Generate candidate combinations

Start: selectors ready
End: BuildCandidates(household, techs, plans)
Steps

Produce combinations: {mobile-only}, {mobile+home}, {mobile+home+tv}.
Test / DoD

Returns non-empty candidate set for typical inputs.

E5. Line-to-plan matching heuristic

Start: none
End: MatchMobilePlansForLines(lines, mobilePlans)
Rule: pick minimal plan meeting expected GB/min; allow overage if cheaper (compare total).
Test / DoD

Lines mapped to plan IDs deterministically.

E6. Price each candidate

Start: cost funcs ready
End: PriceCandidate(candidate)
Steps

Apply line totals, extra line discount, home/tv price, then bundle discount.
Test / DoD

Numeric totals match expected for fixture household.

E7. Sort & top3

Start: pricing ready
End: SelectTop3(candidates)
Test / DoD

Sorted ascending by monthly_total; returns ≤3 items.

F) Backend — HTTP Endpoints
F1. GET /api/users/{id}

Start: repo funcs exist
End: Handler returns { user, household, address_id }
Test / DoD

httptest returns 200 and valid JSON for fixture user_id.

F2. GET /api/catalog

Start: repo funcs exist
End: Return { coverage, mobile_plans, home_plans, tv_plans, bundling_rules, install_slots }
Test / DoD

httptest returns all arrays with length > 0.

F3. POST /api/recommendation

Start: service ready
End: Accept DTO → return { top3: [...] }
Test / DoD

Valid sample payload returns 1–3 combos with monthly_total, savings, reasoning.

Invalid payload → 400.

F4. POST /api/checkout (mock)

Start: none
End: Accept selection+slot → return {status:"ok",order_id}
Test / DoD

Persists a lightweight orders record (optional table) or logs; returns 201.

F5. CORS & error mapping

Start: handlers exist
End: Uniform error JSON structure
Test / DoD

404/500 return {error:{code,message}} with CORS headers.

G) Backend — Sample Data & Golden Tests
G1. Golden JSON for recommendation

Start: endpoint works
End: /backend/testdata/reco_request.json, reco_response.json
Test / DoD

go test compares response to golden file (allowing totals tolerance if needed).

G2. CLI curl script

Start: endpoint works
End: /scripts/curl_reco.sh
Test / DoD

Running the script prints top3 JSON.

H) Frontend — Foundations
H1. Tailwind, fonts, theme

Start: Next created
End: Global styles + light/dark toggle
Test / DoD

Toggle switches class and persists via localStorage.

H2. Type definitions

Start: none
End: frontend/lib/types.ts mirrors backend DTOs
Test / DoD

TS builds; types imported across pages/components.

H3. API client & React Query

Start: none
End: frontend/lib/api.ts with fetcher and hooks
Test / DoD

useQuery('health') demo returns ok.

H4. Wizard context

Start: none
End: frontend/context/WizardContext.tsx
State: userId, addressId, household[], preferTech[]
Test / DoD

Provider wraps app; useWizard() returns defaults; unit test passes.

I) Frontend — Screens & Components
I1. Household form

Start: none
End: components/HouseholdForm.tsx
Fields: member lines (dynamic rows), expected_gb, expected_min, tv_hd_hours.
Test / DoD

Add/remove line works; validation errors show; unit test with @testing-library/react.

I2. Address form

Start: none
End: components/AddressForm.tsx
Fields: city, district (free text), address_id (select or input).
Test / DoD

On address_id blur → calls GET /api/catalog (or separate coverage lookup) and shows badges.

I3. Coverage badge

Start: none
End: components/CoverageBadge.tsx
Test / DoD

Given coverage flags, renders fiber/vdsl/fwa badges.

I4. Wizard page

Start: none
End: app/page.tsx renders Household + Address forms, “Continue”
Test / DoD

Clicking Continue stores data in context and navigates to /recommendations.

I5. Recommendation card

Start: none
End: components/RecommendationCard.tsx
Props: combo_label, items, monthly_total, savings
Test / DoD

Renders cost, tags; “Details” opens modal.

I6. Recommendations page

Start: none
End: app/recommendations/page.tsx
Steps

Collect state from context; call POST /api/recommendation; render top3.
Test / DoD

With known inputs, shows 1–3 cards; handles empty/error.

I7. Summary modal

Start: none
End: components/SummaryModal.tsx
Test / DoD

Displays breakdown (mobile/home/tv, discounts, reasoning); close works.

I8. Slot picker

Start: none
End: components/SlotPicker.tsx
Props: addressId, tech
Test / DoD

Lists slots; selecting one emits onChange(slotId).

I9. Checkout page

Start: none
End: app/checkout/page.tsx
Steps

Receives selection via query or context; shows SlotPicker; “Confirm” posts to /api/checkout.
Test / DoD

On success, shows order id; on error, shows retry.

J) Frontend — Data Hooks (React Query)
J1. useCatalog(addressId)

Start: none
End: Hook returns coverage/plans/slots
Test / DoD

Suspense state & error states verified; cache key includes addressId.

J2. useRecommendation(input)

Start: none
End: Hook posts to API and returns top3
Test / DoD

Calling with fixture input returns array length 3.

J3. useCheckout(selection, slotId)

Start: none
End: Hook posts checkout
Test / DoD

Returns {order_id}; disabled until both args present.

K) End-to-End Manual Test (Scripted)
K1. Happy-path smoke doc

Start: app running
End: scripts/manual_e2e.md
Steps

Enter 3 lines usage → address A1001 → see fiber badge → get top3 → pick first → choose slot → confirm → see MOCK-BUNDLE-001.
Test / DoD

Follow steps exactly and reach confirmation page.

L) QA & Error Cases
L1. No coverage path

Start: coverage service live
End: If no fiber, show vdsl/fwa recommendation fallback
Test / DoD

Use A1002 or A1003 to validate tech fallback.

L2. Single-line household

Start: recommendation live
End: Works with 1 line, no extra line discount
Test / DoD

Verify totals do not apply multi-line discount.

L3. Max discount path (3+ lines + triple bundle)

Start: recommendation live
End: Applies 10% line discount + 15% bundle
Test / DoD

Inspect reasoning text includes both discounts.

M) Docs & DX
M1. API README

Start: endpoints done
End: /backend/README.md with endpoint docs + curl examples
Test / DoD

Copy-paste curl returns valid JSON.

M2. Frontend README

Start: UI done
End: /frontend/README.md with environment, run, tests
Test / DoD

Following steps brings UI up locally.

M3. Seeds & reset docs

Start: seeds done
End: /db/README.md detailing supabase db reset and seed expectations
Test / DoD

Fresh clone can reproduce DB in < 1–2 commands.

N) Optional (Nice-to-Have, Still Tiny)
N1. Orders table for checkout

Start: none
End: Add orders table + insert on checkout
Test / DoD

New row after POST; GET query shows it.

N2. Basic ILP/knapsack toggle (mock)

Start: heuristic exists
End: Flag ?opt=ilp triggers alternative matcher
Test / DoD

Responses differ in edge cases (documented).

N3. Unit price normalization helper

Start: none
End: Function to normalize TL/other (mock)
Test / DoD

Given FX rate param, returns adjusted totals (pure fn tests).

O) Minimal Test Matrix (Execute After Each Area)

DB: migrations apply; seed counts > 0.

Backend unit: go test ./... green; coverage for cost engine > 80% (small suite).

Backend API: httptest passes for 4 endpoints.

Frontend unit: pnpm test green; form add/remove row; hooks handle success/error.

E2E manual: happy path + no-coverage path.
Q) Acceptance Criteria (MVP Complete)

User can input household + address on /.

Coverage badge renders correctly (fiber/vdsl/fwa).

/recommendations displays up to 3 combos sorted by lowest monthly total.

“Details” shows cost breakdown + discounts + reasoning text.

“Select & Continue” leads to /checkout.

Slot selection lists address-specific slots and allows choosing one.

Checkout returns mock order_id and displays confirmation.

All unit tests pass; curl scripts work from a fresh clone after DB reset.