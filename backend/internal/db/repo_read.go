package db

import (
	"context"
	"fmt"

	"app/internal/models"

	"github.com/jackc/pgx/v5"
)

// GetUser retrieves a user by ID with their address information
func (db *DB) GetUser(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT user_id, name, address_id, current_bundle_label, created_at
		FROM users 
		WHERE user_id = $1
	`

	var user models.User
	err := db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Name,
		&user.AddressID,
		&user.CurrentBundleLabel,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetHousehold retrieves all household members for a user
func (db *DB) GetHousehold(ctx context.Context, userID int) ([]models.Household, error) {
	query := `
		SELECT id, user_id, line_id, expected_gb, expected_min, tv_hd_hours
		FROM household 
		WHERE user_id = $1
		ORDER BY line_id
	`

	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query household: %w", err)
	}
	defer rows.Close()

	var household []models.Household
	for rows.Next() {
		var h models.Household
		err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.LineID,
			&h.ExpectedGB,
			&h.ExpectedMin,
			&h.TVHDHours,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan household row: %w", err)
		}
		household = append(household, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate household rows: %w", err)
	}

	return household, nil
}

// GetCoverage retrieves coverage information for an address
func (db *DB) GetCoverage(ctx context.Context, addressID string) (*models.Coverage, error) {
	query := `
		SELECT address_id, city, district, fiber, vdsl, fwa
		FROM coverage 
		WHERE address_id = $1
	`

	var coverage models.Coverage
	err := db.Pool.QueryRow(ctx, query, addressID).Scan(
		&coverage.AddressID,
		&coverage.City,
		&coverage.District,
		&coverage.Fiber,
		&coverage.VDSL,
		&coverage.FWA,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("coverage for address %s not found", addressID)
		}
		return nil, fmt.Errorf("failed to get coverage: %w", err)
	}

	return &coverage, nil
}

// CatalogData represents all catalog information
type CatalogData struct {
	Coverage      []models.Coverage     `json:"coverage"`
	MobilePlans   []models.MobilePlan   `json:"mobile_plans"`
	HomePlans     []models.HomePlan     `json:"home_plans"`
	TVPlans       []models.TVPlan       `json:"tv_plans"`
	BundlingRules []models.BundlingRule `json:"bundling_rules"`
	InstallSlots  []models.InstallSlot  `json:"install_slots"`
}

// GetCatalog retrieves all catalog data (plans, rules, coverage)
func (db *DB) GetCatalog(ctx context.Context) (*CatalogData, error) {
	catalog := &CatalogData{}

	// Get coverage data
	coverageQuery := `SELECT address_id, city, district, fiber, vdsl, fwa FROM coverage ORDER BY city, district`
	rows, err := db.Pool.Query(ctx, coverageQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query coverage: %w", err)
	}

	for rows.Next() {
		var c models.Coverage
		err := rows.Scan(&c.AddressID, &c.City, &c.District, &c.Fiber, &c.VDSL, &c.FWA)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan coverage: %w", err)
		}
		catalog.Coverage = append(catalog.Coverage, c)
	}
	rows.Close()

	// Get mobile plans
	mobileQuery := `SELECT plan_id, plan_name, quota_gb, quota_min, monthly_price, overage_gb, overage_min FROM mobile_plans ORDER BY monthly_price`
	rows, err = db.Pool.Query(ctx, mobileQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query mobile plans: %w", err)
	}

	for rows.Next() {
		var mp models.MobilePlan
		err := rows.Scan(&mp.PlanID, &mp.PlanName, &mp.QuotaGB, &mp.QuotaMin, &mp.MonthlyPrice, &mp.OverageGB, &mp.OverageMin)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan mobile plan: %w", err)
		}
		catalog.MobilePlans = append(catalog.MobilePlans, mp)
	}
	rows.Close()

	// Get home plans
	homeQuery := `SELECT home_id, name, tech, down_mbps, monthly_price, install_fee FROM home_plans ORDER BY tech, monthly_price`
	rows, err = db.Pool.Query(ctx, homeQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query home plans: %w", err)
	}

	for rows.Next() {
		var hp models.HomePlan
		err := rows.Scan(&hp.HomeID, &hp.Name, &hp.Tech, &hp.DownMbps, &hp.MonthlyPrice, &hp.InstallFee)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan home plan: %w", err)
		}
		catalog.HomePlans = append(catalog.HomePlans, hp)
	}
	rows.Close()

	// Get TV plans
	tvQuery := `SELECT tv_id, name, hd_hours_included, monthly_price FROM tv_plans ORDER BY monthly_price`
	rows, err = db.Pool.Query(ctx, tvQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query TV plans: %w", err)
	}

	for rows.Next() {
		var tp models.TVPlan
		err := rows.Scan(&tp.TVID, &tp.Name, &tp.HDHoursIncluded, &tp.MonthlyPrice)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan TV plan: %w", err)
		}
		catalog.TVPlans = append(catalog.TVPlans, tp)
	}
	rows.Close()

	// Get bundling rules
	rulesQuery := `SELECT rule_id, rule_type, description, discount_percent, applies_to FROM bundling_rules ORDER BY rule_type, discount_percent`
	rows, err = db.Pool.Query(ctx, rulesQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query bundling rules: %w", err)
	}

	for rows.Next() {
		var br models.BundlingRule
		err := rows.Scan(&br.RuleID, &br.RuleType, &br.Description, &br.DiscountPercent, &br.AppliesTo)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("failed to scan bundling rule: %w", err)
		}
		catalog.BundlingRules = append(catalog.BundlingRules, br)
	}
	rows.Close()

	return catalog, nil
}

// GetInstallSlots retrieves available installation slots for an address and technology
func (db *DB) GetInstallSlots(ctx context.Context, addressID string, tech string) ([]models.InstallSlot, error) {
	query := `
		SELECT slot_id, address_id, slot_start, slot_end, tech, available
		FROM install_slots 
		WHERE address_id = $1 AND tech = $2 AND available = true
		ORDER BY slot_start
	`

	rows, err := db.Pool.Query(ctx, query, addressID, tech)
	if err != nil {
		return nil, fmt.Errorf("failed to query install slots: %w", err)
	}
	defer rows.Close()

	var slots []models.InstallSlot
	for rows.Next() {
		var slot models.InstallSlot
		err := rows.Scan(
			&slot.SlotID,
			&slot.AddressID,
			&slot.SlotStart,
			&slot.SlotEnd,
			&slot.Tech,
			&slot.Available,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan install slot: %w", err)
		}
		slots = append(slots, slot)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate install slot rows: %w", err)
	}

	return slots, nil
}
