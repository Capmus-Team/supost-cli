package config

import "time"

import "github.com/spf13/viper"

// Config holds all application configuration.
// All config loads through this package. No os.Getenv() elsewhere.
// See AGENTS.md §5.2.
type Config struct {
	Verbose     bool   `json:"verbose"`
	Format      string `json:"format"`
	DatabaseURL string `json:"database_url"` // postgresql://user:pass@host:port/dbname
	Port        int    `json:"port"`

	// Supabase (used by future Next.js frontend — shared .env)
	SupabaseURL            string `json:"supabase_url"`
	SupabaseAnonKey        string `json:"supabase_anon_key"`
	SupabasePublishableKey string `json:"supabase_publishable_key"`
	SupabaseSecretKey      string `json:"supabase_secret_key"`

	// Mailgun + publish-link URL (used by post create submit flow)
	MailgunDomain      string        `json:"mailgun_domain"`
	MailgunAPIKey      string        `json:"mailgun_api_key"`
	MailgunFromEmail   string        `json:"mailgun_from_email"`
	MailgunAPIBase     string        `json:"mailgun_api_base"`
	MailgunSendTimeout time.Duration `json:"mailgun_send_timeout"`
	SupostBaseURL      string        `json:"supost_base_url"`

	// S3 photo upload settings (used by post create when --photo is provided)
	S3PhotoBucket     string `json:"s3_photo_bucket"`
	S3PhotoPrefix     string `json:"s3_photo_prefix"`
	S3PhotoRegion     string `json:"s3_photo_region"`
	S3PhotoAWSProfile string `json:"s3_photo_aws_profile"`
}

// Load reads configuration from viper (merges file + env + flags).
func Load() (*Config, error) {
	supabaseSecretKey := viper.GetString("supabase_secret_key")
	if supabaseSecretKey == "" {
		// Backward-compatible alias used in many Supabase examples/docs.
		supabaseSecretKey = viper.GetString("supabase_service_role_key")
	}

	return &Config{
		Verbose:                viper.GetBool("verbose"),
		Format:                 viper.GetString("format"),
		DatabaseURL:            viper.GetString("database_url"),
		Port:                   viper.GetInt("port"),
		SupabaseURL:            viper.GetString("supabase_url"),
		SupabaseAnonKey:        viper.GetString("supabase_anon_key"),
		SupabasePublishableKey: viper.GetString("supabase_publishable_key"),
		SupabaseSecretKey:      supabaseSecretKey,
		MailgunDomain:          viper.GetString("mailgun_domain"),
		MailgunAPIKey:          viper.GetString("mailgun_api_key"),
		MailgunFromEmail:       viper.GetString("mailgun_from_email"),
		MailgunAPIBase:         viper.GetString("mailgun_api_base"),
		MailgunSendTimeout:     viper.GetDuration("mailgun_send_timeout"),
		SupostBaseURL:          viper.GetString("supost_base_url"),
		S3PhotoBucket:          viper.GetString("s3_photo_bucket"),
		S3PhotoPrefix:          viper.GetString("s3_photo_prefix"),
		S3PhotoRegion:          viper.GetString("s3_photo_region"),
		S3PhotoAWSProfile:      viper.GetString("s3_photo_aws_profile"),
	}, nil
}
