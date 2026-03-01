package repository

import (
	"context"
	"fmt"
	"strings"
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
		IP:             submission.IP,
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

func (r *InMemory) SavePostPhotos(_ context.Context, photos []domain.PostCreateSavedPhoto) error {
	if len(photos) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, photo := range photos {
		if photo.PostID <= 0 {
			return fmt.Errorf("post_id must be positive")
		}
		s3Key := strings.TrimSpace(photo.S3Key)
		if s3Key == "" {
			return fmt.Errorf("s3_key is required")
		}
		r.photos = append(r.photos, domain.PostCreateSavedPhoto{
			PostID:      photo.PostID,
			S3Key:       s3Key,
			TickerS3Key: strings.TrimSpace(photo.TickerS3Key),
			Position:    photo.Position,
		})

		for idx, post := range r.posts {
			if post.ID != photo.PostID {
				continue
			}
			post.HasImage = true
			r.posts[idx] = post
			break
		}
	}
	return nil
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
