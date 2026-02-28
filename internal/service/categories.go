package service

import (
	"context"
	"sort"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// CategoryRepository defines category data access where consumed.
type CategoryRepository interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	ListSubcategories(ctx context.Context) ([]domain.Subcategory, error)
}

// CategoryService orchestrates category listing use-cases.
type CategoryService struct {
	repo CategoryRepository
}

// NewCategoryService constructs CategoryService.
func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// ListCategoriesWithSubcategories returns categories with nested subcategories.
func (s *CategoryService) ListCategoriesWithSubcategories(ctx context.Context) ([]domain.CategoryWithSubcategories, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	subcategories, err := s.repo.ListSubcategories(ctx)
	if err != nil {
		return nil, err
	}

	byCategoryID := make(map[int64][]domain.Subcategory, len(categories))
	for _, sub := range subcategories {
		byCategoryID[sub.CategoryID] = append(byCategoryID[sub.CategoryID], sub)
	}

	out := make([]domain.CategoryWithSubcategories, 0, len(categories))
	for _, category := range categories {
		subcategories := append([]domain.Subcategory(nil), byCategoryID[category.ID]...)
		if subcategories == nil {
			subcategories = make([]domain.Subcategory, 0)
		}

		group := domain.CategoryWithSubcategories{
			ID:            category.ID,
			Name:          category.Name,
			ShortName:     category.ShortName,
			Subcategories: subcategories,
		}

		sort.Slice(group.Subcategories, func(i, j int) bool {
			if group.Subcategories[i].ID == group.Subcategories[j].ID {
				return group.Subcategories[i].Name < group.Subcategories[j].Name
			}
			return group.Subcategories[i].ID < group.Subcategories[j].ID
		})
		out = append(out, group)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].ID == out[j].ID {
			return out[i].Name < out[j].Name
		}
		return out[i].ID < out[j].ID
	})

	return out, nil
}
