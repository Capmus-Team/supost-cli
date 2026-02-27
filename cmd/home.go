package cmd

import (
	"fmt"

	"github.com/Capmus-Team/supost-cli/internal/adapters"
	"github.com/Capmus-Team/supost-cli/internal/config"
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

		format := cfg.Format
		if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
			return adapters.RenderHomePosts(cmd.OutOrStdout(), posts)
		}
		if format == "text" || format == "table" {
			return adapters.RenderHomePosts(cmd.OutOrStdout(), posts)
		}
		return adapters.Render(format, posts)
	},
}

func init() {
	rootCmd.AddCommand(homeCmd)
	homeCmd.Flags().Int("limit", 50, "number of recent active posts to show")
}
