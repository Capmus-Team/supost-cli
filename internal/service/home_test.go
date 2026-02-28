package service

import (
	"context"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockHomeRepo struct {
	receivedLimit int
	posts         []domain.Post
	sections      []domain.HomeCategorySection
}

func (m *mockHomeRepo) ListRecentActivePosts(_ context.Context, limit int) ([]domain.Post, error) {
	m.receivedLimit = limit
	return m.posts, nil
}

func (m *mockHomeRepo) ListHomeCategorySections(_ context.Context) ([]domain.HomeCategorySection, error) {
	return m.sections, nil
}

func TestHomeService_ListRecentActive_UsesDefaultLimit(t *testing.T) {
	repo := &mockHomeRepo{}
	svc := NewHomeService(repo)

	if _, err := svc.ListRecentActive(context.Background(), 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.receivedLimit != defaultHomeLimit {
		t.Fatalf("expected default limit %d, got %d", defaultHomeLimit, repo.receivedLimit)
	}
}

func TestHomeService_ListRecentActive_UsesProvidedLimit(t *testing.T) {
	repo := &mockHomeRepo{}
	svc := NewHomeService(repo)

	if _, err := svc.ListRecentActive(context.Background(), 12); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.receivedLimit != 12 {
		t.Fatalf("expected limit 12, got %d", repo.receivedLimit)
	}
}

func TestHomeService_ListCategorySections(t *testing.T) {
	repo := &mockHomeRepo{
		sections: []domain.HomeCategorySection{
			{CategoryID: domain.CategoryHousing, CategoryName: "housing"},
		},
	}
	svc := NewHomeService(repo)

	sections, err := svc.ListCategorySections(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sections))
	}
	if sections[0].CategoryID != domain.CategoryHousing {
		t.Fatalf("expected housing category id, got %d", sections[0].CategoryID)
	}
}
