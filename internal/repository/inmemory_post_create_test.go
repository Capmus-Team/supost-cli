package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestInMemorySavePostPhotos_SetsHasImageAndStoresRows(t *testing.T) {
	repo := NewInMemory()
	now := time.Now()

	persisted, err := repo.CreatePendingPost(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Email:         "wientjes@alumni.stanford.edu",
		Name:          "Bike",
		Body:          "Body",
		PostedAt:      now,
	})
	if err != nil {
		t.Fatalf("creating pending post: %v", err)
	}

	err = repo.SavePostPhotos(context.Background(), []domain.PostCreateSavedPhoto{
		{PostID: persisted.PostID, S3Key: "v2/posts/1/a.jpg", Position: 0},
		{PostID: persisted.PostID, S3Key: "v2/posts/1/b.jpg", Position: 1},
	})
	if err != nil {
		t.Fatalf("saving photos: %v", err)
	}

	post, err := repo.GetPostByID(context.Background(), persisted.PostID)
	if err != nil {
		t.Fatalf("getting post: %v", err)
	}
	if !post.HasImage {
		t.Fatalf("expected post.HasImage to be true")
	}
	if len(repo.photos) < 2 {
		t.Fatalf("expected photo rows to be stored")
	}
}

func TestInMemorySavePostPhotos_RejectsBlankS3Key(t *testing.T) {
	repo := NewInMemory()

	err := repo.SavePostPhotos(context.Background(), []domain.PostCreateSavedPhoto{
		{PostID: 123, S3Key: "   ", Position: 0},
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
