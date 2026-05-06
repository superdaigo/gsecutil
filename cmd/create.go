package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create SECRET_NAME",
	Short: "Create a new secret in Google Secret Manager",
	Long: `Create a new secret in Google Secret Manager.
You can provide the secret value via --data flag, from a file using --data-file,
or interactively (prompt).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		data, _ := cmd.Flags().GetString("data")
		dataFile, _ := cmd.Flags().GetString("data-file")
		labels, _ := cmd.Flags().GetStringSlice("labels")
		title, _ := cmd.Flags().GetString("title")

		// Merge default labels from config with user-provided labels
		labels = mergeLabelsWithDefaults(labels)

		// Create command should fail for existing secrets.
		exists, err := secretExists(secretName, project)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("secret '%s' already exists. Use `gsecutil update %s` to create a new version", secretName, userInputName)
		}

		// Get secret value
		secretValue, err := getSecretInput(data, dataFile, "Enter secret value: ")
		if err != nil {
			return err
		}

		// Build gcloud command to create secret
		gcloudArgs := []string{"secrets", "create", secretName}

		if project != "" {
			gcloudArgs = append(gcloudArgs, "--project", project)
		}

		// Add labels if provided
		for _, label := range labels {
			gcloudArgs = append(gcloudArgs, "--labels", label)
		}

		gcloudArgs = append(gcloudArgs, "--data-file", "-")

		// Execute gcloud command
		gcloudCmd := exec.Command("gcloud", gcloudArgs...)
		gcloudCmd.Stdin = strings.NewReader(secretValue)

		output, err := gcloudCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("gcloud command failed: %s", string(output))
		}

		fmt.Printf("Secret '%s' created successfully\n", secretName)

		// Save title to config file if provided
		if title != "" {
			if err := saveTitleToConfig(userInputName, title); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to save title to config: %v\n", err)
			} else {
				fmt.Printf("Title saved to configuration file\n")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("data", "d", "", "Secret data to store")
	createCmd.Flags().String("data-file", "", "Path to file containing secret data")
	createCmd.Flags().StringSlice("labels", []string{}, "Labels to apply to the secret (format: key=value)")
	createCmd.Flags().StringP("title", "t", "", "Title for the secret (saved to config file)")
}

func secretExists(secretName, project string) (bool, error) {
	gcloudArgs := []string{"secrets", "describe", secretName, "--format", "value(name)"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.CombinedOutput()
	if err == nil {
		return true, nil
	}

	combinedOutput := strings.ToLower(string(output))
	if strings.Contains(combinedOutput, "not found") || strings.Contains(combinedOutput, "notfound") {
		return false, nil
	}

	return false, formatGcloudError(string(output))
}

// saveTitleToConfig saves the secret title to the configuration file
func saveTitleToConfig(secretName, title string) error {
	config, err := loadOrCreateConfig()
	if err != nil {
		return err
	}

	updateConfigWithMetadata(config, secretName, title, nil)

	return saveConfig(config)
}

// mergeLabelsWithDefaults merges default labels from config with user-provided labels.
// User-provided labels take precedence over default labels.
func mergeLabelsWithDefaults(userLabels []string) []string {
	config := GetConfig()

	// If no default labels configured, return user labels as-is
	if len(config.Defaults.Labels) == 0 {
		return userLabels
	}

	// Parse user-provided labels into a map
	userLabelMap := make(map[string]string)
	for _, label := range userLabels {
		parts := strings.SplitN(label, "=", 2)
		if len(parts) == 2 {
			userLabelMap[parts[0]] = parts[1]
		}
	}

	// Start with default labels
	mergedMap := make(map[string]string)
	for key, value := range config.Defaults.Labels {
		mergedMap[key] = value
	}

	// Override with user-provided labels
	for key, value := range userLabelMap {
		mergedMap[key] = value
	}

	// Convert back to string slice format
	result := make([]string, 0, len(mergedMap))
	for key, value := range mergedMap {
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}

	return result
}
