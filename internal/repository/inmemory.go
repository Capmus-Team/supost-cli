package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// InMemory implements ListingStore using an in-memory map.
// Perfect for prototyping and testing — zero external dependencies.
// Swap to Postgres when ready. See AGENTS.md §6.5.
type InMemory struct {
	mu       sync.RWMutex
	listings map[string]domain.Listing
}

// NewInMemory creates a new in-memory repository pre-loaded with seed data.
func NewInMemory() *InMemory {
	repo := &InMemory{
		listings: make(map[string]domain.Listing),
	}
	repo.loadSeedData()
	return repo
}

func (r *InMemory) ListActive(_ context.Context) ([]domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var active []domain.Listing
	for _, l := range r.listings {
		if l.Status == "active" {
			active = append(active, l)
		}
	}
	return active, nil
}

func (r *InMemory) GetByID(_ context.Context, id string) (*domain.Listing, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	l, ok := r.listings[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &l, nil
}

func (r *InMemory) Create(_ context.Context, listing *domain.Listing) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if listing.ID == "" {
		listing.ID = fmt.Sprintf("mem-%d", len(r.listings)+1)
	}
	now := time.Now()
	listing.CreatedAt = now
	listing.UpdatedAt = now
	if listing.Status == "" {
		listing.Status = "active"
	}

	r.listings[listing.ID] = *listing
	return nil
}

// loadSeedData populates the repository with sample data.
// In a more advanced setup, this could read from testdata/seed/*.json.
func (r *InMemory) loadSeedData() {
	now := time.Now()
	seeds := []domain.Listing{
		{
			ID:          "seed-1",
			UserID:      "user-1",
			Title:       "Used Calculus Textbook",
			Description: "Stewart Calculus, 8th edition. Some highlighting.",
			Price:       4500,
			Category:    "textbooks",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "seed-2",
			UserID:      "user-1",
			Title:       "IKEA Desk",
			Description: "MALM desk, white. Good condition. Pickup only.",
			Price:       6000,
			Category:    "furniture",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "seed-3",
			UserID:      "user-2",
			Title:       "Trek Road Bike",
			Description: "2021 Domane AL 2, 56cm. Low miles.",
			Price:       55000,
			Category:    "bikes",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, s := range seeds {
		r.listings[s.ID] = s
	}
}
