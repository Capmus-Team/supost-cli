package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func (r *Postgres) SearchActivePosts(ctx context.Context, queryText string, categoryID, subcategoryID int64, page, perPage int) ([]domain.Post, bool, error) {
	queryText = strings.TrimSpace(queryText)
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}
	limit := perPage + 1
	offset := (page - 1) * perPage

	const sqlQuery = `
SELECT
	id,
	COALESCE(category_id, 0) AS category_id,
	COALESCE(subcategory_id, 0) AS subcategory_id,
	COALESCE(email, '') AS email,
	COALESCE(name, '') AS name,
	COALESCE(status, 0) AS status,
	COALESCE(time_posted, 0) AS time_posted,
	COALESCE(time_posted_at, to_timestamp(0)) AS time_posted_at,
	COALESCE(price::float8, 0) AS price,
	(price IS NOT NULL) AS has_price,
	(
		COALESCE(photo1_file_name, '') <> '' OR
		COALESCE(photo2_file_name, '') <> '' OR
		COALESCE(photo3_file_name, '') <> '' OR
		COALESCE(photo4_file_name, '') <> '' OR
		COALESCE(image_source1, '') <> '' OR
		COALESCE(image_source2, '') <> '' OR
		COALESCE(image_source3, '') <> '' OR
		COALESCE(image_source4, '') <> ''
	) AS has_image,
	COALESCE(created_at, now()) AS created_at,
	COALESCE(updated_at, created_at, now()) AS updated_at
FROM public.post
WHERE status = $1
  AND ($2 = 0 OR category_id = $2)
  AND ($3 = 0 OR subcategory_id = $3)
  AND (
    $4 = '' OR
    to_tsvector('simple', COALESCE(name, '') || ' ' || COALESCE(body, '')) @@ plainto_tsquery('simple', $4)
  )
ORDER BY time_posted DESC NULLS LAST, id DESC
LIMIT $5 OFFSET $6
`

	rows, err := r.db.QueryContext(ctx, sqlQuery, domain.PostStatusActive, categoryID, subcategoryID, queryText, limit, offset)
	if err != nil {
		return nil, false, fmt.Errorf("querying search posts: %w", err)
	}
	defer rows.Close()

	posts := make([]domain.Post, 0, limit)
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(
			&post.ID,
			&post.CategoryID,
			&post.SubcategoryID,
			&post.Email,
			&post.Name,
			&post.Status,
			&post.TimePosted,
			&post.TimePostedAt,
			&post.Price,
			&post.HasPrice,
			&post.HasImage,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, false, fmt.Errorf("scanning search row: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("iterating search rows: %w", err)
	}

	hasMore := len(posts) > perPage
	if hasMore {
		posts = posts[:perPage]
	}
	return posts, hasMore, nil
}
