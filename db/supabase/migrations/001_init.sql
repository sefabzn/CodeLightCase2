-- Initial schema for Turkcell Ev+Mobil Paket Danışmanı
-- Creates all tables for the package recommendation system

-- Users table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address_id VARCHAR(50) NOT NULL,
    current_bundle_label VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Coverage table - technology availability by address
CREATE TABLE coverage (
    address_id VARCHAR(50) PRIMARY KEY,
    city VARCHAR(100) NOT NULL,
    district VARCHAR(100) NOT NULL,
    fiber BOOLEAN DEFAULT FALSE,
    vdsl BOOLEAN DEFAULT FALSE,
    fwa BOOLEAN DEFAULT FALSE
);

-- Mobile plans catalog
CREATE TABLE mobile_plans (
    plan_id SERIAL PRIMARY KEY,
    plan_name VARCHAR(255) NOT NULL,
    quota_gb NUMERIC(10,2) NOT NULL,
    quota_min NUMERIC(10,2) NOT NULL,
    monthly_price NUMERIC(10,2) NOT NULL,
    overage_gb NUMERIC(10,2) DEFAULT 0,
    overage_min NUMERIC(10,2) DEFAULT 0
);

-- Home internet plans catalog
CREATE TABLE home_plans (
    home_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tech VARCHAR(50) NOT NULL, -- 'fiber', 'vdsl', 'fwa'
    down_mbps INTEGER NOT NULL,
    monthly_price NUMERIC(10,2) NOT NULL,
    install_fee NUMERIC(10,2) DEFAULT 0
);

-- TV plans catalog
CREATE TABLE tv_plans (
    tv_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    hd_hours_included NUMERIC(10,2) NOT NULL,
    monthly_price NUMERIC(10,2) NOT NULL
);

-- Bundling rules for discounts
CREATE TABLE bundling_rules (
    rule_id SERIAL PRIMARY KEY,
    rule_type VARCHAR(50) NOT NULL, -- 'line_discount', 'bundle_discount'
    description TEXT NOT NULL,
    discount_percent NUMERIC(5,2) NOT NULL,
    applies_to VARCHAR(100) NOT NULL -- 'mobile', 'home', 'tv', 'total'
);

-- Household members and their usage patterns
CREATE TABLE household (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    line_id VARCHAR(50) NOT NULL, -- identifier for this line/member
    expected_gb NUMERIC(10,2) NOT NULL,
    expected_min NUMERIC(10,2) NOT NULL,
    tv_hd_hours NUMERIC(10,2) DEFAULT 0
);

-- Current services the user has (for comparison)
CREATE TABLE current_services (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    has_home BOOLEAN DEFAULT FALSE,
    home_tech VARCHAR(50), -- 'fiber', 'vdsl', 'fwa'
    home_speed INTEGER,
    has_tv BOOLEAN DEFAULT FALSE,
    mobile_plan_ids TEXT -- JSON array of current mobile plan IDs
);

-- Installation slots for booking
CREATE TABLE install_slots (
    slot_id SERIAL PRIMARY KEY,
    address_id VARCHAR(50) NOT NULL,
    slot_start TIMESTAMP WITH TIME ZONE NOT NULL,
    slot_end TIMESTAMP WITH TIME ZONE NOT NULL,
    tech VARCHAR(50) NOT NULL, -- 'fiber', 'vdsl', 'fwa'
    available BOOLEAN DEFAULT TRUE
);

-- Add some basic constraints
ALTER TABLE household ADD CONSTRAINT positive_usage CHECK (expected_gb >= 0 AND expected_min >= 0 AND tv_hd_hours >= 0);
ALTER TABLE mobile_plans ADD CONSTRAINT positive_mobile_prices CHECK (monthly_price >= 0 AND overage_gb >= 0 AND overage_min >= 0);
ALTER TABLE home_plans ADD CONSTRAINT positive_home_prices CHECK (monthly_price >= 0 AND install_fee >= 0);
ALTER TABLE tv_plans ADD CONSTRAINT positive_tv_prices CHECK (monthly_price >= 0);
ALTER TABLE home_plans ADD CONSTRAINT valid_tech CHECK (tech IN ('fiber', 'vdsl', 'fwa'));
ALTER TABLE install_slots ADD CONSTRAINT valid_slot_tech CHECK (tech IN ('fiber', 'vdsl', 'fwa'));
ALTER TABLE install_slots ADD CONSTRAINT valid_time_slot CHECK (slot_end > slot_start);
