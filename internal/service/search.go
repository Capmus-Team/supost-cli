package service

import (
	"context"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const (
	defaultSearchPage    = 1
	defaultSearchPerPage = 100
	maxSearchPerPage     = 100
)

// SearchRepository defines search read operations where consumed.
type SearchRepository interface {
	SearchActivePosts(ctx context.Context, categoryID, subcategoryID int64, page, perPage int) ([]domain.Post, bool, error)
}

// SearchService orchestrates search page retrieval.
type SearchService struct {
	repo SearchRepository
}

// NewSearchService constructs SearchService.
func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

// Search returns paginated active posts for optional category/subcategory filters.
func (s *SearchService) Search(ctx context.Context, categoryID, subcategoryID int64, page, perPage int) (domain.SearchResultPage, error) {
	page = normalizeSearchPage(page)
	perPage = normalizeSearchPerPage(perPage)

	posts, hasMore, err := s.repo.SearchActivePosts(ctx, categoryID, subcategoryID, page, perPage)
	if err != nil {
		return domain.SearchResultPage{}, err
	}

	return domain.SearchResultPage{
		CategoryID:    categoryID,
		SubcategoryID: subcategoryID,
		Page:          page,
		PerPage:       perPage,
		HasMore:       hasMore,
		Posts:         posts,
	}, nil
}

func normalizeSearchPage(page int) int {
	if page < 1 {
		return defaultSearchPage
	}
	return page
}

func normalizeSearchPerPage(perPage int) int {
	if perPage <= 0 {
		return defaultSearchPerPage
	}
	if perPage > maxSearchPerPage {
		return maxSearchPerPage
	}
	return perPage
}
