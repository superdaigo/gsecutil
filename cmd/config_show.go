package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configShowCmd = &cobra.Command{
	Use:   "show [config-file]",
	Short: "Show gsecutil configuration file contents",
	Long: `Show the contents of a gsecutil configuration file.

This command displays the configuration file in a human-readable format,
showing all settings including project, prefix, list attributes, and
credentials count.

Use --show-credentials to display detailed credentials information in table format.

If no file path is provided, shows the default configuration file.`,
	Example: `  gsecutil config show
  gsecutil config show /path/to/config.yaml
  gsecutil config show --show-credentials     # Show credentials table`,
	Args: cobra.MaximumNArgs(1),
	RunE: runConfigShow,
}

var (
	configShowCredentials bool
)

func init() {
	configCmd.AddCommand(configShowCmd)
	configShowCmd.Flags().BoolVarP(&configShowCredentials, "show-credentials", "c", false, "Show credentials table")
}

func runConfigShow(cmd *cobra.Command, args []string) error {
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

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Display configuration
	fmt.Printf("Configuration file: %s\n", configPath)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Project with source indication
	projectID, source := getProjectWithSource(cmd, &config)
	if projectID != "" {
		fmt.Printf("📦 Project ID: %s (%s)\n", projectID, source)
		fmt.Println()
	} else {
		fmt.Println("📦 Project ID: (not set)")
		fmt.Println()
	}

	// Prefix
	if config.Prefix != "" {
		fmt.Printf("🏷️  Secret Name Prefix:\n")
		fmt.Printf("   %s\n", config.Prefix)
		fmt.Println()
	} else {
		fmt.Println("🏷️  Secret Name Prefix: (not set)")
		fmt.Println()
	}

	// List Configuration
	if len(config.List.Attributes) > 0 {
		fmt.Println("📋 List Display Attributes:")
		for i, attr := range config.List.Attributes {
			fmt.Printf("   %d. %s\n", i+1, attr)
		}
		fmt.Println()
	} else {
		fmt.Println("📋 List Display Attributes: (not set)")
		fmt.Println()
	}

	// Default Labels
	if len(config.Defaults.Labels) > 0 {
		fmt.Println("🏷️  Default Labels:")
		for key, value := range config.Defaults.Labels {
			fmt.Printf("   %s: %s\n", key, value)
		}
		fmt.Println()
	}

	// Credentials
	fmt.Printf("🔐 Credentials: %d entries\n", len(config.Credentials))
	fmt.Println()

	// Show credentials table if requested
	if configShowCredentials && len(config.Credentials) > 0 {
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println()
		displayCredentialsTable(&config)
		fmt.Println()
	}

	fmt.Println("═══════════════════════════════════════════════════════════")

	return nil
}

// getProjectWithSource returns the project ID and its source
func getProjectWithSource(cmd *cobra.Command, config *Config) (string, string) {
	// 1. Check --project flag
	if projectFlag, _ := cmd.Flags().GetString("project"); projectFlag != "" {
		return projectFlag, "from --project flag"
	}

	// 2. Check configuration file
	if config.Project != "" {
		return config.Project, "from config file"
	}

	// 3. Check environment variable
	if envProject := os.Getenv("GSECUTIL_PROJECT"); envProject != "" {
		return envProject, "from GSECUTIL_PROJECT env"
	}

	// 4. Check gcloud default
	cmd2 := exec.Command("gcloud", "config", "get-value", "project")
	output, err := cmd2.Output()
	if err == nil {
		gcloudProject := strings.TrimSpace(string(output))
		if gcloudProject != "" && gcloudProject != "(unset)" {
			return gcloudProject, "from gcloud default"
		}
	}

	return "", ""
}

func displayCredentialsTable(config *Config) {
	if len(config.Credentials) == 0 {
		return
	}

	// Collect all unique attribute keys
	attributeKeys := make(map[string]bool)
	for _, cred := range config.Credentials {
		for key := range cred.Attributes {
			attributeKeys[key] = true
		}
	}

	// Prepare column headers
	headers := []string{"NAME", "TITLE"}
	for key := range attributeKeys {
		headers = append(headers, strings.ToUpper(key))
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	// Calculate widths based on data
	for _, cred := range config.Credentials {
		if len(cred.Name) > colWidths[0] {
			colWidths[0] = len(cred.Name)
		}
		if len(cred.Title) > colWidths[1] {
			colWidths[1] = len(cred.Title)
		}

		for i, key := range headers[2:] {
			keyLower := strings.ToLower(key)
			if val, exists := cred.Attributes[keyLower]; exists {
				valStr := fmt.Sprintf("%v", val)
				if len(valStr) > colWidths[i+2] {
					colWidths[i+2] = len(valStr)
				}
			}
		}
	}

	// Print header
	for i, header := range headers {
		fmt.Printf("%-*s", colWidths[i]+2, header)
	}
	fmt.Println()

	// Print separator
	for i := range headers {
		fmt.Print(strings.Repeat("-", colWidths[i]+2))
	}
	fmt.Println()

	// Print rows
	for _, cred := range config.Credentials {
		// Name
		fmt.Printf("%-*s", colWidths[0]+2, cred.Name)

		// Title
		title := cred.Title
		if title == "" {
			title = "(no title)"
		}
		fmt.Printf("%-*s", colWidths[1]+2, title)

		// Attributes
		for i, key := range headers[2:] {
			keyLower := strings.ToLower(key)
			if val, exists := cred.Attributes[keyLower]; exists {
				valStr := fmt.Sprintf("%v", val)
				fmt.Printf("%-*s", colWidths[i+2]+2, valStr)
			} else {
				fmt.Printf("%-*s", colWidths[i+2]+2, "(unknown)")
			}
		}
		fmt.Println()
	}
}
