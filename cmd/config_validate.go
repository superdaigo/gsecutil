package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configValidateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate gsecutil configuration file",
	Long: `Validate a gsecutil configuration file for correctness.

This command checks the configuration file for:
- Valid YAML syntax
- Valid configuration structure
- Prefix format (no spaces)
- No duplicate credential names
- No empty credential names
- Valid attribute references

If no file path is provided, validates the default configuration file.`,
	Example: `  gsecutil config validate
  gsecutil config validate /path/to/config.yaml
  gsecutil config validate --verbose  # Show detailed validation results`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConfigValidate,
}

var (
	configValidateVerbose bool
)

func init() {
	configCmd.AddCommand(configValidateCmd)
	configValidateCmd.Flags().BoolVarP(&configValidateVerbose, "verbose", "v", false, "Show detailed validation results")
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	// Determine config file path
	var configPath string
	if len(args) > 0 {
		configPath = args[0]
	} else {
		configPath = getDefaultConfigPath()
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file does not exist: %s", configPath)
	}

	fmt.Printf("Validating configuration file: %s\n", configPath)
	fmt.Println()

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Println("❌ YAML Syntax: FAILED")
		return fmt.Errorf("invalid YAML syntax: %w", err)
	}
	fmt.Println("✓ YAML Syntax: OK")

	// Validate configuration
	validationErrors := []string{}

	// Validate prefix
	if config.Prefix != "" {
		hasSpace := false
		for _, char := range config.Prefix {
			if char == ' ' {
				hasSpace = true
				break
			}
		}
		if hasSpace {
			validationErrors = append(validationErrors, fmt.Sprintf("prefix contains spaces: '%s'", config.Prefix))
		}
	}

	// Validate credentials
	if len(config.Credentials) > 0 {
		seenNames := make(map[string]bool)
		emptyNames := 0

		for i, cred := range config.Credentials {
			// Check for empty names
			if cred.Name == "" {
				emptyNames++
				validationErrors = append(validationErrors, fmt.Sprintf("credential at index %d has empty name", i))
			} else {
				// Check for duplicate names
				if seenNames[cred.Name] {
					validationErrors = append(validationErrors, fmt.Sprintf("duplicate credential name: '%s'", cred.Name))
				}
				seenNames[cred.Name] = true
			}
		}
	}

	// Validate list attributes reference existing credential fields
	if len(config.List.Attributes) > 0 && len(config.Credentials) > 0 {
		// Collect all available attributes from credentials
		availableAttrs := make(map[string]bool)
		availableAttrs["name"] = true
		availableAttrs["title"] = true

		for _, cred := range config.Credentials {
			for attrName := range cred.Attributes {
				availableAttrs[attrName] = true
			}
		}

		// Check if requested attributes exist
		for _, attr := range config.List.Attributes {
			if !availableAttrs[attr] {
				if configValidateVerbose {
					validationErrors = append(validationErrors, fmt.Sprintf("list attribute '%s' not found in any credential", attr))
				}
			}
		}
	}

	// Report validation results
	fmt.Println()
	if len(validationErrors) == 0 {
		fmt.Println("✓ Configuration Structure: OK")
		fmt.Println()
		fmt.Println("✅ Configuration is valid!")
		fmt.Println()

		// Show summary if verbose
		if configValidateVerbose {
			fmt.Println("Configuration summary:")
			if config.Project != "" {
				fmt.Printf("  Project: %s\n", config.Project)
			} else {
				fmt.Println("  Project: (not set)")
			}
			if config.Prefix != "" {
				fmt.Printf("  Prefix: %s\n", config.Prefix)
			} else {
				fmt.Println("  Prefix: (not set)")
			}
			if len(config.List.Attributes) > 0 {
				fmt.Printf("  List attributes: %v\n", config.List.Attributes)
			} else {
				fmt.Println("  List attributes: (not set)")
			}
			fmt.Printf("  Credentials: %d entries\n", len(config.Credentials))
			if len(config.Defaults.Labels) > 0 {
				fmt.Printf("  Default labels: %d entries\n", len(config.Defaults.Labels))
			}
			fmt.Println()
		}

		return nil
	}

	// Show errors
	fmt.Println("❌ Configuration Structure: FAILED")
	fmt.Println()
	fmt.Printf("Found %d validation error(s):\n", len(validationErrors))
	for i, err := range validationErrors {
		fmt.Printf("  %d. %s\n", i+1, err)
	}
	fmt.Println()

	return fmt.Errorf("configuration validation failed with %d error(s)", len(validationErrors))
}
