package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockPostCreateRepo struct {
	categories    []domain.Category
	subcategories []domain.Subcategory
}

func (m *mockPostCreateRepo) ListCategories(_ context.Context) ([]domain.Category, error) {
	return m.categories, nil
}

func (m *mockPostCreateRepo) ListSubcategories(_ context.Context) ([]domain.Subcategory, error) {
	return m.subcategories, nil
}

func (m *mockPostCreateRepo) CreatePendingPost(_ context.Context, _ domain.PostCreateSubmission) (domain.PostCreatePersisted, error) {
	return domain.PostCreatePersisted{}, nil
}

func TestPostCreateService_BuildPage_ChooseCategory(t *testing.T) {
	svc := NewPostCreateService(&mockPostCreateRepo{
		categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
			{ID: 8, Name: "personals/dating", ShortName: "personals"},
		},
	})

	page, err := svc.BuildPage(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Stage != domain.PostCreateStageChooseCategory {
		t.Fatalf("unexpected stage %q", page.Stage)
	}
	if len(page.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(page.Categories))
	}
}

func TestPostCreateService_BuildPage_ChooseSubcategory(t *testing.T) {
	svc := NewPostCreateService(&mockPostCreateRepo{
		categories: []domain.Category{
			{ID: 8, Name: "personals/dating", ShortName: "personals"},
		},
		subcategories: []domain.Subcategory{
			{ID: 131, CategoryID: 8, Name: "girl wants girl"},
			{ID: 130, CategoryID: 8, Name: "friendship"},
		},
	})

	page, err := svc.BuildPage(context.Background(), 8, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Stage != domain.PostCreateStageChooseSubcategory {
		t.Fatalf("unexpected stage %q", page.Stage)
	}
	if len(page.Subcategories) != 2 {
		t.Fatalf("expected 2 subcategories, got %d", len(page.Subcategories))
	}
	if page.Subcategories[0].Name != "friendship" {
		t.Fatalf("expected alphabetical ordering, got first %q", page.Subcategories[0].Name)
	}
}

func TestPostCreateService_BuildPage_FormStage(t *testing.T) {
	svc := NewPostCreateService(&mockPostCreateRepo{
		categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
		},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	})

	page, err := svc.BuildPage(context.Background(), 5, 14)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Stage != domain.PostCreateStageForm {
		t.Fatalf("unexpected stage %q", page.Stage)
	}
	if page.CategoryName != "for sale" || page.SubcategoryName != "furniture" {
		t.Fatalf("unexpected names: category=%q subcategory=%q", page.CategoryName, page.SubcategoryName)
	}
}

func TestPostCreateService_BuildPage_InfersCategoryFromSubcategory(t *testing.T) {
	svc := NewPostCreateService(&mockPostCreateRepo{
		categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
		},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	})

	page, err := svc.BuildPage(context.Background(), 0, 14)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.Stage != domain.PostCreateStageForm {
		t.Fatalf("unexpected stage %q", page.Stage)
	}
	if page.CategoryID != 5 || page.SubcategoryID != 14 {
		t.Fatalf("unexpected IDs category=%d subcategory=%d", page.CategoryID, page.SubcategoryID)
	}
}

func TestPostCreateService_BuildPage_InvalidCategory(t *testing.T) {
	svc := NewPostCreateService(&mockPostCreateRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
	})

	_, err := svc.BuildPage(context.Background(), 99, 0)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
