package cmd

import (
	"fmt"
	"strings"

	"github.com/Capmus-Team/supost-cli/internal/adapters"
	"github.com/Capmus-Team/supost-cli/internal/config"
	"github.com/Capmus-Team/supost-cli/internal/domain"
	"github.com/Capmus-Team/supost-cli/internal/service"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new user account with Supabase Auth",
	Long:  "Sign up a user with display name, email, phone, and password via Supabase Auth.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		displayName, err := cmd.Flags().GetString("display-name")
		if err != nil {
			return fmt.Errorf("reading display-name flag: %w", err)
		}
		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return fmt.Errorf("reading email flag: %w", err)
		}
		phone, err := cmd.Flags().GetString("phone")
		if err != nil {
			return fmt.Errorf("reading phone flag: %w", err)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return fmt.Errorf("reading password flag: %w", err)
		}

		apiKey := strings.TrimSpace(cfg.SupabasePublishableKey)
		if apiKey == "" {
			apiKey = strings.TrimSpace(cfg.SupabaseAnonKey)
		}

		provider, err := adapters.NewSupabaseAuthSignupClient(cfg.SupabaseURL, apiKey, cfg.SupabaseSecretKey)
		if err != nil {
			return fmt.Errorf("configuring supabase auth signup: %w", err)
		}

		svc := service.NewUserSignupService(provider)
		result, err := svc.SignUp(cmd.Context(), domain.UserSignupSubmission{
			DisplayName: displayName,
			Email:       email,
			Phone:       phone,
			Password:    password,
		})
		if err != nil {
			return fmt.Errorf("signing up user: %w", err)
		}

		return adapters.Render(cfg.Format, result)
	},
}

func init() {
	rootCmd.AddCommand(signupCmd)
	signupCmd.Flags().String("display-name", "", "user display name")
	signupCmd.Flags().String("email", "", "user email")
	signupCmd.Flags().String("phone", "", "user phone in international format, e.g. +16505551234")
	signupCmd.Flags().String("password", "", "user password (minimum 8 characters)")
	_ = signupCmd.MarkFlagRequired("display-name")
	_ = signupCmd.MarkFlagRequired("email")
	_ = signupCmd.MarkFlagRequired("phone")
	_ = signupCmd.MarkFlagRequired("password")
}
