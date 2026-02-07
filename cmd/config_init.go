package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gsecutil configuration",
	Long:  `Manage gsecutil configuration file and settings.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gsecutil configuration file interactively",
	Long: `Initialize gsecutil configuration file interactively.

This command creates a new configuration file with guided prompts for:
- Google Cloud project ID
- Secret name prefix (for team organization)
- Default list attributes to display
- Example credential entries

The configuration file will be created at the default location unless
--output is specified.`,
	Example: `  gsecutil config init
  gsecutil config init --output /path/to/custom/config.yaml
  gsecutil config init --force  # Overwrite existing config`,
	RunE: runConfigInit,
}

var (
	configInitOutput string
	configInitForce  bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configInitCmd.Flags().StringVarP(&configInitOutput, "output", "o", "", "Output path for configuration file (default: $HOME/.config/gsecutil/gsecutil.conf)")
	configInitCmd.Flags().BoolVarP(&configInitForce, "force", "f", false, "Overwrite existing configuration file")
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	// Determine output path
	outputPath := configInitOutput
	if outputPath == "" {
		outputPath = getDefaultConfigPath()
	}

	// Check if file already exists
	if _, err := os.Stat(outputPath); err == nil && !configInitForce {
		return fmt.Errorf("configuration file already exists at '%s'. Use --force to overwrite", outputPath)
	}

	fmt.Println("Welcome to gsecutil configuration setup!")
	fmt.Println("This will guide you through creating a configuration file.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	config := &Config{}

	// Project ID
	fmt.Print("Google Cloud Project ID (optional, press Enter to skip): ")
	projectInput, _ := reader.ReadString('\n')
	config.Project = strings.TrimSpace(projectInput)

	// Prefix
	fmt.Println()
	fmt.Println("Secret name prefix helps organize secrets for teams.")
	fmt.Println("Example: 'team-shared-' will make 'database-password' become 'team-shared-database-password'")
	fmt.Print("Secret name prefix (optional, press Enter to skip): ")
	prefixInput, _ := reader.ReadString('\n')
	config.Prefix = strings.TrimSpace(prefixInput)

	// List attributes
	fmt.Println()
	fmt.Println("Default attributes to display in 'list' command.")
	fmt.Println("Common attributes: title, owner, environment, description")
	fmt.Print("Default list attributes (comma-separated, press Enter for 'title,owner,environment'): ")
	attributesInput, _ := reader.ReadString('\n')
	attributesInput = strings.TrimSpace(attributesInput)

	if attributesInput == "" {
		config.List.Attributes = []string{"title", "owner", "environment"}
	} else {
		config.List.Attributes = ParseShowAttributes(attributesInput)
	}

	// Ask if they want to add example credentials
	fmt.Println()
	fmt.Print("Add example credential entries? (y/N): ")
	addExamplesInput, _ := reader.ReadString('\n')
	addExamples := strings.ToLower(strings.TrimSpace(addExamplesInput))

	if addExamples == "y" || addExamples == "yes" {
		config.Credentials = []CredentialInfo{
			{
				Name:  "database-password",
				Title: "Production Database Password",
				Attributes: map[string]interface{}{
					"description": "MySQL root password for production database",
					"environment": "production",
					"owner":       "backend-team",
					"rotation":    "quarterly",
				},
			},
			{
				Name:  "api-key",
				Title: "External API Key",
				Attributes: map[string]interface{}{
					"description": "Production API key for payment processing",
					"environment": "production",
					"owner":       "api-team",
					"sensitive":   "high",
				},
			},
		}
	}

	// Create directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	fmt.Println()
	fmt.Printf("✓ Configuration file created successfully at: %s\n", outputPath)
	fmt.Println()
	fmt.Println("You can now:")
	fmt.Println("  - Edit the file manually to customize settings")
	fmt.Println("  - Run 'gsecutil list' to see your configuration in action")
	fmt.Println("  - Use 'gsecutil --config <path>' to use a different config file")
	fmt.Println()

	return nil
}
