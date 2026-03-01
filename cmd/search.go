package cmd

import (
	"fmt"

	"github.com/Capmus-Team/supost-cli/internal/adapters"
	"github.com/Capmus-Team/supost-cli/internal/config"
	"github.com/Capmus-Team/supost-cli/internal/domain"
	"github.com/Capmus-Team/supost-cli/internal/repository"
	"github.com/Capmus-Team/supost-cli/internal/service"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Render all-post search results",
	Long:  "Show paginated active posts grouped by posting date.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		categoryID, err := cmd.Flags().GetInt64("category")
		if err != nil {
			return fmt.Errorf("reading category flag: %w", err)
		}
		subcategoryID, err := cmd.Flags().GetInt64("subcategory")
		if err != nil {
			return fmt.Errorf("reading subcategory flag: %w", err)
		}
		page, err := cmd.Flags().GetInt("page")
		if err != nil {
			return fmt.Errorf("reading page flag: %w", err)
		}
		perPage, err := cmd.Flags().GetInt("per-page")
		if err != nil {
			return fmt.Errorf("reading per-page flag: %w", err)
		}

		var (
			repo      service.SearchRepository
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

		svc := service.NewSearchService(repo)
		result, err := svc.Search(cmd.Context(), categoryID, subcategoryID, page, perPage)
		if err != nil {
			return fmt.Errorf("fetching search results: %w", err)
		}

		return renderSearchOutput(cmd, cfg.Format, result)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Int64("category", 0, "filter by category id")
	searchCmd.Flags().Int64("subcategory", 0, "filter by subcategory id")
	searchCmd.Flags().Int("page", 1, "page number (1-based)")
	searchCmd.Flags().Int("per-page", 100, "posts per page (max 100)")
}

func renderSearchOutput(cmd *cobra.Command, format string, result domain.SearchResultPage) error {
	if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
		return adapters.RenderSearchResults(cmd.OutOrStdout(), result)
	}
	if format == "text" || format == "table" {
		return adapters.RenderSearchResults(cmd.OutOrStdout(), result)
	}
	return adapters.Render(format, result)
}
