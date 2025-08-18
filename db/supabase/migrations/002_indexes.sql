-- Indexes for Turkcell Ev+Mobil Paket Danışmanı
-- Adds performance indexes for frequently queried columns

-- Index on coverage table for address lookups
CREATE INDEX idx_coverage_address_id ON coverage(address_id);

-- Index on install_slots for address and tech filtering
CREATE INDEX idx_install_slots_address_id ON install_slots(address_id);
CREATE INDEX idx_install_slots_tech ON install_slots(tech);
CREATE INDEX idx_install_slots_available ON install_slots(available);
CREATE INDEX idx_install_slots_time ON install_slots(slot_start, slot_end);

-- Index on household for user lookups
CREATE INDEX idx_household_user_id ON household(user_id);

-- Index on current_services for user lookups
CREATE INDEX idx_current_services_user_id ON current_services(user_id);

-- Index on users for address lookups
CREATE INDEX idx_users_address_id ON users(address_id);

-- Composite indexes for common query patterns
CREATE INDEX idx_household_user_line ON household(user_id, line_id);
CREATE INDEX idx_install_slots_addr_tech_avail ON install_slots(address_id, tech, available);

-- Ensure unique constraints are properly indexed (most are already via PRIMARY KEY)
-- Add unique constraint on household per user+line combination
ALTER TABLE household ADD CONSTRAINT unique_user_line UNIQUE (user_id, line_id);
