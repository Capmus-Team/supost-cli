package repository

import (
	"context"
	"database/sql"
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

// ListHomeCategorySections returns category/subcategory taxonomy and latest
// active post time per category for sidebar rendering.
func (r *Postgres) ListHomeCategorySections(ctx context.Context) ([]domain.HomeCategorySection, error) {
	const query = `
SELECT
	c.id AS category_id,
	COALESCE(NULLIF(c.short_name, ''), NULLIF(c.name, ''), '') AS category_name,
	COALESCE(s.name, '') AS subcategory_name,
	COALESCE(lat.latest_time_posted, 0) AS latest_time_posted
FROM public.category c
LEFT JOIN public.subcategory s ON s.category_id = c.id
LEFT JOIN LATERAL (
	SELECT COALESCE(p.time_posted, 0) AS latest_time_posted
	FROM public.post p
	WHERE p.status = $1 AND p.category_id = c.id
	ORDER BY p.time_posted DESC NULLS LAST, p.id DESC
	LIMIT 1
) lat ON true
ORDER BY c.id ASC, s.id ASC
`

	rows, err := r.db.QueryContext(ctx, query, domain.PostStatusActive)
	if err != nil {
		return nil, fmt.Errorf("querying home category sections: %w", err)
	}
	defer rows.Close()

	sectionsByID := make(map[int64]*domain.HomeCategorySection, 16)
	order := make([]int64, 0, 16)

	for rows.Next() {
		var (
			categoryID       int64
			categoryName     string
			subcategoryName  string
			latestTimePosted int64
		)
		if err := rows.Scan(&categoryID, &categoryName, &subcategoryName, &latestTimePosted); err != nil {
			return nil, fmt.Errorf("scanning home category section row: %w", err)
		}

		section, ok := sectionsByID[categoryID]
		if !ok {
			section = &domain.HomeCategorySection{
				CategoryID:   categoryID,
				CategoryName: strings.TrimSpace(categoryName),
			}
			if latestTimePosted > 0 {
				section.LastPostedAt = time.Unix(latestTimePosted, 0)
			}
			sectionsByID[categoryID] = section
			order = append(order, categoryID)
		}

		sub := strings.TrimSpace(subcategoryName)
		if sub != "" {
			section.SubcategoryNames = append(section.SubcategoryNames, sub)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating home category section rows: %w", err)
	}

	sections := make([]domain.HomeCategorySection, 0, len(order))
	for _, categoryID := range order {
		sections = append(sections, *sectionsByID[categoryID])
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
