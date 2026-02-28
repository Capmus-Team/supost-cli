package cmd

import (
	"fmt"

	"github.com/Capmus-Team/supost-cli/internal/adapters"
	"github.com/Capmus-Team/supost-cli/internal/config"
	"github.com/Capmus-Team/supost-cli/internal/repository"
	"github.com/Capmus-Team/supost-cli/internal/service"
	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List categories and subcategories",
	Long:  "Display all categories with nested subcategories.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		var (
			repo      service.CategoryRepository
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

		svc := service.NewCategoryService(repo)
		categories, err := svc.ListCategoriesWithSubcategories(cmd.Context())
		if err != nil {
			return fmt.Errorf("listing categories: %w", err)
		}

		return adapters.Render(cfg.Format, categories)
	},
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}
