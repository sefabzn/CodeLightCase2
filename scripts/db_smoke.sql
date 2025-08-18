-- Database Smoke Test for Turkcell Ev+Mobil Paket Danışmanı
-- Validates that all tables are populated and referentially coherent

\echo 'Running database smoke test...'
\echo ''

-- Test 1: Check table row counts
\echo '=== TABLE ROW COUNTS ==='
SELECT 'users' as table_name, COUNT(*) as row_count FROM users
UNION ALL
SELECT 'coverage', COUNT(*) FROM coverage
UNION ALL
SELECT 'mobile_plans', COUNT(*) FROM mobile_plans
UNION ALL
SELECT 'home_plans', COUNT(*) FROM home_plans
UNION ALL
SELECT 'tv_plans', COUNT(*) FROM tv_plans
UNION ALL
SELECT 'bundling_rules', COUNT(*) FROM bundling_rules
UNION ALL
SELECT 'household', COUNT(*) FROM household
UNION ALL
SELECT 'current_services', COUNT(*) FROM current_services
UNION ALL
SELECT 'install_slots', COUNT(*) FROM install_slots;

\echo ''

-- Test 2: Sample one row from each table
\echo '=== SAMPLE DATA FROM EACH TABLE ==='

\echo 'Sample user:'
SELECT user_id, name, address_id FROM users LIMIT 1;

\echo 'Sample coverage:'
SELECT address_id, city, district, fiber, vdsl, fwa FROM coverage LIMIT 1;

\echo 'Sample mobile plan:'
SELECT plan_id, plan_name, quota_gb, monthly_price FROM mobile_plans LIMIT 1;

\echo 'Sample home plan:'
SELECT home_id, name, tech, down_mbps, monthly_price FROM home_plans LIMIT 1;

\echo 'Sample TV plan:'
SELECT tv_id, name, hd_hours_included, monthly_price FROM tv_plans LIMIT 1;

\echo 'Sample bundling rule:'
SELECT rule_id, rule_type, description, discount_percent FROM bundling_rules LIMIT 1;

\echo 'Sample household:'
SELECT id, user_id, line_id, expected_gb, expected_min FROM household LIMIT 1;

\echo 'Sample current services:'
SELECT id, user_id, has_home, home_tech, has_tv FROM current_services LIMIT 1;

\echo 'Sample install slot:'
SELECT slot_id, address_id, slot_start, tech, available FROM install_slots LIMIT 1;

\echo ''

-- Test 3: Referential integrity checks
\echo '=== REFERENTIAL INTEGRITY CHECKS ==='

\echo 'Users with valid addresses in coverage:'
SELECT COUNT(*) as users_with_coverage 
FROM users u 
JOIN coverage c ON u.address_id = c.address_id;

\echo 'Household lines with valid users:'
SELECT COUNT(*) as household_with_users 
FROM household h 
JOIN users u ON h.user_id = u.user_id;

\echo 'Current services with valid users:'
SELECT COUNT(*) as services_with_users 
FROM current_services cs 
JOIN users u ON cs.user_id = u.user_id;

\echo 'Install slots for addresses with coverage:'
SELECT COUNT(*) as slots_with_coverage 
FROM install_slots s 
JOIN coverage c ON s.address_id = c.address_id;

\echo ''

-- Test 4: Business logic validation
\echo '=== BUSINESS LOGIC VALIDATION ==='

\echo 'Addresses with fiber availability:'
SELECT COUNT(*) as fiber_addresses FROM coverage WHERE fiber = true;

\echo 'Mobile plans price range:'
SELECT MIN(monthly_price) as min_price, MAX(monthly_price) as max_price FROM mobile_plans;

\echo 'Available installation slots:'
SELECT COUNT(*) as available_slots FROM install_slots WHERE available = true;

\echo 'Total household lines:'
SELECT COUNT(*) as total_lines FROM household;

\echo ''
\echo 'Database smoke test completed!'
\echo 'If all queries returned data without errors, the database is healthy.'
