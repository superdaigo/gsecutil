package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var importCmd = &cobra.Command{
	Use:   "import CSV_FILE",
	Short: "Import secrets from CSV file",
	Long: `Import secrets from a CSV file.

The CSV file should have a header row with column names. Required columns:
- name: Secret name
- value: Secret value (required for creation)

Optional columns:
- title: Secret title (stored in config)
- label:*: Labels to apply (e.g., label:env, label:team)
- Any other columns are treated as config attributes

The CSV format supports Excel-style multi-line cells (cells containing newlines
are properly quoted and escaped).

By default, existing secrets are skipped. Use --update to update existing secrets
or --upsert to create new secrets and update existing ones.

The --update-config flag will update the configuration file with titles and
attributes from the CSV.

When a prefix is configured, CSV names are treated as bare names and the prefix
is applied automatically. Names that cannot be resolved to valid prefixed secrets
are skipped.`,
	Example: `  gsecutil import secrets.csv
  gsecutil import secrets.csv --update
  gsecutil import secrets.csv --upsert
  gsecutil import secrets.csv --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().Bool("update", false, "Update existing secrets")
	importCmd.Flags().Bool("upsert", false, "Create or update secrets (upsert)")
	importCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")
	importCmd.Flags().Bool("update-config", false, "Update configuration file with metadata from CSV")
}

func runImport(cmd *cobra.Command, args []string) error {
	project, _ := cmd.Flags().GetString("project")
	project = GetProject(project)
	prefix := GetPrefix()
	importUpdate, _ := cmd.Flags().GetBool("update")
	importUpsert, _ := cmd.Flags().GetBool("upsert")
	importDryRun, _ := cmd.Flags().GetBool("dry-run")
	importUpdateConfig, _ := cmd.Flags().GetBool("update-config")

	csvFile := args[0]

	// Read CSV file
	records, header, err := readCsvFile(csvFile)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Println("No records found in CSV file")
		return nil
	}

	// Validate header
	nameIdx, valueIdx, err := validateHeader(header)
	if err != nil {
		return err
	}

	// Get existing secrets
	existingSecrets, err := getExistingSecretNames(project, prefix)
	if err != nil {
		return fmt.Errorf("failed to get existing secrets: %w", err)
	}

	// Load or create config if update-config is enabled
	var config *Config
	if importUpdateConfig {
		var err error
		config, err = loadOrCreateConfig()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
	}

	// Process records
	stats := &importStats{}
	for i, record := range records {
		if len(record) != len(header) {
			fmt.Printf("Warning: Row %d has %d columns, expected %d. Skipping.\n", i+2, len(record), len(header))
			stats.skipped++
			continue
		}

		userInputName := strings.TrimSpace(record[nameIdx])
		if userInputName == "" {
			fmt.Printf("Warning: Row %d has empty name. Skipping.\n", i+2)
			stats.skipped++
			continue
		}

		resolvedName, bareName, skip, skipReason := resolveImportSecretName(userInputName, prefix)
		if skip {
			fmt.Printf("Warning: Row %d skipped: %s\n", i+2, skipReason)
			stats.skipped++
			continue
		}

		value := ""
		if valueIdx >= 0 {
			value = record[valueIdx]
		}
		exists := existingSecrets[resolvedName]

		// Determine action
		action := ""
		if exists {
			if importUpsert {
				action = "update"
			} else if importUpdate {
				action = "update"
			} else {
				fmt.Printf("Secret '%s' already exists. Skipping. (Use --update or --upsert to update)\n", resolvedName)
				stats.skipped++
				continue
			}
		} else {
			if importUpdate {
				fmt.Printf("Secret '%s' does not exist. Skipping. (Use --upsert to create)\n", resolvedName)
				stats.skipped++
				continue
			} else {
				action = "create"
			}
		}

		// Extract labels and attributes
		labels, title, attributes := extractColumnsData(header, record, nameIdx, valueIdx)

		// Update config if requested (use CSV name as-is)
		if importUpdateConfig && config != nil && !importDryRun {
			updateConfigWithMetadata(config, bareName, title, attributes)
		}

		// Perform action
		if importDryRun {
			fmt.Printf("[DRY-RUN] Would %s secret: %s\n", action, resolvedName)
			stats.processed++
		} else {
			if err := performSecretAction(action, resolvedName, value, labels, project); err != nil {
				fmt.Printf("Error %sing secret '%s': %v\n", action, resolvedName, err)
				stats.failed++
			} else {
				actionDone := map[string]string{"create": "Created", "update": "Updated"}[action]
				fmt.Printf("%s secret: %s\n", actionDone, resolvedName)
				if action == "create" {
					stats.created++
				} else {
					stats.updated++
				}
			}
		}
	}

	// Save config if updated
	if importUpdateConfig && config != nil && !importDryRun {
		if err := saveConfig(config); err != nil {
			fmt.Printf("Warning: Failed to save configuration file: %v\n", err)
		} else {
			fmt.Println("Configuration file updated with metadata from CSV")
		}
	}

	// Print summary
	fmt.Println()
	fmt.Println("Import Summary:")
	if importDryRun {
		fmt.Printf("  Would process: %d\n", stats.processed)
	} else {
		fmt.Printf("  Created: %d\n", stats.created)
		fmt.Printf("  Updated: %d\n", stats.updated)
		fmt.Printf("  Failed: %d\n", stats.failed)
	}
	fmt.Printf("  Skipped: %d\n", stats.skipped)

	return nil
}

type importStats struct {
	created   int
	updated   int
	failed    int
	skipped   int
	processed int
}

func readCsvFile(filename string) ([][]string, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Enable support for multi-line fields (Excel format)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Read all records
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read CSV record: %w", err)
		}
		records = append(records, record)
	}

	return records, header, nil
}

func validateHeader(header []string) (nameIdx, valueIdx int, err error) {
	nameIdx = -1
	valueIdx = -1

	// Check for duplicate column names (case-insensitive)
	columnsSeen := make(map[string][]int)
	for i, col := range header {
		colNormalized := strings.ToLower(strings.TrimSpace(col))
		if colNormalized == "" {
			return -1, -1, fmt.Errorf("CSV header contains empty column name at position %d", i+1)
		}
		columnsSeen[colNormalized] = append(columnsSeen[colNormalized], i+1)
	}

	// Report duplicates
	var duplicates []string
	for col, positions := range columnsSeen {
		if len(positions) > 1 {
			posStr := make([]string, len(positions))
			for i, pos := range positions {
				posStr[i] = fmt.Sprintf("%d", pos)
			}
			duplicates = append(duplicates, fmt.Sprintf("'%s' (columns %s)", col, strings.Join(posStr, ", ")))
		}
	}

	if len(duplicates) > 0 {
		sort.Strings(duplicates)
		return -1, -1, fmt.Errorf("CSV header contains duplicate column names: %s", strings.Join(duplicates, "; "))
	}

	// Find required columns
	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		if col == "name" {
			nameIdx = i
		} else if col == "value" {
			valueIdx = i
		}
	}

	if nameIdx == -1 {
		return -1, -1, fmt.Errorf("CSV must have 'name' column")
	}

	return nameIdx, valueIdx, nil
}

func getExistingSecretNames(project, prefix string) (map[string]bool, error) {
	gcloudArgs := []string{"secrets", "list", "--format", "value(name)"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, formatGcloudError(string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	secrets := make(map[string]bool)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" {
			// Extract just the secret name from full path
			parts := strings.Split(name, "/")
			if len(parts) > 0 {
				name = parts[len(parts)-1]
			}
			if prefix != "" && !strings.HasPrefix(name, prefix) {
				continue
			}
			secrets[name] = true
		}
	}

	return secrets, nil
}

func resolveImportSecretName(userInputName, prefix string) (resolvedName, bareName string, skip bool, skipReason string) {
	if prefix == "" {
		return userInputName, userInputName, false, ""
	}

	if strings.Contains(userInputName, "/") {
		return "", "", true, fmt.Sprintf("name '%s' is invalid; use bare secret names (without resource path)", userInputName)
	}

	if strings.HasPrefix(userInputName, prefix) {
		bareName = strings.TrimPrefix(userInputName, prefix)
		if bareName == "" {
			return "", "", true, fmt.Sprintf("name '%s' is invalid; no secret name remains after removing prefix '%s'", userInputName, prefix)
		}
		return userInputName, bareName, false, ""
	}

	resolvedName = prefix + userInputName
	if !strings.HasPrefix(resolvedName, prefix) {
		return "", "", true, fmt.Sprintf("name '%s' does not match configured prefix '%s'", userInputName, prefix)
	}

	return resolvedName, userInputName, false, ""
}

func extractColumnsData(header, record []string, nameIdx, valueIdx int) (map[string]string, string, map[string]string) {
	labels := make(map[string]string)
	attributes := make(map[string]string)
	title := ""

	for i, col := range header {
		if i == nameIdx || i == valueIdx {
			continue
		}

		colLower := strings.ToLower(strings.TrimSpace(col))
		value := record[i]

		if value == "" {
			continue
		}

		// Check if it's a label column
		if strings.HasPrefix(colLower, "label:") {
			labelKey := strings.TrimPrefix(colLower, "label:")
			labels[labelKey] = value
		} else if colLower == "title" {
			title = value
		} else {
			// Other columns are treated as attributes
			attributes[col] = value
		}
	}

	return labels, title, attributes
}

func performSecretAction(action, name, value string, labels map[string]string, project string) error {
	if action == "create" {
		return createSecretFromImport(name, value, labels, project)
	} else if action == "update" {
		return updateSecretFromImport(name, value, project)
	}
	return fmt.Errorf("unknown action: %s", action)
}

func createSecretFromImport(name, value string, labels map[string]string, project string) error {
	gcloudArgs := []string{"secrets", "create", name}

	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	// Add labels
	for key, val := range labels {
		gcloudArgs = append(gcloudArgs, "--labels", fmt.Sprintf("%s=%s", key, val))
	}

	gcloudArgs = append(gcloudArgs, "--data-file", "-")

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	gcloudCmd.Stdin = strings.NewReader(value)

	output, err := gcloudCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}

func updateSecretFromImport(name, value string, project string) error {
	gcloudArgs := []string{"secrets", "versions", "add", name}

	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	gcloudArgs = append(gcloudArgs, "--data-file", "-")

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	gcloudCmd.Stdin = strings.NewReader(value)

	output, err := gcloudCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}

	return nil
}

func loadOrCreateConfig() (*Config, error) {
	config, err := LoadConfig("")
	if err != nil {
		// If config doesn't exist, create new one
		config = &Config{
			Credentials: []CredentialInfo{},
		}
	}
	return config, nil
}

func updateConfigWithMetadata(config *Config, name, title string, attributes map[string]string) {
	// Find existing credential or create new one
	var credInfo *CredentialInfo
	for i := range config.Credentials {
		if config.Credentials[i].Name == name {
			credInfo = &config.Credentials[i]
			break
		}
	}

	if credInfo == nil {
		// Create new credential entry
		newCred := CredentialInfo{
			Name:       name,
			Attributes: make(map[string]interface{}),
		}
		config.Credentials = append(config.Credentials, newCred)
		credInfo = &config.Credentials[len(config.Credentials)-1]
	}

	// Update title
	if title != "" {
		credInfo.Title = title
	}

	// Update attributes
	if credInfo.Attributes == nil {
		credInfo.Attributes = make(map[string]interface{})
	}
	for key, value := range attributes {
		credInfo.Attributes[strings.ToLower(key)] = value
	}
}

func saveConfig(config *Config) error {
	// Use the path of the loaded config file, or default path if not loaded
	// This ensures we write to the same file that was read
	configPath := configFilePath
	if configPath == "" {
		// No config was loaded yet, use default path
		configPath = getDefaultConfigPath()
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
