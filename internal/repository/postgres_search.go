package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const sqlQuerySearchSelect = `
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
		COALESCE(p.image_source4, '') <> '' OR
		EXISTS (SELECT 1 FROM public.photo ph WHERE ph.post_id = p.id)
	) AS has_image,
	COALESCE(p.created_at, now()) AS created_at,
	COALESCE(p.updated_at, p.created_at, now()) AS updated_at
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

	args := make([]any, 0, 6)
	args = append(args, domain.PostStatusActive)
	whereClauses := []string{fmt.Sprintf("p.status = $%d", len(args))}

	if categoryID > 0 {
		args = append(args, categoryID)
		whereClauses = append(whereClauses, fmt.Sprintf("p.category_id = $%d", len(args)))
	}
	if subcategoryID > 0 {
		args = append(args, subcategoryID)
		whereClauses = append(whereClauses, fmt.Sprintf("p.subcategory_id = $%d", len(args)))
	}

	orderBy := "p.time_posted DESC, p.id DESC"
	fromClause := "FROM public.post p"
	if queryText != "" {
		args = append(args, queryText)
		queryPos := len(args)
		fromClause = fmt.Sprintf("FROM public.post p\nCROSS JOIN plainto_tsquery('english', $%d) q", queryPos)
		whereClauses = append(whereClauses, "p.fts @@ q")
		orderBy = "ts_rank(p.fts, q) DESC, p.time_posted DESC, p.id DESC"
	}

	limit := perPage + 1
	offset := (page - 1) * perPage
	args = append(args, limit)
	limitPos := len(args)
	args = append(args, offset)
	offsetPos := len(args)

	querySQL := sqlQuerySearchSelect +
		"\n" + fromClause +
		"\nWHERE " + strings.Join(whereClauses, "\n  AND ") +
		"\nORDER BY " + orderBy +
		fmt.Sprintf("\nLIMIT $%d OFFSET $%d\n", limitPos, offsetPos)

	return querySQL, args, perPage
}
