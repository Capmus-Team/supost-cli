package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/adapters"
	"github.com/Capmus-Team/supost-cli/internal/config"
	"github.com/Capmus-Team/supost-cli/internal/domain"
	"github.com/Capmus-Team/supost-cli/internal/repository"
	"github.com/Capmus-Team/supost-cli/internal/service"
	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post <post_id>",
	Short: "View a single post",
	Long:  "Render a single post page with header/footer, photos, body, and message poster panel.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		postID, err := parsePostIDArg(args[0])
		if err != nil {
			return err
		}

		var (
			repo      service.PostRepository
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

		svc := service.NewPostService(repo)
		post, err := svc.GetByID(cmd.Context(), postID)
		if err != nil {
			if err == domain.ErrNotFound {
				return fmt.Errorf("post %d not found", postID)
			}
			return fmt.Errorf("fetching post %d: %w", postID, err)
		}

		return renderPostOutput(cmd, cfg.Format, post)
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
}

func renderPostOutput(cmd *cobra.Command, format string, post domain.Post) error {
	if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
		return adapters.RenderPostPage(cmd.OutOrStdout(), post)
	}
	if format == "text" || format == "table" {
		return adapters.RenderPostPage(cmd.OutOrStdout(), post)
	}
	return adapters.Render(format, post)
}

func parsePostIDArg(raw string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid post id %q", raw)
	}
	return id, nil
}
