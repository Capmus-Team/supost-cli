// Package repository handles all data access.
// The interface is defined in service/ (where consumed). This package
// provides concrete implementations: inmemory.go for prototyping,
// postgres.go for production.
// See AGENTS.md §2.4, §5.8, §6.5.
package repository

import (
	"context"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// ListingStore is the shared interface for listing data access.
// Defined here because both inmemory and postgres adapters implement it,
// and cmd/ needs to reference the concrete constructors.
type ListingStore interface {
	ListActive(ctx context.Context) ([]domain.Listing, error)
	GetByID(ctx context.Context, id string) (*domain.Listing, error)
	Create(ctx context.Context, listing *domain.Listing) error
}
