package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "supost",
	Short: "supost — a CLI tool",
	Long:  `supost is a command-line application built with Cobra.`,
}

// Execute is called by main.go — the single entrypoint into the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".supost.yaml", "config file (default is .supost.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().String("format", "json", "output format: json, table, text")
	cobra.CheckErr(viper.BindPFlags(rootCmd.PersistentFlags()))
}

func initConfig() {
	_ = gotenv.Load(".env")

	if strings.TrimSpace(cfgFile) == "" {
		cfgFile = ".supost.yaml"
	}
	viper.SetConfigFile(cfgFile)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if verbose, _ := rootCmd.Flags().GetBool("verbose"); verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
