package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is updated on every release. See AGENTS.md §11.
var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("v%s\n", Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
