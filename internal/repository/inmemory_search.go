package repository

import (
	"context"
	"sort"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *InMemory) SearchActivePosts(_ context.Context, query string, categoryID, subcategoryID int64, page, perPage int) ([]domain.Post, bool, error) {
	query = strings.TrimSpace(query)
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
		if !matchesPostQuery(post, query) {
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

func matchesPostQuery(post domain.Post, query string) bool {
	if query == "" {
		return true
	}

	nameLower := strings.ToLower(post.Name)
	bodyLower := strings.ToLower(post.Body)
	for _, term := range strings.Fields(strings.ToLower(query)) {
		if term == "" {
			continue
		}
		if !strings.Contains(nameLower, term) && !strings.Contains(bodyLower, term) {
			return false
		}
	}
	return true
}
