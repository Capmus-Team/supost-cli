package repository

import (
	"context"
	"fmt"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *Postgres) CreateResponseMessage(ctx context.Context, postID int64, replyToEmail, message, ip, userAgent string) (domain.Message, error) {
	const query = `
INSERT INTO app_private.message (
	message,
	post_id,
	ip,
	email,
	raw_email,
	source,
	status,
	user_agent,
	scammed,
	created_at,
	updated_at
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$4,
	'cli',
	'queued',
	NULLIF($5, ''),
	false,
	now(),
	now()
)
RETURNING
	COALESCE(id, 0) AS id,
	COALESCE(post_id, 0) AS post_id,
	COALESCE(message, '') AS message,
	COALESCE(ip::text, '') AS ip,
	COALESCE(email::text, '') AS email,
	COALESCE(raw_email::text, '') AS raw_email,
	COALESCE(source, '') AS source,
	COALESCE(status, '') AS status,
	COALESCE(user_agent, '') AS user_agent,
	COALESCE(scammed, false) AS scammed,
	COALESCE(created_at, now()) AS created_at,
	COALESCE(updated_at, created_at, now()) AS updated_at
`

	var out domain.Message
	err := r.db.QueryRowContext(ctx, query, message, postID, nullIfEmpty(ip), replyToEmail, userAgent).Scan(
		&out.ID,
		&out.PostID,
		&out.Message,
		&out.IP,
		&out.Email,
		&out.RawEmail,
		&out.Source,
		&out.Status,
		&out.UserAgent,
		&out.Scammed,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return domain.Message{}, fmt.Errorf("inserting response message: %w", err)
	}
	return out, nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}
