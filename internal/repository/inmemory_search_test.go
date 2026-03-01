package repository

import (
	"context"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestInMemorySearchActivePosts_ReturnsOnlyActiveNewestFirst(t *testing.T) {
	t.Parallel()

	repo := &InMemory{
		posts: []domain.Post{
			{ID: 10, Status: domain.PostStatusActive, TimePosted: 200, Name: "red chair", Body: "great shape"},
			{ID: 11, Status: 0, TimePosted: 500, Name: "red bike", Body: "inactive"},
			{ID: 12, Status: domain.PostStatusActive, TimePosted: 300, Name: "bike", Body: "red and fast"},
			{ID: 13, Status: domain.PostStatusActive, TimePosted: 300, Name: "desk", Body: "campus pickup"},
		},
	}

	posts, hasMore, err := repo.SearchActivePosts(context.Background(), "", 0, 0, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasMore {
		t.Fatalf("expected hasMore false")
	}

	if len(posts) != 3 {
		t.Fatalf("expected 3 active posts, got %d", len(posts))
	}

	wantIDs := []int64{13, 12, 10}
	for idx, wantID := range wantIDs {
		if posts[idx].ID != wantID {
			t.Fatalf("expected post id %d at index %d, got %d", wantID, idx, posts[idx].ID)
		}
		if posts[idx].Status != domain.PostStatusActive {
			t.Fatalf("expected active status at index %d, got %d", idx, posts[idx].Status)
		}
	}
}

func TestInMemorySearchActivePosts_Paginates(t *testing.T) {
	t.Parallel()

	repo := &InMemory{
		posts: []domain.Post{
			{ID: 30, Status: domain.PostStatusActive, TimePosted: 300, Name: "desk"},
			{ID: 20, Status: domain.PostStatusActive, TimePosted: 200, Name: "chair"},
			{ID: 10, Status: domain.PostStatusActive, TimePosted: 100, Name: "lamp"},
		},
	}

	firstPage, hasMore, err := repo.SearchActivePosts(context.Background(), "", 0, 0, 1, 2)
	if err != nil {
		t.Fatalf("unexpected first-page error: %v", err)
	}
	if !hasMore {
		t.Fatalf("expected hasMore true on first page")
	}
	if len(firstPage) != 2 {
		t.Fatalf("expected 2 posts on first page, got %d", len(firstPage))
	}
	if firstPage[0].ID != 30 || firstPage[1].ID != 20 {
		t.Fatalf("unexpected first-page ids: %d, %d", firstPage[0].ID, firstPage[1].ID)
	}

	secondPage, hasMore, err := repo.SearchActivePosts(context.Background(), "", 0, 0, 2, 2)
	if err != nil {
		t.Fatalf("unexpected second-page error: %v", err)
	}
	if hasMore {
		t.Fatalf("expected hasMore false on second page")
	}
	if len(secondPage) != 1 || secondPage[0].ID != 10 {
		t.Fatalf("unexpected second-page result: %+v", secondPage)
	}
}

func TestInMemorySearchActivePosts_QueryMatchesNameOrBody(t *testing.T) {
	t.Parallel()

	repo := &InMemory{
		posts: []domain.Post{
			{ID: 40, Status: domain.PostStatusActive, TimePosted: 400, Name: "Red bike", Body: "Pick up on campus"},
			{ID: 30, Status: domain.PostStatusActive, TimePosted: 300, Name: "Blue couch", Body: "Stanford poster included"},
			{ID: 20, Status: domain.PostStatusActive, TimePosted: 200, Name: "Bike lock", Body: "metal"},
			{ID: 10, Status: 0, TimePosted: 500, Name: "Red bike", Body: "inactive post"},
		},
	}

	posts, hasMore, err := repo.SearchActivePosts(context.Background(), "red bike", 0, 0, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasMore {
		t.Fatalf("expected hasMore false")
	}
	if len(posts) != 1 || posts[0].ID != 40 {
		t.Fatalf("expected only active post 40, got %+v", posts)
	}

	posts, hasMore, err = repo.SearchActivePosts(context.Background(), "stanford poster", 0, 0, 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasMore {
		t.Fatalf("expected hasMore false")
	}
	if len(posts) != 1 || posts[0].ID != 30 {
		t.Fatalf("expected post 30 match from body terms, got %+v", posts)
	}
}
