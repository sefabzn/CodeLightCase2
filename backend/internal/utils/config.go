package utils

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	Port               string
	DatabaseURL        string
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Port:               getEnvWithDefault("PORT", "8000"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validate checks that all required configuration is present
func (c *Config) validate() error {
	var missingVars []string

	if c.DatabaseURL == "" {
		missingVars = append(missingVars, "DATABASE_URL")
	}
	if c.SupabaseURL == "" {
		missingVars = append(missingVars, "SUPABASE_URL")
	}
	if c.SupabaseAnonKey == "" {
		missingVars = append(missingVars, "SUPABASE_ANON_KEY")
	}
	if c.SupabaseServiceKey == "" {
		missingVars = append(missingVars, "SUPABASE_SERVICE_ROLE_KEY")
	}

	// Validate port is a valid number
	if _, err := strconv.Atoi(c.Port); err != nil {
		missingVars = append(missingVars, "PORT (must be a valid number)")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing or invalid environment variables: %v", missingVars)
	}

	return nil
}

// GetAddr returns the server address with port
func (c *Config) GetAddr() string {
	return ":" + c.Port
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
