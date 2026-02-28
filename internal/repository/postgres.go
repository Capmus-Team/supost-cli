package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const maxRecentActivePosts = 50

// Postgres implements repository methods backed by PostgreSQL.
type Postgres struct {
	db *sql.DB
}

// NewPostgres initializes the PostgreSQL adapter.
func NewPostgres(databaseURL string) (*Postgres, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, fmt.Errorf("database_url is required")
	}

	connString := ensurePoolerSafeConnectionString(databaseURL)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Postgres{db: db}, nil
}

// Close closes the DB pool.
func (r *Postgres) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

// ListRecentActivePosts returns status=1 posts sorted by newest first.
func (r *Postgres) ListRecentActivePosts(ctx context.Context, limit int) ([]domain.Post, error) {
	limit = clampRecentLimit(limit)

	const query = `
SELECT
	id,
	COALESCE(category_id, 0) AS category_id,
	COALESCE(subcategory_id, 0) AS subcategory_id,
	COALESCE(email, '') AS email,
	COALESCE(name, '') AS name,
	COALESCE(status, 0) AS status,
	COALESCE(time_posted, 0) AS time_posted,
	COALESCE(time_modified, 0) AS time_modified,
	COALESCE(time_posted_at, to_timestamp(0)) AS time_posted_at,
	COALESCE(time_modified_at, to_timestamp(0)) AS time_modified_at,
	COALESCE(access_token, '') AS access_token,
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
ORDER BY time_posted DESC NULLS LAST, id DESC
LIMIT $2
`

	rows, err := r.db.QueryContext(ctx, query, domain.PostStatusActive, limit)
	if err != nil {
		return nil, fmt.Errorf("querying recent active posts: %w", err)
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
			&post.TimeModified,
			&post.TimePostedAt,
			&post.TimeModifiedAt,
			&post.AccessToken,
			&post.Price,
			&post.HasPrice,
			&post.HasImage,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning post row: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating post rows: %w", err)
	}

	return posts, nil
}

// GetPostByID returns one post by identifier.
func (r *Postgres) GetPostByID(ctx context.Context, postID int64) (domain.Post, error) {
	const query = `
SELECT
	id,
	COALESCE(category_id, 0) AS category_id,
	COALESCE(subcategory_id, 0) AS subcategory_id,
	COALESCE(email, '') AS email,
	COALESCE(name, '') AS name,
	COALESCE(body, '') AS body,
	COALESCE(photo1_file_name, '') AS photo1_file_name,
	COALESCE(photo2_file_name, '') AS photo2_file_name,
	COALESCE(photo3_file_name, '') AS photo3_file_name,
	COALESCE(photo4_file_name, '') AS photo4_file_name,
	COALESCE(image_source1, '') AS image_source1,
	COALESCE(image_source2, '') AS image_source2,
	COALESCE(image_source3, '') AS image_source3,
	COALESCE(image_source4, '') AS image_source4,
	COALESCE(status, 0) AS status,
	COALESCE(time_posted, 0) AS time_posted,
	COALESCE(time_modified, 0) AS time_modified,
	COALESCE(time_posted_at, to_timestamp(0)) AS time_posted_at,
	COALESCE(time_modified_at, to_timestamp(0)) AS time_modified_at,
	COALESCE(access_token, '') AS access_token,
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
WHERE id = $1
`

	var post domain.Post
	err := r.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.CategoryID,
		&post.SubcategoryID,
		&post.Email,
		&post.Name,
		&post.Body,
		&post.Photo1File,
		&post.Photo2File,
		&post.Photo3File,
		&post.Photo4File,
		&post.ImageSource1,
		&post.ImageSource2,
		&post.ImageSource3,
		&post.ImageSource4,
		&post.Status,
		&post.TimePosted,
		&post.TimeModified,
		&post.TimePostedAt,
		&post.TimeModifiedAt,
		&post.AccessToken,
		&post.Price,
		&post.HasPrice,
		&post.HasImage,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Post{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Post{}, fmt.Errorf("querying post by id: %w", err)
	}

	return post, nil
}

// ListRecentActivePostsByCategory returns active posts in one category sorted by newest first.
func (r *Postgres) ListRecentActivePostsByCategory(ctx context.Context, categoryID int64, limit int) ([]domain.Post, error) {
	limit = clampRecentLimit(limit)

	const query = `
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
  AND category_id = $2
ORDER BY time_posted DESC NULLS LAST, id DESC
LIMIT $3
`

	rows, err := r.db.QueryContext(ctx, query, domain.PostStatusActive, categoryID, limit)
	if err != nil {
		return nil, fmt.Errorf("querying recent active posts by category: %w", err)
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
			return nil, fmt.Errorf("scanning post row: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating post rows: %w", err)
	}

	return posts, nil
}

