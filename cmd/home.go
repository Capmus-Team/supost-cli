package cmd

import (
	"fmt"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/adapters"
	"github.com/Capmus-Team/supost-cli/internal/config"
	"github.com/Capmus-Team/supost-cli/internal/domain"
	"github.com/Capmus-Team/supost-cli/internal/repository"
	"github.com/Capmus-Team/supost-cli/internal/service"
	"github.com/spf13/cobra"
)

var homeCmd = &cobra.Command{
	Use:   "home",
	Short: "Render the SUPost home feed",
	Long:  "Show recently posted active posts from the post table.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return fmt.Errorf("reading limit flag: %w", err)
		}
		cacheTTL, err := cmd.Flags().GetDuration("cache-ttl")
		if err != nil {
			return fmt.Errorf("reading cache-ttl flag: %w", err)
		}
		refresh, err := cmd.Flags().GetBool("refresh")
		if err != nil {
			return fmt.Errorf("reading refresh flag: %w", err)
		}

		if cfg.DatabaseURL != "" {
			cachedPosts, ok, err := getCachedHomePosts(cfg.DatabaseURL, refresh, cacheTTL, limit)
			if err == nil && ok {
				return renderHomeOutput(cmd, cfg.Format, cachedPosts)
			}
		}

		var (
			repo      service.HomeRepository
			closeRepo func() error
		)
		if cfg.DatabaseURL != "" {
			pgRepo, err := repository.NewPostgres(cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("connecting to postgres: %w", err)
			}
			repo = pgRepo
			closeRepo = pgRepo.Close
		} else {
			repo = repository.NewInMemory()
		}
		if closeRepo != nil {
			defer func() {
				_ = closeRepo()
			}()
		}

		svc := service.NewHomeService(repo)
		posts, err := svc.ListRecentActive(cmd.Context(), limit)
		if err != nil {
			return fmt.Errorf("fetching recent active posts: %w", err)
		}

		if cfg.DatabaseURL != "" && cacheTTL > 0 {
			_ = adapters.SaveHomePostsCache(posts)
		}

		return renderHomeOutput(cmd, cfg.Format, posts)
	},
}

func init() {
	rootCmd.AddCommand(homeCmd)
	homeCmd.Flags().Int("limit", 50, "number of recent active posts to show")
	homeCmd.Flags().Duration("cache-ttl", 30*time.Second, "cache TTL for home feed when using database")
	homeCmd.Flags().Bool("refresh", false, "bypass cache and fetch fresh data from database")
}

func getCachedHomePosts(databaseURL string, refresh bool, ttl time.Duration, limit int) ([]domain.Post, bool, error) {
	if databaseURL == "" || refresh || ttl <= 0 {
		return nil, false, nil
	}
	posts, ok, err := adapters.LoadHomePostsCache(ttl, limit)
	if err != nil {
		return nil, false, err
	}
	return posts, ok, nil
}

func renderHomeOutput(cmd *cobra.Command, format string, posts []domain.Post) error {
	if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
		return adapters.RenderHomePosts(cmd.OutOrStdout(), posts)
	}
	if format == "text" || format == "table" {
		return adapters.RenderHomePosts(cmd.OutOrStdout(), posts)
	}
	return adapters.Render(format, posts)
}
