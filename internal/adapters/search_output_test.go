package adapters

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestRenderSearchResults_GroupsByDateAndShowsNextPage(t *testing.T) {
	var out bytes.Buffer
	base := time.Date(2026, time.February, 27, 20, 37, 0, 0, time.UTC)
	result := domain.SearchResultPage{
		CategoryID: 3,
		Page:       1,
		PerPage:    100,
		HasMore:    true,
		Posts: []domain.Post{
			{
				ID:            1,
				Name:          "Shared House - $700",
				CategoryID:    3,
				SubcategoryID: 59,
				Email:         "person@stanford.edu",
				HasPrice:      true,
				Price:         700,
				TimePostedAt:  base,
			},
			{
				ID:            2,
				Name:          "Large furnished room in ...",
				CategoryID:    3,
				SubcategoryID: 59,
				Email:         "person@stanford.edu",
				TimePostedAt:  base.Add(-24 * time.Hour),
			},
		},
	}

	if err := RenderSearchResults(&out, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{
		"SUPost » Stanford, California » housing",
		"all posts",
		"Fri, Feb 27",
		"Thu, Feb 26",
		"Shared House - $700",
		"next 100 posts",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in rendered search page", needle)
		}
	}
}

func TestRenderSearchResults_ShowsPreviousAndNextWhenPageGreaterThanOne(t *testing.T) {
	var out bytes.Buffer
	result := domain.SearchResultPage{
		CategoryID: 5,
		Page:       2,
		PerPage:    20,
		HasMore:    true,
		Posts: []domain.Post{
			{
				ID:            1,
				Name:          "Desk - $80",
				CategoryID:    5,
				SubcategoryID: 14,
				Email:         "person@stanford.edu",
				TimePostedAt:  time.Date(2026, time.February, 27, 20, 37, 0, 0, time.UTC),
			},
		},
	}

	if err := RenderSearchResults(&out, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	for _, needle := range []string{"previous 20 posts", "next 20 posts"} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in rendered search page", needle)
		}
	}
}

func TestRenderSearchResults_ShowsEmptyStateForNoPosts(t *testing.T) {
	var out bytes.Buffer
	result := domain.SearchResultPage{
		SubcategoryID: 14,
		Page:          3,
		PerPage:       30,
		Posts:         []domain.Post{},
	}

	if err := RenderSearchResults(&out, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	if !strings.Contains(plain, "No posts found for this page.") {
		t.Fatalf("missing empty-state message in %q", plain)
	}
	if !strings.Contains(plain, "previous 30 posts") {
		t.Fatalf("missing previous-page label in %q", plain)
	}
}

func TestRenderSearchResults_SubcategoryOnlyInfersParentCategoryInBreadcrumb(t *testing.T) {
	var out bytes.Buffer
	result := domain.SearchResultPage{
		SubcategoryID: 14,
		Page:          1,
		PerPage:       100,
		Posts:         []domain.Post{},
	}

	if err := RenderSearchResults(&out, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	plain := stripANSI(out.String())
	if !strings.Contains(plain, "SUPost » Stanford, California » for sale » furniture") {
		t.Fatalf("missing inferred category+subcategory breadcrumb in %q", plain)
	}
	if !strings.Contains(plain, "all posts") {
		t.Fatalf("missing all-posts title in %q", plain)
	}
}
