package service

import (
	"context"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockSearchRepo struct {
	categoryID    int64
	subcategoryID int64
	page          int
	perPage       int
	posts         []domain.Post
	hasMore       bool
}

func (m *mockSearchRepo) SearchActivePosts(_ context.Context, categoryID, subcategoryID int64, page, perPage int) ([]domain.Post, bool, error) {
	m.categoryID = categoryID
	m.subcategoryID = subcategoryID
	m.page = page
	m.perPage = perPage
	return m.posts, m.hasMore, nil
}

func TestSearchService_Search_NormalizesPaging(t *testing.T) {
	repo := &mockSearchRepo{
		posts:   []domain.Post{{ID: 1}},
		hasMore: true,
	}
	svc := NewSearchService(repo)

	result, err := svc.Search(context.Background(), 3, 59, 0, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.page != 1 {
		t.Fatalf("expected normalized page 1, got %d", repo.page)
	}
	if repo.perPage != 100 {
		t.Fatalf("expected normalized per_page 100, got %d", repo.perPage)
	}
	if result.Page != 1 || result.PerPage != 100 {
		t.Fatalf("unexpected result paging: page=%d per_page=%d", result.Page, result.PerPage)
	}
	if !result.HasMore {
		t.Fatalf("expected has_more true")
	}
}

func TestSearchService_Search_ForwardsFilters(t *testing.T) {
	repo := &mockSearchRepo{}
	svc := NewSearchService(repo)

	_, err := svc.Search(context.Background(), 5, 9, 2, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.categoryID != 5 || repo.subcategoryID != 9 {
		t.Fatalf("expected forwarded category/subcategory, got %d/%d", repo.categoryID, repo.subcategoryID)
	}
}
