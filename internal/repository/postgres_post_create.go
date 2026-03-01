package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *Postgres) CreatePendingPost(ctx context.Context, submission domain.PostCreateSubmission) (domain.PostCreatePersisted, error) {
	postedAt := submission.PostedAt
	if postedAt.IsZero() {
		postedAt = time.Now()
	}
	postedUnix := postedAt.Unix()
	var ipValue any
	if submission.IP != "" {
		ipValue = submission.IP
	}
	var priceValue any
	if submission.PriceProvided {
		priceValue = submission.Price
	}

	const query = `
INSERT INTO public.post (
	college_id,
	category_id,
	subcategory_id,
	email,
	ip,
	name,
	body,
	status,
	time_posted,
	time_modified,
	time_posted_at,
	time_modified_at,
	access_token,
	price,
	created_at,
	updated_at
) VALUES (
	1,
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	0,
	$7,
	$7,
	to_timestamp($7),
	to_timestamp($7),
	$8,
	$9,
	now(),
	now()
)
RETURNING
	COALESCE(id, 0) AS id,
	COALESCE(access_token, '') AS access_token,
	COALESCE(time_posted_at, to_timestamp($7)) AS time_posted_at
`

	var persisted domain.PostCreatePersisted
	err := r.db.QueryRowContext(
		ctx,
		query,
		submission.CategoryID,
		submission.SubcategoryID,
		submission.Email,
		ipValue,
		submission.Name,
		submission.Body,
		postedUnix,
		submission.AccessToken,
		priceValue,
	).Scan(
		&persisted.PostID,
		&persisted.AccessToken,
		&persisted.PostedAt,
	)
	if err != nil {
		return domain.PostCreatePersisted{}, fmt.Errorf("inserting post: %w", err)
	}
	return persisted, nil
}
