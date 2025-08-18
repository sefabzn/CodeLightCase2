package services

import (
	"context"
	"fmt"

	"app/internal/db"
)

// CoverageService handles coverage-related operations
type CoverageService struct {
	db db.DatabaseInterface
}

// NewCoverageService creates a new coverage service
func NewCoverageService(database db.DatabaseInterface) *CoverageService {
	return &CoverageService{
		db: database,
	}
}

// ComputeCoverage determines available technologies for an address
// Returns technologies ordered by preference: fiber > vdsl > fwa
func (s *CoverageService) ComputeCoverage(ctx context.Context, addressID string) ([]string, error) {
	coverage, err := s.db.GetCoverage(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get coverage for address %s: %w", addressID, err)
	}

	var availableTech []string

	// Add technologies in preference order
	if coverage.Fiber {
		availableTech = append(availableTech, "fiber")
	}

	if coverage.VDSL {
		availableTech = append(availableTech, "vdsl")
	}

	if coverage.FWA {
		availableTech = append(availableTech, "fwa")
	}

	return availableTech, nil
}

// GetCoverageInfo returns detailed coverage information for an address
func (s *CoverageService) GetCoverageInfo(ctx context.Context, addressID string) (*CoverageInfo, error) {
	coverage, err := s.db.GetCoverage(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get coverage info for address %s: %w", addressID, err)
	}

	availableTech, _ := s.ComputeCoverage(ctx, addressID)

	return &CoverageInfo{
		AddressID:     coverage.AddressID,
		City:          coverage.City,
		District:      coverage.District,
		Fiber:         coverage.Fiber,
		VDSL:          coverage.VDSL,
		FWA:           coverage.FWA,
		AvailableTech: availableTech,
	}, nil
}

// CoverageInfo represents coverage information with available technologies
type CoverageInfo struct {
	AddressID     string   `json:"address_id"`
	City          string   `json:"city"`
	District      string   `json:"district"`
	Fiber         bool     `json:"fiber"`
	VDSL          bool     `json:"vdsl"`
	FWA           bool     `json:"fwa"`
	AvailableTech []string `json:"available_tech"`
}
