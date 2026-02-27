package service

import (
	"context"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockHomeRepo struct {
	receivedLimit int
	posts         []domain.Post
}

func (m *mockHomeRepo) ListRecentActivePosts(_ context.Context, limit int) ([]domain.Post, error) {
	m.receivedLimit = limit
	return m.posts, nil
}

func TestHomeService_ListRecentActive_UsesDefaultLimit(t *testing.T) {
	repo := &mockHomeRepo{}
	svc := NewHomeService(repo)

	if _, err := svc.ListRecentActive(context.Background(), 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.receivedLimit != defaultHomeLimit {
		t.Fatalf("expected default limit %d, got %d", defaultHomeLimit, repo.receivedLimit)
	}
}

func TestHomeService_ListRecentActive_UsesProvidedLimit(t *testing.T) {
	repo := &mockHomeRepo{}
	svc := NewHomeService(repo)

	if _, err := svc.ListRecentActive(context.Background(), 12); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.receivedLimit != 12 {
		t.Fatalf("expected limit 12, got %d", repo.receivedLimit)
	}
}
