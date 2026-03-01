package repository

import (
	"context"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *InMemory) CreateResponseMessage(_ context.Context, postID int64, replyToEmail, message, ip, userAgent string) (domain.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	id := r.nextMessageIDLocked()
	record := domain.Message{
		ID:        id,
		PostID:    postID,
		Message:   message,
		IP:        ip,
		Email:     replyToEmail,
		RawEmail:  replyToEmail,
		Source:    "cli",
		Status:    "queued",
		UserAgent: userAgent,
		Scammed:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.messages = append(r.messages, record)
	return record, nil
}

func (r *InMemory) nextMessageIDLocked() int64 {
	var maxID int64
	for _, message := range r.messages {
		if message.ID > maxID {
			maxID = message.ID
		}
	}
	if maxID <= 0 {
		return 1
	}
	return maxID + 1
}
