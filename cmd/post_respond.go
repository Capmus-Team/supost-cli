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

var postRespondCmd = &cobra.Command{
	Use:   "respond <post_id>",
	Short: "Send a response to a post owner",
	Long:  "Send a response email to the post owner and persist the message in app_private.message.",
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
		message, err := cmd.Flags().GetString("message")
		if err != nil {
			return fmt.Errorf("reading message flag: %w", err)
		}
		replyTo, err := cmd.Flags().GetString("reply-to")
		if err != nil {
			return fmt.Errorf("reading reply-to flag: %w", err)
		}
		ip, err := cmd.Flags().GetString("ip")
		if err != nil {
			return fmt.Errorf("reading ip flag: %w", err)
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return fmt.Errorf("reading dry-run flag: %w", err)
		}

		var (
			repo      service.PostRespondRepository
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

		var sender service.PostRespondEmailSender
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

		svc := service.NewPostRespondService(repo)
		result, err := svc.Respond(
			cmd.Context(),
			domain.PostRespondSubmission{
				PostID:    postID,
				Message:   strings.TrimSpace(message),
				ReplyTo:   strings.TrimSpace(replyTo),
				IP:        strings.TrimSpace(ip),
				UserAgent: "supost-cli",
			},
			dryRun,
			cfg.SupostBaseURL,
			cfg.MailgunFromEmail,
			sender,
		)
		if err != nil {
			return fmt.Errorf("responding to post %d: %w", postID, err)
		}
		return renderPostRespondOutput(cmd, cfg.Format, result)
	},
}

func init() {
	postCmd.AddCommand(postRespondCmd)
	postRespondCmd.Flags().String("message", "", "response message body")
	postRespondCmd.Flags().String("reply-to", "", "reply-to email")
	postRespondCmd.Flags().String("ip", "", "sender IP address (optional)")
	postRespondCmd.Flags().Bool("dry-run", false, "validate and render email without sending or persisting")
	_ = postRespondCmd.MarkFlagRequired("message")
	_ = postRespondCmd.MarkFlagRequired("reply-to")
}

func renderPostRespondOutput(cmd *cobra.Command, format string, result domain.PostRespondResult) error {
	if !cmd.Flags().Changed("format") && (format == "" || format == "json") {
		return adapters.RenderPostRespondResult(cmd.OutOrStdout(), result)
	}
	if format == "text" || format == "table" {
		return adapters.RenderPostRespondResult(cmd.OutOrStdout(), result)
	}
	return adapters.Render(format, result)
}
