package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configSetTitleCmd = &cobra.Command{
	Use:   "set-title SECRET_NAME TITLE",
	Short: "Set or update the title for a secret in the configuration file",
	Long: `Set or update the title for a secret in the configuration file.

This updates the local configuration file only and does not affect the
secret stored in Google Secret Manager.

Examples:
  gsecutil config set-title database-password "Production Database Password"
  gsecutil config set-title api-key "External API Key"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretName := args[0]
		title := args[1]

		// Add prefix if configured
		secretName = AddPrefixToSecretName(secretName)

		// Load or create config
		config, err := loadOrCreateConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Update the title
		updateConfigWithMetadata(config, secretName, title, nil)

		// Save config
		if err := saveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		// Determine which config file was updated
		configPath := configFilePath
		if configPath == "" {
			configPath = getDefaultConfigPath()
		}

		fmt.Printf("✓ Title updated for secret '%s'\n", secretName)
		fmt.Printf("  Configuration file: %s\n", configPath)

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetTitleCmd)
}
