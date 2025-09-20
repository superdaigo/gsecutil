package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets in Google Secret Manager",
	Long: `List all secrets in the specified Google Cloud project.
You can filter results and control the output format. By default, includes labels.

Examples:
  gsecutil list                     # List all secrets with labels
  gsecutil list --no-labels         # List secrets without labels
  gsecutil list --format json       # Raw JSON output
  gsecutil list --filter "labels.env=prod"  # Filter by label`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		filter, _ := cmd.Flags().GetString("filter")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		noLabels, _ := cmd.Flags().GetBool("no-labels")

		// If user specified a custom format, use the original gcloud passthrough approach
		if format != "" && format != "table" {
			return runOriginalGcloudList(project, filter, format, limit)
		}

		// Enhanced list with labels
		return listSecretsWithLabels(project, filter, limit, !noLabels)
	},
}

// runOriginalGcloudList runs the original gcloud list command for custom formats
func runOriginalGcloudList(project, filter, format string, limit int) error {
	// Build gcloud command
	gcloudArgs := []string{"secrets", "list"}

	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	if filter != "" {
		gcloudArgs = append(gcloudArgs, "--filter", filter)
	}

	if format != "" {
		gcloudArgs = append(gcloudArgs, "--format", format)
	}

	if limit > 0 {
		gcloudArgs = append(gcloudArgs, "--limit", fmt.Sprintf("%d", limit))
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

// listSecretsWithLabels lists secrets with enhanced formatting including labels
func listSecretsWithLabels(project, filter string, limit int, showLabels bool) error {
	// Build gcloud command to get JSON output
	gcloudArgs := []string{"secrets", "list", "--format", "json"}

	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	if filter != "" {
		gcloudArgs = append(gcloudArgs, "--filter", filter)
	}

	if limit > 0 {
		gcloudArgs = append(gcloudArgs, "--limit", fmt.Sprintf("%d", limit))
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	// Parse JSON response
	var secrets []SecretInfo
	if err := json.Unmarshal(output, &secrets); err != nil {
		return fmt.Errorf("failed to parse secrets list: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found.")
		return nil
	}

	// Sort secrets by name for consistent output
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].Name < secrets[j].Name
	})

	// Display secrets
	if showLabels {
		displaySecretsWithLabels(secrets)
	} else {
		displaySecretsSimple(secrets)
	}

	return nil
}

// displaySecretsWithLabels displays secrets in a table format with labels
func displaySecretsWithLabels(secrets []SecretInfo) {
	// Calculate column widths
	maxNameWidth := 4    // "NAME"
	maxLabelsWidth := 6  // "LABELS"
	maxCreatedWidth := 7 // "CREATED"

	for _, secret := range secrets {
		name := extractSecretName(secret.Name)
		if len(name) > maxNameWidth {
			maxNameWidth = len(name)
		}

		labelsStr := formatLabels(secret.Labels)
		if len(labelsStr) > maxLabelsWidth {
			maxLabelsWidth = len(labelsStr)
		}

		createdStr := secret.CreateTime.Format("2006-01-02")
		if len(createdStr) > maxCreatedWidth {
			maxCreatedWidth = len(createdStr)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s\n", maxNameWidth, "NAME", maxLabelsWidth, "LABELS", maxCreatedWidth, "CREATED")
	fmt.Printf("%s  %s  %s\n", strings.Repeat("-", maxNameWidth), strings.Repeat("-", maxLabelsWidth), strings.Repeat("-", maxCreatedWidth))

	// Print secrets
	for _, secret := range secrets {
		name := extractSecretName(secret.Name)
		labelsStr := formatLabels(secret.Labels)
		createdStr := secret.CreateTime.Format("2006-01-02")
		fmt.Printf("%-*s  %-*s  %-*s\n", maxNameWidth, name, maxLabelsWidth, labelsStr, maxCreatedWidth, createdStr)
	}
}

// displaySecretsSimple displays secrets without labels (similar to original gcloud output)
func displaySecretsSimple(secrets []SecretInfo) {
	maxNameWidth := 4    // "NAME"
	maxCreatedWidth := 7 // "CREATED"

	for _, secret := range secrets {
		name := extractSecretName(secret.Name)
		if len(name) > maxNameWidth {
			maxNameWidth = len(name)
		}

		createdStr := secret.CreateTime.Format("2006-01-02 15:04:05")
		if len(createdStr) > maxCreatedWidth {
			maxCreatedWidth = len(createdStr)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s\n", maxNameWidth, "NAME", maxCreatedWidth, "CREATE_TIME")
	fmt.Printf("%s  %s\n", strings.Repeat("-", maxNameWidth), strings.Repeat("-", maxCreatedWidth))

	// Print secrets
	for _, secret := range secrets {
		name := extractSecretName(secret.Name)
		createdStr := secret.CreateTime.Format("2006-01-02T15:04:05Z")
		fmt.Printf("%-*s  %-*s\n", maxNameWidth, name, maxCreatedWidth, createdStr)
	}
}

// extractSecretName extracts the secret name from the full resource name
func extractSecretName(fullName string) string {
	// Full name format: "projects/PROJECT_ID/secrets/SECRET_NAME"
	parts := strings.Split(fullName, "/")
	if len(parts) >= 4 {
		return parts[3] // Return just the secret name
	}
	return fullName // Fallback to full name if parsing fails
}

// formatLabels formats labels as key=value pairs separated by commas
func formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return "-"
	}

	// Sort labels by key for consistent output
	var keys []string
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build label string
	var labelPairs []string
	for _, key := range keys {
		labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", key, labels[key]))
	}

	return strings.Join(labelPairs, ",")
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("filter", "", "Filter expression to apply to the list")
	listCmd.Flags().String("format", "", "Output format (e.g., table, json, yaml) - custom formats bypass label display")
	listCmd.Flags().Int("limit", 0, "Maximum number of secrets to list (0 for no limit)")
	listCmd.Flags().Bool("no-labels", false, "Hide labels in output")
}
