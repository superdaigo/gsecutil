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
}
