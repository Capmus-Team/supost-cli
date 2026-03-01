package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const sqlQuerySearchDefault = `
SELECT
	p.id,
	COALESCE(p.category_id, 0) AS category_id,
	COALESCE(p.subcategory_id, 0) AS subcategory_id,
	COALESCE(p.email, '') AS email,
	COALESCE(p.name, '') AS name,
	COALESCE(p.status, 0) AS status,
	COALESCE(p.time_posted, 0) AS time_posted,
	COALESCE(p.time_posted_at, to_timestamp(0)) AS time_posted_at,
	COALESCE(p.price::float8, 0) AS price,
	(p.price IS NOT NULL) AS has_price,
	(
		COALESCE(p.photo1_file_name, '') <> '' OR
		COALESCE(p.photo2_file_name, '') <> '' OR
		COALESCE(p.photo3_file_name, '') <> '' OR
		COALESCE(p.photo4_file_name, '') <> '' OR
		COALESCE(p.image_source1, '') <> '' OR
		COALESCE(p.image_source2, '') <> '' OR
		COALESCE(p.image_source3, '') <> '' OR
		COALESCE(p.image_source4, '') <> ''
	) AS has_image,
	COALESCE(p.created_at, now()) AS created_at,
	COALESCE(p.updated_at, p.created_at, now()) AS updated_at
FROM public.post p
WHERE p.status = $1
  AND ($2 = 0 OR p.category_id = $2)
  AND ($3 = 0 OR p.subcategory_id = $3)
ORDER BY p.time_posted DESC NULLS LAST, p.id DESC
LIMIT $4 OFFSET $5
`

const sqlQuerySearchFTS = `
SELECT
	p.id,
	COALESCE(p.category_id, 0) AS category_id,
	COALESCE(p.subcategory_id, 0) AS subcategory_id,
	COALESCE(p.email, '') AS email,
	COALESCE(p.name, '') AS name,
	COALESCE(p.status, 0) AS status,
	COALESCE(p.time_posted, 0) AS time_posted,
	COALESCE(p.time_posted_at, to_timestamp(0)) AS time_posted_at,
	COALESCE(p.price::float8, 0) AS price,
	(p.price IS NOT NULL) AS has_price,
	(
		COALESCE(p.photo1_file_name, '') <> '' OR
		COALESCE(p.photo2_file_name, '') <> '' OR
		COALESCE(p.photo3_file_name, '') <> '' OR
		COALESCE(p.photo4_file_name, '') <> '' OR
		COALESCE(p.image_source1, '') <> '' OR
		COALESCE(p.image_source2, '') <> '' OR
		COALESCE(p.image_source3, '') <> '' OR
		COALESCE(p.image_source4, '') <> ''
	) AS has_image,
	COALESCE(p.created_at, now()) AS created_at,
	COALESCE(p.updated_at, p.created_at, now()) AS updated_at
FROM public.post p, plainto_tsquery('english', $4) q
WHERE p.status = $1
  AND ($2 = 0 OR p.category_id = $2)
  AND ($3 = 0 OR p.subcategory_id = $3)
  AND p.fts @@ q
ORDER BY ts_rank(p.fts, q) DESC, p.time_posted DESC NULLS LAST, p.id DESC
LIMIT $5 OFFSET $6
`

func (r *Postgres) SearchActivePosts(ctx context.Context, queryText string, categoryID, subcategoryID int64, page, perPage int) ([]domain.Post, bool, error) {
	querySQL, queryArgs, normalizedPerPage := buildSearchActivePostsStatement(queryText, categoryID, subcategoryID, page, perPage)

	rows, err := r.db.QueryContext(ctx, querySQL, queryArgs...)
	if err != nil {
		return nil, false, fmt.Errorf("querying search posts: %w", err)
	}
	defer rows.Close()

	posts := make([]domain.Post, 0, normalizedPerPage+1)
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

	hasMore := len(posts) > normalizedPerPage
	if hasMore {
		posts = posts[:normalizedPerPage]
	}
	return posts, hasMore, nil
}

func buildSearchActivePostsStatement(queryText string, categoryID, subcategoryID int64, page, perPage int) (string, []any, int) {
	queryText = strings.TrimSpace(queryText)
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}

	limit := perPage + 1
	offset := (page - 1) * perPage
	if queryText == "" {
		return sqlQuerySearchDefault, []any{domain.PostStatusActive, categoryID, subcategoryID, limit, offset}, perPage
	}

	return sqlQuerySearchFTS, []any{domain.PostStatusActive, categoryID, subcategoryID, queryText, limit, offset}, perPage
}
