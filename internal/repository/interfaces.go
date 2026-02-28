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

// PostStore is the shared interface for single-post read operations.
type PostStore interface {
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
}

// HomePostStore is the read contract for homepage posts.
type HomePostStore interface {
	ListRecentActivePosts(ctx context.Context, limit int) ([]domain.Post, error)
	ListRecentActivePostsByCategory(ctx context.Context, categoryID int64, limit int) ([]domain.Post, error)
	ListHomeCategorySections(ctx context.Context) ([]domain.HomeCategorySection, error)
}

// CategoryStore is the shared interface for category taxonomy reads.
type CategoryStore interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	ListSubcategories(ctx context.Context) ([]domain.Subcategory, error)
}
