package repository

import (
	"context"
	"sort"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *InMemory) SearchActivePosts(_ context.Context, categoryID, subcategoryID int64, page, perPage int) ([]domain.Post, bool, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.Status != domain.PostStatusActive {
			continue
		}
		if categoryID > 0 && post.CategoryID != categoryID {
			continue
		}
		if subcategoryID > 0 && post.SubcategoryID != subcategoryID {
			continue
		}
		filtered = append(filtered, post)
	}

	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].TimePosted == filtered[j].TimePosted {
			return filtered[i].ID > filtered[j].ID
		}
		return filtered[i].TimePosted > filtered[j].TimePosted
	})

	offset := (page - 1) * perPage
	if offset >= len(filtered) {
		return []domain.Post{}, false, nil
	}

	end := offset + perPage
	hasMore := end < len(filtered)
	if end > len(filtered) {
		end = len(filtered)
	}

	out := make([]domain.Post, end-offset)
	copy(out, filtered[offset:end])
	return out, hasMore, nil
}
