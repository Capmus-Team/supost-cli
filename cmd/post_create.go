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

var postCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Start creating a new post",
	Long:  "Render staged post creation pages: choose category, choose subcategory, then form fields.",
	Args:  cobra.NoArgs,
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

		var (
			repo      service.PostCreateRepository
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

		svc := service.NewPostCreateService(repo)
		page, err := svc.BuildPage(cmd.Context(), categoryID, subcategoryID)
		if err != nil {
			return fmt.Errorf("building post create page: %w", err)
		}

		return renderPostCreateOutput(cmd, cfg.Format, page)
	},
}

func init() {
	postCmd.AddCommand(postCreateCmd)
	postCreateCmd.Flags().Int64("category", 0, "selected category id")
	postCreateCmd.Flags().Int64("subcategory", 0, "selected subcategory id")
}

func renderPostCreateOutput(cmd *cobra.Command, format string, page domain.PostCreatePage) error {
	if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
		return adapters.RenderPostCreatePage(cmd.OutOrStdout(), page)
	}
	if format == "text" || format == "table" {
		return adapters.RenderPostCreatePage(cmd.OutOrStdout(), page)
	}
	return adapters.Render(format, page)
}