// ListCategories returns categories from public.category sorted by id.
func (r *Postgres) ListCategories(ctx context.Context) ([]domain.Category, error) {
	const query = `
SELECT
	COALESCE(id, 0) AS id,
	COALESCE(name, '') AS name,
	COALESCE(short_name, '') AS short_name
FROM public.category
ORDER BY id ASC
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying categories: %w", err)
	}
	defer rows.Close()

	categories := make([]domain.Category, 0, 32)
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.ShortName); err != nil {
			return nil, fmt.Errorf("scanning category row: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating category rows: %w", err)
	}

	return categories, nil
}

// ListSubcategories returns subcategories from public.subcategory sorted by category/id.
func (r *Postgres) ListSubcategories(ctx context.Context) ([]domain.Subcategory, error) {
	const query = `
SELECT
	COALESCE(id, 0) AS id,
	COALESCE(category_id, 0) AS category_id,
	COALESCE(name, '') AS name
FROM public.subcategory
ORDER BY category_id ASC, id ASC
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying subcategories: %w", err)
	}
	defer rows.Close()

	subcategories := make([]domain.Subcategory, 0, 128)
	for rows.Next() {
		var subcategory domain.Subcategory
		if err := rows.Scan(&subcategory.ID, &subcategory.CategoryID, &subcategory.Name); err != nil {
			return nil, fmt.Errorf("scanning subcategory row: %w", err)
		}
		subcategories = append(subcategories, subcategory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating subcategory rows: %w", err)
	}

	return subcategories, nil
}

// ListHomeCategorySections returns latest active post time per category.
// Taxonomy (category/subcategory names) is loaded from cached/local data.
func (r *Postgres) ListHomeCategorySections(ctx context.Context) ([]domain.HomeCategorySection, error) {
	const query = `
SELECT
	COALESCE(category_id, 0) AS category_id,
	MAX(
		COALESCE(
			time_posted_at,
			CASE WHEN COALESCE(time_posted, 0) > 0 THEN to_timestamp(time_posted) END
		)
	) AS last_posted_at
FROM public.post
WHERE status = $1
  AND category_id IS NOT NULL
GROUP BY category_id
ORDER BY category_id ASC
`

	rows, err := r.db.QueryContext(ctx, query, domain.PostStatusActive)
	if err != nil {
		return nil, fmt.Errorf("querying home category sections: %w", err)
	}
	defer rows.Close()

	sections := make([]domain.HomeCategorySection, 0, 16)

	for rows.Next() {
		var (
			categoryID   int64
			lastPostedAt sql.NullTime
		)
		if err := rows.Scan(&categoryID, &lastPostedAt); err != nil {
			return nil, fmt.Errorf("scanning home category section row: %w", err)
		}

		section := domain.HomeCategorySection{
			CategoryID: categoryID,
		}
		if lastPostedAt.Valid {
			section.LastPostedAt = lastPostedAt.Time
		}
		sections = append(sections, section)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating home category section rows: %w", err)
	}
	return sections, nil
}

func clampRecentLimit(limit int) int {
	if limit <= 0 || limit > maxRecentActivePosts {
		return maxRecentActivePosts
	}
	return limit
}

func ensurePoolerSafeConnectionString(databaseURL string) string {
	trimmed := strings.TrimSpace(databaseURL)
	if trimmed == "" {
		return trimmed
	}

	// URL-style DSN (postgres://... or postgresql://...)
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Scheme != "" {
		query := parsed.Query()
		if query.Get("default_query_exec_mode") == "" {
			query.Set("default_query_exec_mode", "simple_protocol")
		}
		if query.Get("statement_cache_capacity") == "" {
			query.Set("statement_cache_capacity", "0")
		}
		if query.Get("description_cache_capacity") == "" {
			query.Set("description_cache_capacity", "0")
		}
		parsed.RawQuery = query.Encode()
		return parsed.String()
	}

	// Key/value DSN fallback.
	if !strings.Contains(trimmed, "default_query_exec_mode=") {
		trimmed += " default_query_exec_mode=simple_protocol"
	}
	if !strings.Contains(trimmed, "statement_cache_capacity=") {
		trimmed += " statement_cache_capacity=0"
	}
	if !strings.Contains(trimmed, "description_cache_capacity=") {
		trimmed += " description_cache_capacity=0"
	}
	return trimmed
}
