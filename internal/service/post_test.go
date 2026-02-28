package service

import (
	"context"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockPostRepo struct {
	post domain.Post
}

func (m *mockPostRepo) GetPostByID(_ context.Context, _ int64) (domain.Post, error) {
	return m.post, nil
}

func TestPostService_GetByID(t *testing.T) {
	repo := &mockPostRepo{
		post: domain.Post{ID: 130031961, Name: "Shared House"},
	}
	svc := NewPostService(repo)

	post, err := svc.GetByID(context.Background(), 130031961)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != 130031961 {
		t.Fatalf("expected post id 130031961, got %d", post.ID)
	}
}
