package service

import (
	"context"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// PostRepository defines read-only post access where consumed.
type PostRepository interface {
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
}

// PostService orchestrates post page retrieval.
type PostService struct {
	repo PostRepository
}

// NewPostService constructs PostService.
func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

// GetByID fetches one post by identifier.
func (s *PostService) GetByID(ctx context.Context, postID int64) (domain.Post, error) {
	return s.repo.GetPostByID(ctx, postID)
}
