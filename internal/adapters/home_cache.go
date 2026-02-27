package adapters

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type homePostsCache struct {
	CachedAt time.Time     `json:"cached_at"`
	Posts    []domain.Post `json:"posts"`
}

// LoadHomePostsCache reads recent posts from local cache when still valid.
func LoadHomePostsCache(ttl time.Duration, limit int) ([]domain.Post, bool, error) {
	if ttl <= 0 {
		return nil, false, nil
	}

	cachePath, err := homePostsCachePath()
	if err != nil {
		return nil, false, fmt.Errorf("resolving cache path: %w", err)
	}

	payload, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("reading cache file: %w", err)
	}

	var cache homePostsCache
	if err := json.Unmarshal(payload, &cache); err != nil {
		return nil, false, fmt.Errorf("decoding cache JSON: %w", err)
	}

	if cache.CachedAt.IsZero() || time.Since(cache.CachedAt) > ttl {
		return nil, false, nil
	}

	posts := cache.Posts
	if limit > 0 && len(posts) > limit {
		posts = posts[:limit]
	}
	return posts, true, nil
}

// SaveHomePostsCache stores recent posts to local cache.
func SaveHomePostsCache(posts []domain.Post) error {
	cachePath, err := homePostsCachePath()
	if err != nil {
		return fmt.Errorf("resolving cache path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return fmt.Errorf("creating cache directory: %w", err)
	}

	cache := homePostsCache{
		CachedAt: time.Now().UTC(),
		Posts:    posts,
	}

	payload, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("encoding cache JSON: %w", err)
	}

	if err := os.WriteFile(cachePath, payload, 0o644); err != nil {
		return fmt.Errorf("writing cache file: %w", err)
	}
	return nil
}

func homePostsCachePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	return filepath.Join(cacheDir, "supost-cli", "home_recent_active_posts.json"), nil
}
