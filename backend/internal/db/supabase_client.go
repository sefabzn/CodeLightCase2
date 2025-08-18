package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"app/internal/models"
)

// SupabaseClient represents a client for Supabase REST API
type SupabaseClient struct {
	baseURL    string
	anonKey    string
	serviceKey string
	httpClient *http.Client
}

// NewSupabaseClient creates a new Supabase client
func NewSupabaseClient(baseURL, anonKey, serviceKey string) *SupabaseClient {
	return &SupabaseClient{
		baseURL:    baseURL,
		anonKey:    anonKey,
		serviceKey: serviceKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Health checks if the Supabase service is available
func (s *SupabaseClient) Health(ctx context.Context) error {
	url := s.baseURL + "/rest/v1/users?limit=1"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("apikey", s.anonKey)
	req.Header.Set("Authorization", "Bearer "+s.anonKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// Close closes the client (no-op for HTTP client)
func (s *SupabaseClient) Close() {
	// No resources to close for HTTP client
}

// Helper method to make GET requests
func (s *SupabaseClient) get(ctx context.Context, endpoint string, result interface{}) error {
	url := s.baseURL + "/rest/v1/" + endpoint

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.serviceKey)
	req.Header.Set("Authorization", "Bearer "+s.serviceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID
func (s *SupabaseClient) GetUser(ctx context.Context, userID int) (*models.User, error) {
	endpoint := fmt.Sprintf("users?user_id=eq.%d&limit=1", userID)

	var users []models.User
	if err := s.get(ctx, endpoint, &users); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user %d not found", userID)
	}

	return &users[0], nil
}

// GetCoverage retrieves coverage information for an address
func (s *SupabaseClient) GetCoverage(ctx context.Context, addressID string) (*models.Coverage, error) {
	endpoint := fmt.Sprintf("coverage?address_id=eq.%s&limit=1", addressID)

	var coverage []models.Coverage
	if err := s.get(ctx, endpoint, &coverage); err != nil {
		return nil, fmt.Errorf("failed to get coverage: %w", err)
	}

	if len(coverage) == 0 {
		return nil, fmt.Errorf("coverage for address %s not found", addressID)
	}

	return &coverage[0], nil
}

// GetHousehold retrieves household information for a user
func (s *SupabaseClient) GetHousehold(ctx context.Context, userID int) ([]models.Household, error) {
	endpoint := fmt.Sprintf("household?user_id=eq.%d", userID)

	var household []models.Household
	if err := s.get(ctx, endpoint, &household); err != nil {
		return nil, fmt.Errorf("failed to get household: %w", err)
	}

	return household, nil
}

// GetInstallSlots retrieves install slots for an address and technology
func (s *SupabaseClient) GetInstallSlots(ctx context.Context, addressID, tech string) ([]models.InstallSlot, error) {
	endpoint := fmt.Sprintf("install_slots?address_id=eq.%s&tech=eq.%s&available=eq.true&order=slot_start", addressID, tech)

	var slots []models.InstallSlot
	if err := s.get(ctx, endpoint, &slots); err != nil {
		return nil, fmt.Errorf("failed to get install slots: %w", err)
	}

	return slots, nil
}

// GetCatalog retrieves all plan catalogs
func (s *SupabaseClient) GetCatalog(ctx context.Context) (*models.Catalog, error) {
	var catalog models.Catalog

	// Get mobile plans
	if err := s.get(ctx, "mobile_plans", &catalog.MobilePlans); err != nil {
		return nil, fmt.Errorf("failed to get mobile plans: %w", err)
	}

	// Get home plans
	if err := s.get(ctx, "home_plans", &catalog.HomePlans); err != nil {
		return nil, fmt.Errorf("failed to get home plans: %w", err)
	}

	// Get TV plans
	if err := s.get(ctx, "tv_plans", &catalog.TVPlans); err != nil {
		return nil, fmt.Errorf("failed to get TV plans: %w", err)
	}

	// Get bundling rules
	if err := s.get(ctx, "bundling_rules", &catalog.BundlingRules); err != nil {
		return nil, fmt.Errorf("failed to get bundling rules: %w", err)
	}

	return &catalog, nil
}
