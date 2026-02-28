package repository

import (
	"context"
	"fmt"
	"sort"
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
	posts    []domain.Post
}

// NewInMemory creates a new in-memory repository pre-loaded with seed data.
func NewInMemory() *InMemory {
	repo := &InMemory{
		listings: make(map[string]domain.Listing),
		posts:    make([]domain.Post, 0),
	}
	repo.loadSeedData()
	repo.loadPostSeedData()
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

func (r *InMemory) ListRecentActivePosts(_ context.Context, limit int) ([]domain.Post, error) {
	if limit <= 0 {
		limit = 50
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	active := make([]domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.Status == domain.PostStatusActive {
			active = append(active, post)
		}
	}

	sort.Slice(active, func(i, j int) bool {
		if active[i].TimePosted == active[j].TimePosted {
			return active[i].ID > active[j].ID
		}
		return active[i].TimePosted > active[j].TimePosted
	})

	if len(active) > limit {
		active = active[:limit]
	}

	out := make([]domain.Post, len(active))
	copy(out, active)
	return out, nil
}

// ListHomeCategorySections returns nil for in-memory mode; home renderer falls back
// to baked category metadata when database taxonomy is unavailable.
func (r *InMemory) ListHomeCategorySections(_ context.Context) ([]domain.HomeCategorySection, error) {
	return nil, nil
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

func (r *InMemory) loadPostSeedData() {
	now := time.Now()
	r.posts = append(r.posts,
		domain.Post{
			ID:            130031901,
			CategoryID:    3,
			SubcategoryID: 14,
			Email:         "alex@stanford.edu",
			Name:          "Sublet room in EVGR premium 2b2b",
			Status:        domain.PostStatusActive,
			TimePosted:    now.Add(-2 * time.Hour).Unix(),
			TimePostedAt:  now.Add(-2 * time.Hour),
			Price:         2000,
			HasPrice:      true,
			HasImage:      false,
			CreatedAt:     now.Add(-2 * time.Hour),
			UpdatedAt:     now.Add(-2 * time.Hour),
		},
		domain.Post{
			ID:            130031900,
			CategoryID:    3,
			SubcategoryID: 14,
			Email:         "casey@stanford.edu",
			Name:          "Shared House",
			Status:        domain.PostStatusActive,
			TimePosted:    now.Add(-3 * time.Hour).Unix(),
			TimePostedAt:  now.Add(-3 * time.Hour),
			Price:         700,
			HasPrice:      true,
			HasImage:      true,
			CreatedAt:     now.Add(-3 * time.Hour),
			UpdatedAt:     now.Add(-3 * time.Hour),
		},
		domain.Post{
			ID:            130031899,
			CategoryID:    5,
			SubcategoryID: 20,
			Email:         "morgan@stanford.edu",
			Name:          "Ikea Stackable beds(2) + 2 Mattresses - Pickup in MV FRIDAY / SATURDAY",
			Status:        domain.PostStatusActive,
			TimePosted:    now.Add(-5 * time.Hour).Unix(),
			TimePostedAt:  now.Add(-5 * time.Hour),
			Price:         0,
			HasPrice:      true,
			HasImage:      true,
			CreatedAt:     now.Add(-5 * time.Hour),
			UpdatedAt:     now.Add(-5 * time.Hour),
		},
		domain.Post{
			ID:            130031898,
			CategoryID:    9,
			SubcategoryID: 90,
			Email:         "sam@stanford.edu",
			Name:          "English tutoring (book club)",
			Status:        domain.PostStatusActive,
			TimePosted:    now.Add(-6 * time.Hour).Unix(),
			TimePostedAt:  now.Add(-6 * time.Hour),
			HasPrice:      false,
			HasImage:      false,
			CreatedAt:     now.Add(-6 * time.Hour),
			UpdatedAt:     now.Add(-6 * time.Hour),
		},
		domain.Post{
			ID:            130031897,
			CategoryID:    5,
			SubcategoryID: 20,
			Email:         "jamie@stanford.edu",
			Name:          "Table and chair",
			Status:        domain.PostStatusActive,
			TimePosted:    now.Add(-15 * time.Hour).Unix(),
			TimePostedAt:  now.Add(-15 * time.Hour),
			Price:         0,
			HasPrice:      true,
			HasImage:      true,
			CreatedAt:     now.Add(-15 * time.Hour),
			UpdatedAt:     now.Add(-15 * time.Hour),
		},
		domain.Post{
			ID:            130031896,
			CategoryID:    5,
			SubcategoryID: 20,
			Email:         "pat@stanford.edu",
			Name:          "Apple Magic Keyboard & Magic Mouse 2 Bundle (Lightning)",
			Status:        0,
			TimePosted:    now.Add(-20 * time.Hour).Unix(),
			TimePostedAt:  now.Add(-20 * time.Hour),
			Price:         65,
			HasPrice:      true,
			HasImage:      true,
			CreatedAt:     now.Add(-20 * time.Hour),
			UpdatedAt:     now.Add(-20 * time.Hour),
		},
	)
}
