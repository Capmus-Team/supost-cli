package repository

import (
	"context"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *InMemory) CreatePendingPost(_ context.Context, submission domain.PostCreateSubmission) (domain.PostCreatePersisted, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := submission.PostedAt
	if now.IsZero() {
		now = time.Now()
	}

	postID := r.nextPostIDLocked()
	r.posts = append(r.posts, domain.Post{
		ID:             postID,
		CategoryID:     submission.CategoryID,
		SubcategoryID:  submission.SubcategoryID,
		Email:          submission.Email,
		Name:           submission.Name,
		Body:           submission.Body,
		Status:         0,
		AccessToken:    submission.AccessToken,
		TimePosted:     now.Unix(),
		TimeModified:   now.Unix(),
		TimePostedAt:   now,
		TimeModifiedAt: now,
		Price:          submission.Price,
		HasPrice:       submission.PriceProvided,
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	return domain.PostCreatePersisted{
		PostID:      postID,
		AccessToken: submission.AccessToken,
		PostedAt:    now,
	}, nil
}

func (r *InMemory) nextPostIDLocked() int64 {
	var maxID int64
	for _, post := range r.posts {
		if post.ID > maxID {
			maxID = post.ID
		}
	}
	if maxID <= 0 {
		return 130000001
	}
	return maxID + 1
}
