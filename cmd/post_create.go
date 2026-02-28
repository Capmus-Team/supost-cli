package cmd

import (
	"fmt"
	"strings"

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

		if postCreateSubmitRequested(cmd) {
			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return fmt.Errorf("reading name flag: %w", err)
			}
			body, err := cmd.Flags().GetString("body")
			if err != nil {
				return fmt.Errorf("reading body flag: %w", err)
			}
			email, err := cmd.Flags().GetString("email")
			if err != nil {
				return fmt.Errorf("reading email flag: %w", err)
			}
			price, err := cmd.Flags().GetFloat64("price")
			if err != nil {
				return fmt.Errorf("reading price flag: %w", err)
			}
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return fmt.Errorf("reading dry-run flag: %w", err)
			}

			input := domain.PostCreateSubmission{
				CategoryID:    categoryID,
				SubcategoryID: subcategoryID,
				Name:          strings.TrimSpace(name),
				Body:          strings.TrimSpace(body),
				Email:         strings.TrimSpace(email),
				Price:         price,
				PriceProvided: cmd.Flags().Changed("price"),
			}

			var sender service.PostCreateEmailSender
			if !dryRun {
				mailgunSender, err := adapters.NewMailgunSender(
					cfg.MailgunAPIBase,
					cfg.MailgunDomain,
					cfg.MailgunAPIKey,
					cfg.MailgunFromEmail,
					cfg.MailgunSendTimeout,
				)
				if err != nil {
					return fmt.Errorf("configuring mailgun sender: %w", err)
				}
				sender = mailgunSender
			}

			result, err := svc.Submit(
				cmd.Context(),
				input,
				dryRun,
				cfg.SupostBaseURL,
				cfg.MailgunFromEmail,
				sender,
			)
			if err != nil {
				return fmt.Errorf("submitting post: %w", err)
			}
			return renderPostCreateSubmitOutput(cmd, cfg.Format, result)
		}

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
	postCreateCmd.Flags().String("name", "", "post title")
	postCreateCmd.Flags().String("body", "", "post body")
	postCreateCmd.Flags().String("email", "", "poster email")
	postCreateCmd.Flags().Float64("price", 0, "post price")
	postCreateCmd.Flags().Bool("dry-run", false, "validate and render publish email without inserting/sending")
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

func renderPostCreateSubmitOutput(cmd *cobra.Command, format string, result domain.PostCreateSubmitResult) error {
	if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
		return adapters.RenderPostCreateSubmitResult(cmd.OutOrStdout(), result)
	}
	if format == "text" || format == "table" {
		return adapters.RenderPostCreateSubmitResult(cmd.OutOrStdout(), result)
	}
	return adapters.Render(format, result)
}

func postCreateSubmitRequested(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("name") ||
		cmd.Flags().Changed("body") ||
		cmd.Flags().Changed("email") ||
		cmd.Flags().Changed("price") ||
		cmd.Flags().Changed("dry-run")
}
