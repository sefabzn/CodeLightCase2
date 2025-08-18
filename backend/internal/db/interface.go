package db

import (
	"context"

	"app/internal/models"
)

// DatabaseInterface defines the interface for database operations
type DatabaseInterface interface {
	Health(ctx context.Context) error
	Close()
	GetUser(ctx context.Context, userID int) (*models.User, error)
	GetCoverage(ctx context.Context, addressID string) (*models.Coverage, error)
	GetHousehold(ctx context.Context, userID int) ([]models.Household, error)
	GetInstallSlots(ctx context.Context, addressID, tech string) ([]models.InstallSlot, error)
	GetCatalog(ctx context.Context) (*models.Catalog, error)
}
