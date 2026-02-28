package service

import (
	"context"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockCategoryRepo struct {
	categories    []domain.Category
	subcategories []domain.Subcategory
}

func (m *mockCategoryRepo) ListCategories(_ context.Context) ([]domain.Category, error) {
	return m.categories, nil
}

func (m *mockCategoryRepo) ListSubcategories(_ context.Context) ([]domain.Subcategory, error) {
	return m.subcategories, nil
}

func TestCategoryService_ListCategoriesWithSubcategories(t *testing.T) {
	repo := &mockCategoryRepo{
		categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
			{ID: 2, Name: "off campus jobs", ShortName: "off campus jobs"},
		},
		subcategories: []domain.Subcategory{
			{ID: 160, CategoryID: 2, Name: "paid interns"},
			{ID: 159, CategoryID: 2, Name: "part-time"},
			{ID: 12, CategoryID: 5, Name: "electronics"},
		},
	}

	svc := NewCategoryService(repo)
	got, err := svc.ListCategoriesWithSubcategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(got))
	}
	if got[0].ID != 2 || got[1].ID != 5 {
		t.Fatalf("categories not sorted by id: got [%d, %d]", got[0].ID, got[1].ID)
	}
	if len(got[0].Subcategories) != 2 {
		t.Fatalf("expected 2 subcategories for category 2, got %d", len(got[0].Subcategories))
	}
	if got[0].Subcategories[0].ID != 159 || got[0].Subcategories[1].ID != 160 {
		t.Fatalf("subcategory sort mismatch for category 2")
	}
}
