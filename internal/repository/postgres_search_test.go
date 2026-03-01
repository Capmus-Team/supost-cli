package repository

import (
	"strings"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestBuildSearchActivePostsStatement_UsesDefaultQueryWithoutKeyword(t *testing.T) {
	querySQL, queryArgs, perPage := buildSearchActivePostsStatement("   ", 7, 11, 2, 25)

	if perPage != 25 {
		t.Fatalf("expected per_page 25, got %d", perPage)
	}
	if len(queryArgs) != 5 {
		t.Fatalf("expected 5 args, got %d", len(queryArgs))
	}
	if got, ok := queryArgs[0].(int); !ok || got != domain.PostStatusActive {
		t.Fatalf("expected arg0 active status %d, got %#v", domain.PostStatusActive, queryArgs[0])
	}
	if got, ok := queryArgs[1].(int64); !ok || got != 7 {
		t.Fatalf("expected arg1 category 7, got %#v", queryArgs[1])
	}
	if got, ok := queryArgs[2].(int64); !ok || got != 11 {
		t.Fatalf("expected arg2 subcategory 11, got %#v", queryArgs[2])
	}
	if got, ok := queryArgs[3].(int); !ok || got != 26 {
		t.Fatalf("expected arg3 limit 26, got %#v", queryArgs[3])
	}
	if got, ok := queryArgs[4].(int); !ok || got != 25 {
		t.Fatalf("expected arg4 offset 25, got %#v", queryArgs[4])
	}

	for _, needle := range []string{
		"WHERE p.status = $1",
		"p.category_id = $2",
		"p.subcategory_id = $3",
		"FROM public.post p",
		"ORDER BY p.time_posted DESC, p.id DESC",
		"LIMIT $4 OFFSET $5",
	} {
		if !strings.Contains(querySQL, needle) {
			t.Fatalf("expected SQL to contain %q", needle)
		}
	}
	if strings.Contains(querySQL, " OR ") {
		t.Fatalf("did not expect OR-guard filter pattern in SQL: %q", querySQL)
	}
	if strings.Contains(querySQL, "NULLS LAST") {
		t.Fatalf("did not expect NULLS LAST in SQL: %q", querySQL)
	}
}

func TestBuildSearchActivePostsStatement_UsesFTSQueryWithKeyword(t *testing.T) {
	querySQL, queryArgs, perPage := buildSearchActivePostsStatement("  stanford bike  ", 0, 0, 3, 10)

	if perPage != 10 {
		t.Fatalf("expected per_page 10, got %d", perPage)
	}
	if len(queryArgs) != 4 {
		t.Fatalf("expected 4 args, got %d", len(queryArgs))
	}
	if got, ok := queryArgs[1].(string); !ok || got != "stanford bike" {
		t.Fatalf("expected trimmed keyword query, got %#v", queryArgs[1])
	}
	if got, ok := queryArgs[2].(int); !ok || got != 11 {
		t.Fatalf("expected arg2 limit 11, got %#v", queryArgs[2])
	}
	if got, ok := queryArgs[3].(int); !ok || got != 20 {
		t.Fatalf("expected arg3 offset 20, got %#v", queryArgs[3])
	}

	for _, needle := range []string{
		"CROSS JOIN plainto_tsquery('english', $2) q",
		"p.fts @@ q",
		"ts_rank(p.fts, q)",
		"LIMIT $3 OFFSET $4",
	} {
		if !strings.Contains(querySQL, needle) {
			t.Fatalf("expected FTS SQL to contain %q", needle)
		}
	}
}

func TestBuildSearchActivePostsStatement_NormalizesPagingDefaults(t *testing.T) {
	_, queryArgs, perPage := buildSearchActivePostsStatement("keyword", 0, 0, 0, 0)

	if perPage != 100 {
		t.Fatalf("expected default per_page 100, got %d", perPage)
	}
	if got, ok := queryArgs[2].(int); !ok || got != 101 {
		t.Fatalf("expected default limit 101, got %#v", queryArgs[2])
	}
	if got, ok := queryArgs[3].(int); !ok || got != 0 {
		t.Fatalf("expected default offset 0, got %#v", queryArgs[3])
	}
}

func TestBuildSearchActivePostsStatement_OmitsUnsetFilters(t *testing.T) {
	querySQL, queryArgs, _ := buildSearchActivePostsStatement("", 0, 0, 1, 10)

	if len(queryArgs) != 3 {
		t.Fatalf("expected 3 args, got %d", len(queryArgs))
	}
	if strings.Contains(querySQL, "category_id =") {
		t.Fatalf("did not expect category filter when category is unset")
	}
	if strings.Contains(querySQL, "subcategory_id =") {
		t.Fatalf("did not expect subcategory filter when subcategory is unset")
	}
}
