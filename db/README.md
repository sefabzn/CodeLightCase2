# Turkcell Recommendation System - Database

Comprehensive database setup and management for the Turkcell package recommendation system using Supabase/PostgreSQL.

## 🚀 Quick Start

### Prerequisites
- [Supabase CLI](https://supabase.com/docs/guides/cli) installed
- Docker Desktop running (for local development)
- Git repository cloned

### Fresh Setup (< 2 commands)

```bash
# 1. Initialize and start Supabase
supabase start

# 2. Reset database with migrations and seed data
supabase db reset
```

That's it! Your database is now ready with all tables, indexes, and sample data loaded.

## 📊 Database Schema

### Core Tables

#### 🧑‍💼 Users
```sql
users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address_id VARCHAR(50) NOT NULL,
    current_bundle_label VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)
```

#### 📍 Coverage
```sql
coverage (
    address_id VARCHAR(50) PRIMARY KEY,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    fiber BOOLEAN DEFAULT FALSE,
    vdsl BOOLEAN DEFAULT FALSE,
    fwa BOOLEAN DEFAULT FALSE
)
```

#### 📱 Mobile Plans
```sql
mobile_plans (
    plan_id SERIAL PRIMARY KEY,
    plan_name VARCHAR(255) NOT NULL,
    quota_gb NUMERIC(10,2) NOT NULL,
    quota_min NUMERIC(10,2) NOT NULL,
    monthly_price NUMERIC(10,2) NOT NULL,
    overage_gb NUMERIC(10,2) DEFAULT 0,
    overage_min NUMERIC(10,2) DEFAULT 0
)
```

#### 🏠 Home Plans
```sql
home_plans (
    home_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tech VARCHAR(20) NOT NULL,
    down_mbps INTEGER NOT NULL,
    monthly_price NUMERIC(10,2) NOT NULL,
    install_fee NUMERIC(10,2) DEFAULT 0
)
```

#### 📺 TV Plans
```sql
tv_plans (
    tv_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    hd_hours_included NUMERIC(10,2) NOT NULL,
    monthly_price NUMERIC(10,2) NOT NULL
)
```

#### 🎯 Bundling Rules
```sql
bundling_rules (
    rule_id SERIAL PRIMARY KEY,
    rule_type VARCHAR(50) NOT NULL,
    condition_value VARCHAR(255),
    discount_percentage NUMERIC(5,2) DEFAULT 0,
    description TEXT
)
```

#### 👨‍👩‍👧‍👦 Household
```sql
household (
    user_id INTEGER REFERENCES users(user_id),
    line_id VARCHAR(50) NOT NULL,
    expected_gb NUMERIC(10,2) NOT NULL,
    expected_min NUMERIC(10,2) NOT NULL,
    tv_hd_hours NUMERIC(10,2) DEFAULT 0,
    PRIMARY KEY (user_id, line_id)
)
```

#### 📞 Current Services
```sql
current_services (
    user_id INTEGER REFERENCES users(user_id),
    service_type VARCHAR(50) NOT NULL,
    plan_id INTEGER,
    monthly_cost NUMERIC(10,2),
    PRIMARY KEY (user_id, service_type)
)
```

#### 📅 Install Slots
```sql
install_slots (
    slot_id VARCHAR(100) PRIMARY KEY,
    address_id VARCHAR(50) NOT NULL,
    slot_start TIMESTAMP WITH TIME ZONE NOT NULL,
    slot_end TIMESTAMP WITH TIME ZONE NOT NULL,
    tech VARCHAR(20) NOT NULL,
    available BOOLEAN DEFAULT TRUE
)
```

## 🗂️ Migration Files

### 001_init.sql
- Creates all core tables
- Sets up primary keys and constraints
- Defines data types and defaults

### 002_indexes.sql
- Performance indexes for frequent queries
- Coverage lookups by address_id
- Install slots by address_id and tech
- User and household relationships

### 003_seed.sql
- Sample data insertion
- Representative coverage scenarios
- Diverse plan offerings
- Test user data

## 🌱 Seed Data

### Sample Coverage Areas
```csv
address_id,city,district,fiber,vdsl,fwa
A1001,Istanbul,Kadikoy,1,1,1      # Full coverage
A1002,Ankara,Cankaya,0,1,1        # VDSL + FWA
A1003,Izmir,Bornova,1,0,1         # Fiber + FWA
A1004,Istanbul,Besiktas,1,1,0     # Fiber + VDSL
A1005,Antalya,Muratpasa,1,1,1     # Full coverage
```

### Mobile Plan Examples
```csv
plan_id,plan_name,quota_gb,quota_min,monthly_price,overage_gb,overage_min
1,Basic 5GB,5.00,300.00,99.90,5.00,0.50
2,Standard 10GB,10.00,500.00,149.90,4.00,0.50
3,Premium 20GB,20.00,1000.00,199.90,3.50,0.40
4,Unlimited 50GB,50.00,2000.00,299.90,2.00,0.30
```

### Home Internet Plans
```csv
home_id,name,tech,down_mbps,monthly_price,install_fee
201,Fiber Basic 50,fiber,50,120.00,50.00
202,Fiber Premium 100,fiber,100,150.00,100.00
203,Fiber Ultra 200,fiber,200,200.00,100.00
301,VDSL Standard 25,vdsl,25,90.00,75.00
401,FWA Basic 25,fwa,25,110.00,25.00
```

### TV Packages
```csv
tv_id,name,hd_hours_included,monthly_price
501,TV Basic,20.00,49.90
502,TV Standard,50.00,79.90
503,TV Premium,100.00,119.90
504,TV Ultimate,200.00,159.90
```

## 🛠️ Management Commands

### Database Reset
```bash
# Reset database with fresh migrations and seed data
supabase db reset

# Reset without seed data
supabase db reset --no-seed
```

### Migration Management
```bash
# Create new migration
supabase migration new my_new_migration

# Apply migrations
supabase migration up

# Check migration status
supabase migration list
```

### Data Import/Export
```bash
# Export schema
supabase db dump --schema-only > schema.sql

# Export data
supabase db dump --data-only > data.sql

# Import SQL file
psql -f your_file.sql "postgresql://postgres:postgres@127.0.0.1:54322/postgres"
```

## 🔍 Database Testing

### Smoke Test Queries
```sql
-- Verify all tables have data
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'coverage', COUNT(*) FROM coverage
UNION ALL
SELECT 'mobile_plans', COUNT(*) FROM mobile_plans
UNION ALL
SELECT 'home_plans', COUNT(*) FROM home_plans
UNION ALL
SELECT 'tv_plans', COUNT(*) FROM tv_plans
UNION ALL
SELECT 'install_slots', COUNT(*) FROM install_slots;

-- Test coverage query
SELECT address_id, city, district, 
       CASE WHEN fiber THEN 'fiber' END as fiber,
       CASE WHEN vdsl THEN 'vdsl' END as vdsl,
       CASE WHEN fwa THEN 'fwa' END as fwa
FROM coverage 
WHERE address_id = 'A1001';

-- Test plan availability
SELECT mp.plan_name, mp.quota_gb, mp.monthly_price
FROM mobile_plans mp
WHERE mp.quota_gb >= 10
ORDER BY mp.monthly_price;
```

### Sample Queries for Testing API
```sql
-- Get user with household data
SELECT u.*, h.line_id, h.expected_gb, h.expected_min, h.tv_hd_hours
FROM users u
LEFT JOIN household h ON u.user_id = h.user_id
WHERE u.user_id = 1;

-- Check coverage for address
SELECT * FROM coverage WHERE address_id = 'A1001';

-- Available install slots
SELECT * FROM install_slots 
WHERE address_id = 'A1001' 
  AND tech = 'fiber' 
  AND available = true
  AND slot_start > NOW()
ORDER BY slot_start;

-- Plan recommendations simulation
SELECT 
    mp.plan_name,
    mp.quota_gb,
    mp.monthly_price,
    hp.name as home_plan,
    hp.monthly_price as home_price,
    tp.name as tv_plan,
    tp.monthly_price as tv_price
FROM mobile_plans mp
CROSS JOIN home_plans hp
CROSS JOIN tv_plans tp
WHERE hp.tech IN (
    SELECT CASE WHEN fiber THEN 'fiber'
                WHEN vdsl THEN 'vdsl' 
                WHEN fwa THEN 'fwa' END
    FROM coverage WHERE address_id = 'A1001'
)
LIMIT 5;
```

## 🌐 Local Development

### Supabase Studio
Access the local Supabase Studio at: http://127.0.0.1:54323

**Features:**
- Visual table editor
- SQL query interface
- Real-time data viewing
- Authentication management
- Storage browser

### Database Connection
```bash
# Direct PostgreSQL connection
psql "postgresql://postgres:postgres@127.0.0.1:54322/postgres"

# Connection string for applications
DATABASE_URL="postgresql://postgres:postgres@127.0.0.1:54322/postgres"
```

### API Endpoints
```bash
# Supabase API URL
http://127.0.0.1:54321

# Direct database queries via REST API
curl "http://127.0.0.1:54321/rest/v1/coverage?address_id=eq.A1001" \
  -H "apikey: YOUR_ANON_KEY"
```

## 🚀 Production Setup

### Environment Variables
```bash
# Production Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key

# Connection string format
DATABASE_URL=postgresql://postgres:[password]@db.[project].supabase.co:5432/postgres
```

### Migration Deployment
```bash
# Link to production project
supabase link --project-ref your-project-ref

# Push migrations to production
supabase db push

# Verify deployment
supabase db remote commit
```

## 📈 Performance Tuning

### Indexes for Optimization
```sql
-- Coverage lookups
CREATE INDEX idx_coverage_address ON coverage(address_id);

-- Install slots by address and technology
CREATE INDEX idx_install_slots_address_tech ON install_slots(address_id, tech, available);

-- Plan queries
CREATE INDEX idx_mobile_plans_quota ON mobile_plans(quota_gb, monthly_price);
CREATE INDEX idx_home_plans_tech ON home_plans(tech, down_mbps);

-- User household relationships
CREATE INDEX idx_household_user ON household(user_id);
```

### Query Performance Tips
- Use address_id indexes for coverage lookups
- Filter install_slots by available = true
- Order plans by price for recommendation sorting
- Use EXPLAIN ANALYZE for slow queries

## 🔒 Security & Backup

### Row Level Security (RLS)
```sql
-- Enable RLS on sensitive tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE household ENABLE ROW LEVEL SECURITY;
ALTER TABLE current_services ENABLE ROW LEVEL SECURITY;

-- Example policies (customize for your auth)
CREATE POLICY "Users can view own data" ON users
  FOR SELECT USING (auth.uid()::text = user_id::text);
```

### Backup Strategy
```bash
# Create backup
supabase db dump > backup_$(date +%Y%m%d_%H%M%S).sql

# Automated backups (cron example)
0 2 * * * cd /path/to/project && supabase db dump > "backups/backup_$(date +\%Y\%m\%d).sql"
```

## 🐛 Troubleshooting

### Common Issues

#### "relation does not exist"
```bash
# Reset and apply migrations
supabase db reset
```

#### "connection refused"
```bash
# Ensure Docker is running and start Supabase
docker ps
supabase start
```

#### "no data in tables"
```bash
# Check if seed ran properly
supabase db reset
# or manually run seed
psql -f supabase/migrations/003_seed.sql "postgresql://postgres:postgres@127.0.0.1:54322/postgres"
```

#### "permission denied"
```bash
# Check RLS policies if enabled
SELECT * FROM pg_policies WHERE tablename = 'your_table';
```

### Logs and Debugging
```bash
# View Supabase logs
supabase logs

# PostgreSQL logs
docker logs supabase_db_*

# Check service status
supabase status
```

---

**Database Version**: PostgreSQL 17  
**Supabase CLI**: Latest  
**Last Updated**: December 2024  
**Contact**: Turkcell Database Team

## 📚 Additional Resources

- [Supabase Documentation](https://supabase.com/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [SQL Migration Best Practices](https://supabase.com/docs/guides/database/migrations)
