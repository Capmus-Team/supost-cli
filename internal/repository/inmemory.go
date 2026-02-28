package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// InMemory implements in-memory data access for posts, categories, and messages.
// Perfect for prototyping and testing — zero external dependencies.
// Swap to Postgres when ready. See AGENTS.md §6.5.
type InMemory struct {
	mu            sync.RWMutex
	posts         []domain.Post
	messages      []domain.Message
	categories    []domain.Category
	subcategories []domain.Subcategory
}

// NewInMemory creates a new in-memory repository pre-loaded with seed data.
func NewInMemory() *InMemory {
	repo := &InMemory{
		posts:         make([]domain.Post, 0),
		messages:      make([]domain.Message, 0),
		categories:    make([]domain.Category, 0),
		subcategories: make([]domain.Subcategory, 0),
	}
	repo.loadPostSeedData()
	repo.loadCategorySeedData()
	return repo
}

func (r *InMemory) GetPostByID(_ context.Context, postID int64) (domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, post := range r.posts {
		if post.ID == postID {
			return post, nil
		}
	}

	return domain.Post{}, domain.ErrNotFound
}

func (r *InMemory) ListRecentActivePosts(_ context.Context, limit int) ([]domain.Post, error) {
	return r.listRecentActivePosts(limit, nil), nil
}

func (r *InMemory) ListRecentActivePostsByCategory(_ context.Context, categoryID int64, limit int) ([]domain.Post, error) {
	return r.listRecentActivePosts(limit, &categoryID), nil
}

func (r *InMemory) ListCategories(_ context.Context) ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Category, len(r.categories))
	copy(out, r.categories)
	return out, nil
}

func (r *InMemory) ListSubcategories(_ context.Context) ([]domain.Subcategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Subcategory, len(r.subcategories))
	copy(out, r.subcategories)
	return out, nil
}

