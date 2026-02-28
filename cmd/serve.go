package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Capmus-Team/supost-cli/internal/config"
	"github.com/Capmus-Team/supost-cli/internal/repository"
	"github.com/Capmus-Team/supost-cli/internal/service"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a preview HTTP server",
	Long: `Start a lightweight HTTP server that exposes the service layer as JSON
endpoints. This is for prototyping only — it will be replaced by Next.js
API routes in production. Uses in-memory data by default.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		port := cfg.Port
		if port == 0 {
			port = 8080
		}

		// Composition root: choose the repository adapter.
		repo := repository.NewInMemory()
		homeSvc := service.NewHomeService(repo)

		mux := http.NewServeMux()

		// GET /api/posts — returns recent active posts as JSON
		mux.HandleFunc("GET /api/posts", func(w http.ResponseWriter, r *http.Request) {
			posts, err := homeSvc.ListRecentActive(r.Context(), 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(posts)
		})

		// Health check
		mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		addr := fmt.Sprintf(":%d", port)
		log.Printf("Preview server running at http://localhost%s", addr)
		log.Printf("  GET /api/posts")
		log.Printf("  GET /api/health")
		log.Printf("Press Ctrl+C to stop.")
		return http.ListenAndServe(addr, mux)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "port to listen on")
}
