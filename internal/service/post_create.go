package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// PostCreateRepository defines taxonomy reads for post creation pages.
type PostCreateRepository interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	ListSubcategories(ctx context.Context) ([]domain.Subcategory, error)
	CreatePendingPost(ctx context.Context, submission domain.PostCreateSubmission) (domain.PostCreatePersisted, error)
	SavePostPhotos(ctx context.Context, photos []domain.PostCreateSavedPhoto) error
}

// PostCreateService builds the staged post creation flow.
type PostCreateService struct {
	repo PostCreateRepository
}

// NewPostCreateService constructs PostCreateService.
func NewPostCreateService(repo PostCreateRepository) *PostCreateService {
	return &PostCreateService{repo: repo}
}

// BuildPage returns the requested post-create stage based on selected IDs.
func (s *PostCreateService) BuildPage(ctx context.Context, categoryID, subcategoryID int64) (domain.PostCreatePage, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return domain.PostCreatePage{}, err
	}
	subcategories, err := s.repo.ListSubcategories(ctx)
	if err != nil {
		return domain.PostCreatePage{}, err
	}

	sortCategoriesByID(categories)
	categoryByID := indexCategoriesByID(categories)
	subcategoryByID := indexSubcategoriesByID(subcategories)

	if categoryID <= 0 && subcategoryID > 0 {
		if sub, ok := subcategoryByID[subcategoryID]; ok {
			categoryID = sub.CategoryID
		}
	}

	if categoryID <= 0 {
		return domain.PostCreatePage{
			Stage:      domain.PostCreateStageChooseCategory,
			Categories: append([]domain.Category(nil), categories...),
		}, nil
	}

	category, ok := categoryByID[categoryID]
	if !ok {
		return domain.PostCreatePage{}, fmt.Errorf("category %d not found: %w", categoryID, domain.ErrNotFound)
	}

	filteredSubs := filterSubcategoriesByCategoryID(subcategories, categoryID)
	sortSubcategoriesByName(filteredSubs)

	if subcategoryID <= 0 {
		return domain.PostCreatePage{
			Stage:         domain.PostCreateStageChooseSubcategory,
			CategoryID:    categoryID,
			CategoryName:  normalizedCategoryName(category),
			Subcategories: filteredSubs,
		}, nil
	}

	subcategory, ok := subcategoryByID[subcategoryID]
	if !ok || subcategory.CategoryID != categoryID {
		return domain.PostCreatePage{}, fmt.Errorf("subcategory %d not found in category %d: %w", subcategoryID, categoryID, domain.ErrNotFound)
	}

	return domain.PostCreatePage{
		Stage:           domain.PostCreateStageForm,
		CategoryID:      categoryID,
		SubcategoryID:   subcategoryID,
		CategoryName:    normalizedCategoryName(category),
		SubcategoryName: strings.TrimSpace(subcategory.Name),
	}, nil
}

func sortCategoriesByID(categories []domain.Category) {
	sort.Slice(categories, func(i, j int) bool {
		if categories[i].ID == categories[j].ID {
			return categories[i].Name < categories[j].Name
		}
		return categories[i].ID < categories[j].ID
	})
}

func sortSubcategoriesByName(subcategories []domain.Subcategory) {
	sort.Slice(subcategories, func(i, j int) bool {
		left := strings.ToLower(strings.TrimSpace(subcategories[i].Name))
		right := strings.ToLower(strings.TrimSpace(subcategories[j].Name))
		if left == right {
			return subcategories[i].ID < subcategories[j].ID
		}
		return left < right
	})
}

func indexCategoriesByID(categories []domain.Category) map[int64]domain.Category {
	out := make(map[int64]domain.Category, len(categories))
	for _, category := range categories {
		out[category.ID] = category
	}
	return out
}

func indexSubcategoriesByID(subcategories []domain.Subcategory) map[int64]domain.Subcategory {
	out := make(map[int64]domain.Subcategory, len(subcategories))
	for _, subcategory := range subcategories {
		out[subcategory.ID] = subcategory
	}
	return out
}

func filterSubcategoriesByCategoryID(subcategories []domain.Subcategory, categoryID int64) []domain.Subcategory {
	filtered := make([]domain.Subcategory, 0, len(subcategories))
	for _, subcategory := range subcategories {
		if subcategory.CategoryID == categoryID {
			filtered = append(filtered, subcategory)
		}
	}
	return filtered
}

func normalizedCategoryName(category domain.Category) string {
	if short := strings.TrimSpace(category.ShortName); short != "" {
		return short
	}
	return strings.TrimSpace(category.Name)
}
