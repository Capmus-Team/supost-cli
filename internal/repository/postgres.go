package repository

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// Postgres implements repository methods backed by PostgreSQL.
type Postgres struct {
	databaseURL string
}

const maxRecentActivePosts = 50
const fieldSeparator = "\x1f"
const recordSeparator = "\x1e"

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

// ListRecentActivePosts returns status=1 posts sorted by newest first.
func (r *Postgres) ListRecentActivePosts(ctx context.Context, limit int) ([]domain.Post, error) {
	limit = clampRecentLimit(limit)
	query := buildRecentActivePostsQuery(limit)

	command := exec.CommandContext(
		ctx,
		"psql",
		r.databaseURL,
		"-X",
		"-q",
		"-A",
		"-t",
		"-F", fieldSeparator,
		"-R", recordSeparator,
		"-v", "ON_ERROR_STOP=1",
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

	posts, err := parseRecentActivePosts(stdout.String())
	if err != nil {
		return nil, fmt.Errorf("parsing post rows: %w", err)
	}
	if len(posts) > limit {
		return posts[:limit], nil
	}
	return posts, nil
}

func clampRecentLimit(limit int) int {
	if limit <= 0 || limit > maxRecentActivePosts {
		return maxRecentActivePosts
	}
	return limit
}

func buildRecentActivePostsQuery(limit int) string {
	return fmt.Sprintf(`
SELECT
	COALESCE(id, 0),
	COALESCE(email, ''),
	REPLACE(REPLACE(COALESCE(name, ''), CHR(30), ' '), CHR(31), ' '),
	COALESCE(status, 0),
	COALESCE(time_posted, 0),
	COALESCE(EXTRACT(EPOCH FROM time_posted_at)::bigint, 0),
	COALESCE(price::float8, 0),
	(price IS NOT NULL),
	(
		COALESCE(photo1_file_name, '') <> '' OR
		COALESCE(photo2_file_name, '') <> '' OR
		COALESCE(photo3_file_name, '') <> '' OR
		COALESCE(photo4_file_name, '') <> '' OR
		COALESCE(image_source1, '') <> '' OR
		COALESCE(image_source2, '') <> '' OR
		COALESCE(image_source3, '') <> '' OR
		COALESCE(image_source4, '') <> ''
	)
FROM public.post
WHERE status = 1
ORDER BY time_posted DESC NULLS LAST, id DESC
LIMIT %d
`, limit)
}

func parseRecentActivePosts(raw string) ([]domain.Post, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []domain.Post{}, nil
	}

	records := strings.Split(raw, recordSeparator)
	posts := make([]domain.Post, 0, len(records))
	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		fields := strings.Split(record, fieldSeparator)
		if len(fields) != 9 {
			return nil, fmt.Errorf("expected 9 fields, got %d", len(fields))
		}

		id, err := parseInt64Field(fields[0])
		if err != nil {
			return nil, fmt.Errorf("id: %w", err)
		}
		status64, err := parseInt64Field(fields[3])
		if err != nil {
			return nil, fmt.Errorf("status: %w", err)
		}
		timePosted, err := parseInt64Field(fields[4])
		if err != nil {
			return nil, fmt.Errorf("time_posted: %w", err)
		}
		timePostedAtUnix, err := parseInt64Field(fields[5])
		if err != nil {
			return nil, fmt.Errorf("time_posted_at: %w", err)
		}
		price, err := parseFloat64Field(fields[6])
		if err != nil {
			return nil, fmt.Errorf("price: %w", err)
		}
		hasPrice, err := parseBoolField(fields[7])
		if err != nil {
			return nil, fmt.Errorf("has_price: %w", err)
		}
		hasImage, err := parseBoolField(fields[8])
		if err != nil {
			return nil, fmt.Errorf("has_image: %w", err)
		}

		post := domain.Post{
			ID:         id,
			Email:      fields[1],
			Name:       fields[2],
			Status:     int(status64),
			TimePosted: timePosted,
			Price:      price,
			HasPrice:   hasPrice,
			HasImage:   hasImage,
		}

		if timePostedAtUnix > 0 {
			post.TimePostedAt = time.Unix(timePostedAtUnix, 0).UTC()
		}

		posts = append(posts, post)
	}

	if len(posts) > maxRecentActivePosts {
		posts = posts[:maxRecentActivePosts]
	}
	return posts, nil
}

func parseInt64Field(raw string) (int64, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, nil
	}
	out, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %w", raw, err)
	}
	return out, nil
}

func parseFloat64Field(raw string) (float64, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, nil
	}
	out, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", raw, err)
	}
	return out, nil
}

func parseBoolField(raw string) (bool, error) {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "t", "true", "1":
		return true, nil
	case "f", "false", "0", "":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool %q", raw)
	}
}
