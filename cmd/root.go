package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gsecutil",
	Short: "A Google Secret Manager utility CLI",
	Long: `gsecutil is a command-line utility that provides a simple wrapper
around the gcloud CLI for managing Google Secret Manager secrets.

It allows you to get, create, update, delete, list, and describe secrets
with simplified commands, and also provides the ability to copy secret
values directly to your clipboard.`,
	SilenceErrors: true, // Prevent duplicate error printing
	SilenceUsage:  true, // Don't show usage on every error
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("project", "p", "", "Google Cloud project ID")
	rootCmd.PersistentFlags().String("config", "", "Configuration file path (default: auto-detect: ./gsecutil.conf then $HOME/.config/gsecutil/gsecutil.conf)")

	// Set up pre-run hook to load custom config if specified
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		if configPath != "" {
			if err := SetCustomConfigPath(configPath); err != nil {
				return fmt.Errorf("failed to load config file: %w", err)
			}
		}
		return nil
	}
}
