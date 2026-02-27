package config

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
	SupabaseURL    string `json:"supabase_url"`
	SupabaseAnonKey string `json:"supabase_anon_key"`
}

// Load reads configuration from viper (merges file + env + flags).
func Load() (*Config, error) {
	return &Config{
		Verbose:        viper.GetBool("verbose"),
		Format:         viper.GetString("format"),
		DatabaseURL:    viper.GetString("database_url"),
		Port:           viper.GetInt("port"),
		SupabaseURL:    viper.GetString("supabase_url"),
		SupabaseAnonKey: viper.GetString("supabase_anon_key"),
	}, nil
}
