package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// Postgres implements repository methods backed by PostgreSQL.
type Postgres struct {
	databaseURL string
}

const maxRecentActivePosts = 50

// NewPostgres validates repository configuration.
func NewPostgres(databaseURL string) (*Postgres, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, fmt.Errorf("database_url is required")
	}
	return &Postgres{databaseURL: databaseURL}, nil
}

// Close is a no-op because this adapter shells out to psql per request.
func (r *Postgres) Close() error {
	return nil
}

type postQueryRow struct {
	ID            int64     `json:"id"`
	CategoryID    int64     `json:"category_id"`
	SubcategoryID int64     `json:"subcategory_id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	Status        int       `json:"status"`
	TimePosted    int64     `json:"time_posted"`
	TimePostedAt  time.Time `json:"time_posted_at"`
	Price         float64   `json:"price"`
	HasPrice      bool      `json:"has_price"`
	HasImage      bool      `json:"has_image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListRecentActivePosts returns status=1 posts sorted by newest first.
func (r *Postgres) ListRecentActivePosts(ctx context.Context, limit int) ([]domain.Post, error) {
	const query = `
SELECT
	row_to_json(t)::text
FROM (
	SELECT
		id,
		COALESCE(category_id, 0) AS category_id,
		COALESCE(subcategory_id, 0) AS subcategory_id,
		email,
		COALESCE(name, '') AS name,
		status,
		COALESCE(time_posted, 0) AS time_posted,
		time_posted_at,
		COALESCE(price, 0)::float8 AS price,
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
	WHERE status = 1
	ORDER BY time_posted DESC NULLS LAST, id DESC
	LIMIT 50
) t
`

	command := exec.CommandContext(
		ctx,
		"psql",
		r.databaseURL,
		"-v", "ON_ERROR_STOP=1",
		"-At",
		"-c", query,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		errText := strings.TrimSpace(stderr.String())
		if errText == "" {
			errText = "unknown psql error"
		}
		return nil, fmt.Errorf("running post query: %w (%s)", err, errText)
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "" {
		return []domain.Post{}, nil
	}

	posts := make([]domain.Post, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var row postQueryRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, fmt.Errorf("decoding post row: %w", err)
		}

		post := domain.Post{
			ID:            row.ID,
			CategoryID:    row.CategoryID,
			SubcategoryID: row.SubcategoryID,
			Email:         row.Email,
			Name:          row.Name,
			Status:        row.Status,
			TimePosted:    row.TimePosted,
			TimePostedAt:  row.TimePostedAt,
			Price:         row.Price,
			HasPrice:      row.HasPrice,
			HasImage:      row.HasImage,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}

		posts = append(posts, post)
	}

	if limit > 0 && len(posts) > limit {
		return posts[:limit], nil
	}
	if len(posts) > maxRecentActivePosts {
		return posts[:maxRecentActivePosts], nil
	}
	return posts, nil
}
