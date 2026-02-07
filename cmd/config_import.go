package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configImportCmd = &cobra.Command{
	Use:   "import <source-file>",
	Short: "Import configuration from an existing file",
	Long: `Import configuration from an existing YAML configuration file.

This command reads a configuration file from the specified path and copies
it to the default gsecutil configuration location (or a custom path if
specified with --output).

The source file will be validated before importing to ensure it's a valid
gsecutil configuration file.`,
	Example: `  gsecutil config import /path/to/team-config.yaml
  gsecutil config import ./gsecutil.conf --output ~/.config/gsecutil/gsecutil.conf
  gsecutil config import remote-config.yaml --force  # Overwrite existing config`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigImport,
}

var (
	configImportOutput string
	configImportForce  bool
)

func init() {
	configCmd.AddCommand(configImportCmd)
	configImportCmd.Flags().StringVarP(&configImportOutput, "output", "o", "", "Output path for configuration file (default: $HOME/.config/gsecutil/gsecutil.conf)")
	configImportCmd.Flags().BoolVarP(&configImportForce, "force", "f", false, "Overwrite existing configuration file")
}

func runConfigImport(cmd *cobra.Command, args []string) error {
	sourcePath := args[0]

	// Check if source file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", sourcePath)
	}

	// Read and validate source file
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Parse to validate YAML structure
	var config Config
	if err := yaml.Unmarshal(sourceData, &config); err != nil {
		return fmt.Errorf("invalid configuration file format: %w", err)
	}

	// Validate configuration contents
	if err := validateConfig(&config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Determine output path
	outputPath := configImportOutput
	if outputPath == "" {
		outputPath = getDefaultConfigPath()
	}

	// Check if output file already exists
	if _, err := os.Stat(outputPath); err == nil && !configImportForce {
		return fmt.Errorf("configuration file already exists at '%s'. Use --force to overwrite", outputPath)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write configuration to output path
	if err := os.WriteFile(outputPath, sourceData, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	fmt.Printf("✓ Configuration imported successfully from: %s\n", sourcePath)
	fmt.Printf("✓ Configuration saved to: %s\n", outputPath)
	fmt.Println()

	// Show summary of imported configuration
	fmt.Println("Imported configuration summary:")
	if config.Project != "" {
		fmt.Printf("  Project: %s\n", config.Project)
	}
	if config.Prefix != "" {
		fmt.Printf("  Prefix: %s\n", config.Prefix)
	}
	if len(config.List.Attributes) > 0 {
		fmt.Printf("  List attributes: %v\n", config.List.Attributes)
	}
	if len(config.Credentials) > 0 {
		fmt.Printf("  Credentials: %d entries\n", len(config.Credentials))
	}
	fmt.Println()

	return nil
}

// validateConfig performs validation on configuration structure
func validateConfig(config *Config) error {
	// Validate prefix format (should not contain spaces)
	if config.Prefix != "" {
		for _, char := range config.Prefix {
			if char == ' ' {
				return fmt.Errorf("prefix cannot contain spaces: '%s'", config.Prefix)
			}
		}
	}

	// Validate credentials
	seenNames := make(map[string]bool)
	for i, cred := range config.Credentials {
		// Check for duplicate names
		if cred.Name == "" {
			return fmt.Errorf("credential at index %d has empty name", i)
		}
		if seenNames[cred.Name] {
			return fmt.Errorf("duplicate credential name found: '%s'", cred.Name)
		}
		seenNames[cred.Name] = true
	}

	return nil
}