func (r *InMemory) listRecentActivePosts(limit int, categoryID *int64) []domain.Post {
	if limit <= 0 {
		limit = 50
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	active := make([]domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		if post.Status != domain.PostStatusActive {
			continue
		}
		if categoryID != nil && post.CategoryID != *categoryID {
			continue
		}
		active = append(active, post)
	}

	sort.Slice(active, func(i, j int) bool {
		if active[i].TimePosted == active[j].TimePosted {
			return active[i].ID > active[j].ID
		}
		return active[i].TimePosted > active[j].TimePosted
	})

	if len(active) > limit {
		active = active[:limit]
	}

	out := make([]domain.Post, len(active))
	copy(out, active)
	return out
}

// ListHomeCategorySections returns latest active post times per category.
func (r *InMemory) ListHomeCategorySections(_ context.Context) ([]domain.HomeCategorySection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	latestByCategory := make(map[int64]time.Time, 8)
	for _, post := range r.posts {
		if post.Status != domain.PostStatusActive {
			continue
		}
		var postedAt time.Time
		if !post.TimePostedAt.IsZero() {
			postedAt = post.TimePostedAt
		} else if post.TimePosted > 0 {
			postedAt = time.Unix(post.TimePosted, 0)
		}
		if postedAt.IsZero() {
			continue
		}
		if existing, ok := latestByCategory[post.CategoryID]; !ok || postedAt.After(existing) {
			latestByCategory[post.CategoryID] = postedAt
		}
	}

	sections := make([]domain.HomeCategorySection, 0, len(latestByCategory))
	for categoryID, postedAt := range latestByCategory {
		sections = append(sections, domain.HomeCategorySection{
			CategoryID:   categoryID,
			LastPostedAt: postedAt,
		})
	}
	sort.Slice(sections, func(i, j int) bool {
		return sections[i].CategoryID < sections[j].CategoryID
	})
	return sections, nil
}

func (r *InMemory) loadPostSeedData() {
	now := time.Now()
	examplePostedAt := time.Date(2026, time.February, 26, 17, 32, 0, 0, time.FixedZone("PST", -8*60*60))
	r.posts = append(r.posts,
		domain.Post{
			ID:             130031783,
			CategoryID:     9,
			SubcategoryID:  90,
			Email:          "wientjes@alumni.stanford.edu",
			Name:           "Looking for a buddy to go to the movies",
			Body:           "Movie buddy wanted.",
			Status:         domain.PostStatusActive,
			AccessToken:    "dfc6dbef55489317652434afff4caf287c23a1b287dd934c529092fad939260e",
			TimePosted:     examplePostedAt.Unix(),
			TimeModified:   examplePostedAt.Unix(),
			TimePostedAt:   examplePostedAt,
			TimeModifiedAt: examplePostedAt,
			HasPrice:       false,
			HasImage:       false,
			CreatedAt:      examplePostedAt,
			UpdatedAt:      examplePostedAt,
		},
		domain.Post{
			ID:             130031901,
			CategoryID:     3,
			SubcategoryID:  59,
			Email:          "alex@stanford.edu",
			Name:           "Sublet room in EVGR premium 2b2b",
			Status:         domain.PostStatusActive,
			AccessToken:    "token_130031901",
			TimePosted:     now.Add(-2 * time.Hour).Unix(),
			TimeModified:   now.Add(-2 * time.Hour).Unix(),
			TimePostedAt:   now.Add(-2 * time.Hour),
			TimeModifiedAt: now.Add(-2 * time.Hour),
			Price:          2000,
			HasPrice:       true,
			HasImage:       false,
			CreatedAt:      now.Add(-2 * time.Hour),
			UpdatedAt:      now.Add(-2 * time.Hour),
		},
		domain.Post{
			ID:             130031900,
			CategoryID:     3,
			SubcategoryID:  59,
			Email:          "casey@stanford.edu",
			Name:           "Shared House",
			Body:           "Room w/ Private Bathroom for Rent in Quiet Home | Menlo Park. Please do not message this poster about other commercial services.",
			Photo1File:     "post_130031900a.jpg",
			Photo2File:     "post_130031900b.jpg",
			Status:         domain.PostStatusActive,
			AccessToken:    "token_130031900",
			TimePosted:     now.Add(-3 * time.Hour).Unix(),
			TimeModified:   now.Add(-3 * time.Hour).Unix(),
			TimePostedAt:   now.Add(-3 * time.Hour),
			TimeModifiedAt: now.Add(-3 * time.Hour),
			Price:          700,
			HasPrice:       true,
			HasImage:       true,
			CreatedAt:      now.Add(-3 * time.Hour),
			UpdatedAt:      now.Add(-3 * time.Hour),
		},
		domain.Post{
			ID:             130031899,
			CategoryID:     5,
			SubcategoryID:  9,
			Email:          "morgan@stanford.edu",
			Name:           "Ikea Stackable beds(2) + 2 Mattresses - Pickup in MV FRIDAY / SATURDAY",
			Body:           "Pickup in Mountain View on Friday or Saturday.",
			Photo1File:     "post_130031899a.jpg",
			Status:         domain.PostStatusActive,
			AccessToken:    "token_130031899",
			TimePosted:     now.Add(-5 * time.Hour).Unix(),
			TimeModified:   now.Add(-5 * time.Hour).Unix(),
			TimePostedAt:   now.Add(-5 * time.Hour),
			TimeModifiedAt: now.Add(-5 * time.Hour),
			Price:          0,
			HasPrice:       true,
			HasImage:       true,
			CreatedAt:      now.Add(-5 * time.Hour),
			UpdatedAt:      now.Add(-5 * time.Hour),
		},
		domain.Post{
			ID:             130031898,
			CategoryID:     9,
			SubcategoryID:  90,
			Email:          "sam@stanford.edu",
			Name:           "English tutoring (book club)",
			Status:         domain.PostStatusActive,
			AccessToken:    "token_130031898",
			TimePosted:     now.Add(-6 * time.Hour).Unix(),
			TimeModified:   now.Add(-6 * time.Hour).Unix(),
			TimePostedAt:   now.Add(-6 * time.Hour),
			TimeModifiedAt: now.Add(-6 * time.Hour),
			HasPrice:       false,
			HasImage:       false,
			CreatedAt:      now.Add(-6 * time.Hour),
			UpdatedAt:      now.Add(-6 * time.Hour),
		},
		domain.Post{
			ID:             130031897,
			CategoryID:     5,
			SubcategoryID:  9,
			Email:          "jamie@stanford.edu",
			Name:           "Table and chair",
			Status:         domain.PostStatusActive,
			AccessToken:    "token_130031897",
			TimePosted:     now.Add(-15 * time.Hour).Unix(),
			TimeModified:   now.Add(-15 * time.Hour).Unix(),
			TimePostedAt:   now.Add(-15 * time.Hour),
			TimeModifiedAt: now.Add(-15 * time.Hour),
			Price:          0,
			HasPrice:       true,
			HasImage:       true,
			CreatedAt:      now.Add(-15 * time.Hour),
			UpdatedAt:      now.Add(-15 * time.Hour),
		},
		domain.Post{
			ID:             130031896,
			CategoryID:     5,
			SubcategoryID:  9,
			Email:          "pat@stanford.edu",
			Name:           "Apple Magic Keyboard & Magic Mouse 2 Bundle (Lightning)",
			Status:         0,
			AccessToken:    "token_130031896",
			TimePosted:     now.Add(-20 * time.Hour).Unix(),
			TimeModified:   now.Add(-20 * time.Hour).Unix(),
			TimePostedAt:   now.Add(-20 * time.Hour),
			TimeModifiedAt: now.Add(-20 * time.Hour),
			Price:          65,
			HasPrice:       true,
			HasImage:       true,
			CreatedAt:      now.Add(-20 * time.Hour),
			UpdatedAt:      now.Add(-20 * time.Hour),
		},
	)
}

func (r *InMemory) loadCategorySeedData() {
	categories, err := loadSeedCategories()
	if err != nil || len(categories) == 0 {
		categories = defaultSeedCategories()
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].ID < categories[j].ID
	})

	subcategories, err := loadSeedSubcategories()
	if err != nil {
		subcategories = nil
	}
	sort.Slice(subcategories, func(i, j int) bool {
		if subcategories[i].CategoryID == subcategories[j].CategoryID {
			return subcategories[i].ID < subcategories[j].ID
		}
		return subcategories[i].CategoryID < subcategories[j].CategoryID
	})

	r.categories = categories
	r.subcategories = subcategories
}

func loadSeedCategories() ([]domain.Category, error) {
	path, err := seedFilePath("category_rows.json")
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rows []domain.Category
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func loadSeedSubcategories() ([]domain.Subcategory, error) {
	path, err := seedFilePath("subcategory_rows.json")
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rows []domain.Subcategory
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func seedFilePath(filename string) (string, error) {
	_, srcFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("resolving repository path")
	}

	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(srcFile), "..", ".."))
	return filepath.Join(repoRoot, "testdata", "seed", filename), nil
}

func defaultSeedCategories() []domain.Category {
	return []domain.Category{
		{ID: 1, Name: "campus jobs", ShortName: "campus jobs"},
		{ID: 2, Name: "off campus jobs", ShortName: "off campus jobs"},
		{ID: 3, Name: "housing (offering)", ShortName: "housing"},
		{ID: 4, Name: "housing (need)", ShortName: "need housing"},
		{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
		{ID: 6, Name: "resumes", ShortName: "resumes"},
		{ID: 7, Name: "services offered", ShortName: "services"},
		{ID: 8, Name: "personals/dating", ShortName: "personals"},
		{ID: 9, Name: "community", ShortName: "community"},
		{ID: 10, Name: "events", ShortName: "events"},
	}
}
