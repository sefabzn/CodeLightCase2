-- Seed data for Turkcell Ev+Mobil Paket Danışmanı
-- Loads sample data from CSV files

-- Load coverage data
INSERT INTO coverage (address_id, city, district, fiber, vdsl, fwa) VALUES
('A1001', 'Istanbul', 'Kadikoy', true, true, false),
('A1002', 'Istanbul', 'Besiktas', false, true, true),
('A1003', 'Ankara', 'Cankaya', true, true, false),
('A1004', 'Izmir', 'Konak', false, false, true),
('A1005', 'Istanbul', 'Sisli', true, true, true),
('A1006', 'Ankara', 'Kecioren', false, true, false),
('A1007', 'Izmir', 'Bornova', true, false, false),
('A1008', 'Istanbul', 'Bakirkoy', true, true, false);

-- Load mobile plans
INSERT INTO mobile_plans (plan_id, plan_name, quota_gb, quota_min, monthly_price, overage_gb, overage_min) VALUES
(1, 'Basic 5GB', 5.00, 300.00, 99.90, 5.00, 0.50),
(2, 'Standard 10GB', 10.00, 500.00, 149.90, 4.00, 0.50),
(3, 'Premium 20GB', 20.00, 1000.00, 199.90, 3.50, 0.40),
(4, 'Unlimited 50GB', 50.00, 2000.00, 299.90, 2.00, 0.30),
(5, 'Family 75GB', 75.00, 3000.00, 399.90, 1.50, 0.25);

-- Load home plans
INSERT INTO home_plans (home_id, name, tech, down_mbps, monthly_price, install_fee) VALUES
(1, 'Fiber 50Mbps', 'fiber', 50, 89.90, 0.00),
(2, 'Fiber 100Mbps', 'fiber', 100, 119.90, 0.00),
(3, 'Fiber 200Mbps', 'fiber', 200, 159.90, 0.00),
(4, 'VDSL 25Mbps', 'vdsl', 25, 69.90, 50.00),
(5, 'VDSL 50Mbps', 'vdsl', 50, 89.90, 50.00),
(6, 'FWA 30Mbps', 'fwa', 30, 79.90, 100.00),
(7, 'FWA 50Mbps', 'fwa', 50, 99.90, 100.00);

-- Load TV plans
INSERT INTO tv_plans (tv_id, name, hd_hours_included, monthly_price) VALUES
(1, 'Basic TV', 30.00, 39.90),
(2, 'Standard TV', 60.00, 59.90),
(3, 'Premium TV', 120.00, 89.90),
(4, 'Sports Package', 150.00, 119.90);

-- Load bundling rules
INSERT INTO bundling_rules (rule_id, rule_type, description, discount_percent, applies_to) VALUES
(1, 'line_discount', '2nd line 5% discount', 5.00, 'mobile'),
(2, 'line_discount', '3+ lines 10% discount', 10.00, 'mobile'),
(3, 'bundle_discount', 'Mobile + Home bundle discount', 10.00, 'total'),
(4, 'bundle_discount', 'Mobile + Home + TV triple bundle', 15.00, 'total');

-- Load users
INSERT INTO users (user_id, name, address_id, current_bundle_label) VALUES
(1, 'Ahmet Yilmaz', 'A1001', 'Basic Mobile + VDSL'),
(2, 'Fatma Demir', 'A1002', 'Premium Mobile + FWA'),
(3, 'Mehmet Ozkan', 'A1003', 'Standard Mobile Only'),
(4, 'Ayse Kaya', 'A1004', 'Family Mobile + TV'),
(5, 'Can Celik', 'A1005', 'Premium Triple Bundle');

-- Load household data
INSERT INTO household (id, user_id, line_id, expected_gb, expected_min, tv_hd_hours) VALUES
(1, 1, 'LINE001', 8.00, 450.00, 25.00),
(2, 1, 'LINE002', 3.00, 200.00, 0.00),
(3, 2, 'LINE003', 15.00, 800.00, 40.00),
(4, 3, 'LINE004', 12.00, 600.00, 0.00),
(5, 4, 'LINE005', 6.00, 300.00, 60.00),
(6, 4, 'LINE006', 4.00, 250.00, 30.00),
(7, 4, 'LINE007', 2.00, 100.00, 15.00),
(8, 5, 'LINE008', 25.00, 1200.00, 80.00),
(9, 5, 'LINE009', 18.00, 900.00, 50.00);

-- Load current services
INSERT INTO current_services (id, user_id, has_home, home_tech, home_speed, has_tv, mobile_plan_ids) VALUES
(1, 1, true, 'vdsl', 25, false, '[1,2]'),
(2, 2, true, 'fwa', 30, false, '[3]'),
(3, 3, false, null, null, false, '[2]'),
(4, 4, false, null, null, true, '[1,1,1]'),
(5, 5, true, 'fiber', 100, true, '[4,3]');

-- Load installation slots
INSERT INTO install_slots (slot_id, address_id, slot_start, slot_end, tech, available) VALUES
(1, 'A1001', '2024-01-15 09:00:00+03', '2024-01-15 12:00:00+03', 'fiber', true),
(2, 'A1001', '2024-01-15 13:00:00+03', '2024-01-15 16:00:00+03', 'fiber', true),
(3, 'A1001', '2024-01-16 09:00:00+03', '2024-01-16 12:00:00+03', 'vdsl', true),
(4, 'A1002', '2024-01-15 09:00:00+03', '2024-01-15 12:00:00+03', 'vdsl', true),
(5, 'A1002', '2024-01-15 14:00:00+03', '2024-01-15 17:00:00+03', 'fwa', true),
(6, 'A1003', '2024-01-16 09:00:00+03', '2024-01-16 12:00:00+03', 'fiber', false),
(7, 'A1003', '2024-01-16 13:00:00+03', '2024-01-16 16:00:00+03', 'fiber', true),
(8, 'A1004', '2024-01-17 09:00:00+03', '2024-01-17 12:00:00+03', 'fwa', true);

-- Update sequences to prevent ID conflicts
SELECT setval('users_user_id_seq', (SELECT MAX(user_id) FROM users));
SELECT setval('mobile_plans_plan_id_seq', (SELECT MAX(plan_id) FROM mobile_plans));
SELECT setval('home_plans_home_id_seq', (SELECT MAX(home_id) FROM home_plans));
SELECT setval('tv_plans_tv_id_seq', (SELECT MAX(tv_id) FROM tv_plans));
SELECT setval('bundling_rules_rule_id_seq', (SELECT MAX(rule_id) FROM bundling_rules));
SELECT setval('household_id_seq', (SELECT MAX(id) FROM household));
SELECT setval('current_services_id_seq', (SELECT MAX(id) FROM current_services));
SELECT setval('install_slots_slot_id_seq', (SELECT MAX(slot_id) FROM install_slots));
