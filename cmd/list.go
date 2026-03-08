package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets in Google Secret Manager",
	Long: `List all secrets in the specified Google Cloud project.
You can filter results and control the output format. Supports configuration-based
attribute display and filtering.

When using --show or list.attributes config, the specified custom attributes
are inserted after the NAME column, followed by built-in fields (LABELS, CREATED).
Built-in Secret Manager fields are always preserved and shown.

Examples:
  gsecutil list                             # List secrets with default attributes from config
  gsecutil list --show-labels               # List secrets with labels
  gsecutil list --format json               # Raw JSON output
  gsecutil list --filter "labels.env=prod"  # Filter by Secret Manager labels
  gsecutil list --attr-filter "environment=prod"  # Filter by config attributes
  gsecutil list --show "title,owner,environment"  # Show: NAME + custom attributes + LABELS + CREATED
  gsecutil list --principal user:alice@example.com  # List secrets accessible by a principal`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		filter, _ := cmd.Flags().GetString("filter")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		showLabels, _ := cmd.Flags().GetBool("show-labels")
		principal, _ := cmd.Flags().GetString("principal")
		attrFilter, _ := cmd.Flags().GetString("attr-filter")
		showAttributes, _ := cmd.Flags().GetString("show")
		// Also check --show-attributes for backward compatibility during transition
		if showAttributes == "" {
			showAttributes, _ = cmd.Flags().GetString("show-attributes")
		}
		showUpdated, _ := cmd.Flags().GetBool("show-updated")

		// Use configuration-based project resolution
		project = GetProject(project)

		// If principal is specified, list secrets accessible by that principal
		if principal != "" {
			return listSecretsForPrincipal(principal, project, showLabels, showUpdated)
		}

		// If user specified a custom format, use the original gcloud passthrough approach
		if format != "" && format != "table" {
			return runOriginalGcloudList(project, filter, format, limit)
		}

		// Handle configuration-based filtering
		if attrFilter != "" {
			return listSecretsWithConfigFiltering(project, filter, limit, attrFilter, showAttributes, showLabels, showUpdated)
		}

		// Enhanced list with potential config attributes
		return listSecretsWithConfigAttributes(project, filter, limit, showAttributes, showLabels, showUpdated)
	},
}

// runOriginalGcloudList runs the original gcloud list command for custom formats
func runOriginalGcloudList(project, filter, format string, limit int) error {
	// Inject prefix filter so custom formats respect the configured prefix
	if prefix := GetPrefix(); prefix != "" {
		prefixFilter := fmt.Sprintf("name:'/secrets/%s'", prefix)
		if filter != "" {
			filter = fmt.Sprintf("(%s) AND %s", filter, prefixFilter)
		} else {
			filter = prefixFilter
		}
	}

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
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	fmt.Print(string(output))
	return nil
}

// listSecretsWithLabels lists secrets with enhanced formatting including labels
func listSecretsWithLabels(project, filter string, limit int, showLabels, showUpdated bool) error {
	secrets, err := fetchSecrets(project, filter, limit)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found.")
		return nil
	}

	// Sort secrets by name for consistent output
	sortSecrets(secrets)

	if showUpdated {
		enrichSecretsWithVersionTimes(secrets, project)
	}

	// Display secrets
	if showLabels {
		displaySecretsWithLabels(secrets, showUpdated)
	} else {
		displaySecretsSimple(secrets, showUpdated)
	}

	return nil
}

// displaySecretsWithLabels displays secrets in a table format with labels
func displaySecretsWithLabels(secrets []SecretInfo, showUpdated bool) {
	prefix := GetPrefix()
	// Calculate column widths using terminal display width (wide chars = 2 cols)
	maxNameWidth := 4    // "NAME"
	maxLabelsWidth := 6  // "LABELS"
	maxCreatedWidth := 7 // "CREATED"
	maxUpdatedWidth := 7 // "UPDATED"

	for _, secret := range secrets {
		name := strings.TrimPrefix(extractSecretName(secret.Name), prefix)
		if w := displayWidth(name); w > maxNameWidth {
			maxNameWidth = w
		}
		labelsStr := formatLabels(secret.Labels)
		if w := displayWidth(labelsStr); w > maxLabelsWidth {
			maxLabelsWidth = w
		}
		createdStr := secret.CreateTime.UTC().Format(datetimeFormat)
		if w := displayWidth(createdStr); w > maxCreatedWidth {
			maxCreatedWidth = w
		}
		if showUpdated {
			if w := displayWidth(formatUpdateTime(secret.LatestVersionTime)); w > maxUpdatedWidth {
				maxUpdatedWidth = w
			}
		}
	}

	// Print header
	header := padRight("NAME", maxNameWidth) + "  " + padRight("LABELS", maxLabelsWidth) + "  " + padRight("CREATED (UTC)", maxCreatedWidth)
	sep := strings.Repeat("-", maxNameWidth) + "  " + strings.Repeat("-", maxLabelsWidth) + "  " + strings.Repeat("-", maxCreatedWidth)
	if showUpdated {
		header += "  " + padRight("UPDATED (UTC)", maxUpdatedWidth)
		sep += "  " + strings.Repeat("-", maxUpdatedWidth)
	}
	fmt.Println(header)
	fmt.Println(sep)

	// Print secrets
	for _, secret := range secrets {
		name := strings.TrimPrefix(extractSecretName(secret.Name), prefix)
		labelsStr := formatLabels(secret.Labels)
		createdStr := secret.CreateTime.UTC().Format(datetimeFormat)
		row := padRight(name, maxNameWidth) + "  " + padRight(labelsStr, maxLabelsWidth) + "  " + padRight(createdStr, maxCreatedWidth)
		if showUpdated {
			row += "  " + padRight(formatUpdateTime(secret.LatestVersionTime), maxUpdatedWidth)
		}
		fmt.Println(row)
	}
}

// displaySecretsSimple displays secrets without labels (similar to original gcloud output)
func displaySecretsSimple(secrets []SecretInfo, showUpdated bool) {
	prefix := GetPrefix()
	maxNameWidth := 4    // "NAME"
	maxCreatedWidth := 7 // "CREATED"
	maxUpdatedWidth := 7 // "UPDATED"

	for _, secret := range secrets {
		name := strings.TrimPrefix(extractSecretName(secret.Name), prefix)
		if w := displayWidth(name); w > maxNameWidth {
			maxNameWidth = w
		}
		createdStr := secret.CreateTime.UTC().Format(datetimeFormat)
		if w := displayWidth(createdStr); w > maxCreatedWidth {
			maxCreatedWidth = w
		}
		if showUpdated {
			if w := displayWidth(formatUpdateTime(secret.LatestVersionTime)); w > maxUpdatedWidth {
				maxUpdatedWidth = w
			}
		}
	}

	// Print header
	header := padRight("NAME", maxNameWidth) + "  " + padRight("CREATED (UTC)", maxCreatedWidth)
	sep := strings.Repeat("-", maxNameWidth) + "  " + strings.Repeat("-", maxCreatedWidth)
	if showUpdated {
		header += "  " + padRight("UPDATED (UTC)", maxUpdatedWidth)
		sep += "  " + strings.Repeat("-", maxUpdatedWidth)
	}
	fmt.Println(header)
	fmt.Println(sep)

	// Print secrets
	for _, secret := range secrets {
		name := strings.TrimPrefix(extractSecretName(secret.Name), prefix)
		createdStr := secret.CreateTime.UTC().Format(datetimeFormat)
		row := padRight(name, maxNameWidth) + "  " + padRight(createdStr, maxCreatedWidth)
		if showUpdated {
			row += "  " + padRight(formatUpdateTime(secret.LatestVersionTime), maxUpdatedWidth)
		}
		fmt.Println(row)
	}
}

// enrichSecretsWithVersionTimes fetches the latest version createTime for each secret
// concurrently and stores it in LatestVersionTime. Secrets with no versions show "-".
func enrichSecretsWithVersionTimes(secrets []SecretInfo, project string) {
	const maxConcurrency = 10
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	for i := range secrets {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			name := extractSecretName(secrets[idx].Name)
			if versionInfo, err := getSecretVersionInfo(name, "latest", project); err == nil {
				secrets[idx].LatestVersionTime = versionInfo.CreateTime
			}
		}(i)
	}
	wg.Wait()
}

// datetimeFormat is the standard format for displaying timestamps in list output
const datetimeFormat = "2006-01-02 15:04"

// formatUpdateTime formats an update time for display, returning "-" if not set
func formatUpdateTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.UTC().Format(datetimeFormat)
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

// listSecretsForPrincipal lists all secrets that a principal has access to
func listSecretsForPrincipal(principal, project string, showLabels, showUpdated bool) error {
	// Validate the principal format
	if err := validatePrincipalFormat(principal); err != nil {
		return err
	}

	// Get all secrets in the project
	allSecrets, err := fetchSecrets(project, "", 0)
	if err != nil {
		return err
	}

	if len(allSecrets) == 0 {
		fmt.Printf("No secrets found in the project.\n")
		return nil
	}

	// Filter secrets that the principal has access to
	var accessibleSecrets []SecretInfo
	for _, secret := range allSecrets {
		secretName := extractSecretName(secret.Name)
		hasAccess, err := checkPrincipalAccess(secretName, principal, project)
		if err != nil {
			// Log warning but continue with other secrets
			fmt.Printf("Warning: Could not check access for secret '%s': %v\n", secretName, err)
			continue
		}
		if hasAccess {
			accessibleSecrets = append(accessibleSecrets, secret)
		}
	}

	if len(accessibleSecrets) == 0 {
		fmt.Printf("No secrets found that '%s' has access to.\n", principal)
		fmt.Println("Note: This checks both secret-level and project-level IAM permissions.")
		return nil
	}

	// Sort secrets by name for consistent output
	sortSecrets(accessibleSecrets)

	if showUpdated {
		enrichSecretsWithVersionTimes(accessibleSecrets, project)
	}

	fmt.Printf("Secrets accessible by '%s':\n\n", principal)

	// Display accessible secrets
	if showLabels {
		displaySecretsWithLabels(accessibleSecrets, showUpdated)
	} else {
		displaySecretsSimple(accessibleSecrets, showUpdated)
	}

	return nil
}

// checkPrincipalAccess checks if a principal has access to a specific secret
func checkPrincipalAccess(secretName, principal, project string) (bool, error) {
	// First check secret-level permissions
	hasSecretAccess, err := checkSecretLevelAccess(secretName, principal, project)
	if err != nil {
		return false, err
	}
	if hasSecretAccess {
		return true, nil
	}

	// Then check project-level permissions
	hasProjectAccess, err := checkProjectLevelAccess(principal, project)
	if err != nil {
		return false, err
	}

	return hasProjectAccess, nil
}

// checkSecretLevelAccess checks if a principal has secret-level access
func checkSecretLevelAccess(secretName, principal, project string) (bool, error) {
	// Get IAM policy for the secret
	gcloudArgs := []string{"secrets", "get-iam-policy", secretName, "--format", "json"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		// If we can't get the policy, assume no access
		return false, nil
	}

	// Parse JSON response
	var policy IAMPolicy
	if err := json.Unmarshal(output, &policy); err != nil {
		return false, err
	}

	// Check if the principal is in any of the bindings
	for _, binding := range policy.Bindings {
		for _, member := range binding.Members {
			if member == principal {
				return true, nil
			}
		}
	}

	return false, nil
}

// checkProjectLevelAccess checks if a principal has project-level Secret Manager access
func checkProjectLevelAccess(principal, project string) (bool, error) {
	projectID := getProjectID(project)

	// Get project IAM policy
	gcloudArgs := []string{"projects", "get-iam-policy", projectID, "--format", "json"}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		// If we can't get the project policy, assume no access
		return false, nil
	}

	// Parse JSON response
	var policy IAMPolicy
	if err := json.Unmarshal(output, &policy); err != nil {
		return false, err
	}

	// Define roles that provide Secret Manager access
	secretManagerRoles := map[string]bool{
		"roles/secretmanager.admin":                true,
		"roles/secretmanager.secretAccessor":       true,
		"roles/secretmanager.viewer":               true,
		"roles/secretmanager.secretVersionManager": true,
		"roles/secretmanager.secretVersionAdder":   true,
		"roles/editor":                             true,
		"roles/owner":                              true,
	}

	// Check if the principal has any Secret Manager roles at project level
	for _, binding := range policy.Bindings {
		if secretManagerRoles[binding.Role] {
			for _, member := range binding.Members {
				if member == principal {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// validatePrincipalFormat validates the format of a principal (used in list command)
func validatePrincipalFormat(principal string) error {
	validPrefixes := []string{"user:", "group:", "serviceAccount:", "domain:", "allUsers", "allAuthenticatedUsers"}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(principal, prefix) {
			return nil
		}
	}

	return fmt.Errorf("invalid principal format: %s\nValid formats: user:email@domain.com, group:group@domain.com, serviceAccount:sa@project.iam.gserviceaccount.com, domain:domain.com, allUsers, allAuthenticatedUsers", principal)
}

// listSecretsWithConfigAttributes lists secrets with configuration-based attribute display
func listSecretsWithConfigAttributes(project, filter string, limit int, showAttributes string, showLabels, showUpdated bool) error {
	// Get secrets first
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
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	// Parse JSON response
	var secrets []SecretInfo
	if err := json.Unmarshal(output, &secrets); err != nil {
		return fmt.Errorf("failed to parse secrets list: %w", err)
	}

	// Filter by prefix if configured
	if prefix := GetPrefix(); prefix != "" {
		var filteredSecrets []SecretInfo
		for _, secret := range secrets {
			secretName := extractSecretName(secret.Name)
			if strings.HasPrefix(secretName, prefix) {
				filteredSecrets = append(filteredSecrets, secret)
			}
		}
		secrets = filteredSecrets
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found.")
		return nil
	}

	// Sort secrets by name for consistent output
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].Name < secrets[j].Name
	})

	if showUpdated {
		enrichSecretsWithVersionTimes(secrets, project)
	}

	// Determine which attributes to show
	var attributes []string
	if showAttributes != "" {
		// CLI parameter overrides everything
		attributes = ParseShowAttributes(showAttributes)
	} else {
		// Use config file settings
		attributes = GetListAttributes()
	}

	// Display secrets with or without config attributes
	if len(attributes) > 0 {
		displaySecretsWithConfigAttributes(secrets, attributes, showLabels, showUpdated)
	} else if showLabels {
		displaySecretsWithLabels(secrets, showUpdated)
	} else {
		displaySecretsSimple(secrets, showUpdated)
	}

	return nil
}

// listSecretsWithConfigFiltering
func listSecretsWithConfigFiltering(project, filter string, limit int, filterAttributes, showAttributes string, showLabels, showUpdated bool) error {
	// Parse filter attributes
	filters, err := ParseFilterAttributes(filterAttributes)
	if err != nil {
		return fmt.Errorf("invalid attr-filter: %w", err)
	}

	// Filter credentials based on attributes
	filteredCredentials := FilterCredentialsByAttributes(filters)
	if len(filteredCredentials) == 0 {
		fmt.Println("No secrets match the specified attribute filters.")
		return nil
	}

	// Get all secrets to match against filtered credentials
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
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	// Parse JSON response
	var allSecrets []SecretInfo
	if err := json.Unmarshal(output, &allSecrets); err != nil {
		return fmt.Errorf("failed to parse secrets list: %w", err)
	}

	// Match secrets with filtered credentials
	var matchingSecrets []SecretInfo
	for _, secret := range allSecrets {
		secretName := extractSecretName(secret.Name)
		bareName := strings.TrimPrefix(secretName, GetPrefix()) // config stores bare names
		for _, cred := range filteredCredentials {
			if cred.Name == bareName {
				matchingSecrets = append(matchingSecrets, secret)
				break
			}
		}
	}

	if len(matchingSecrets) == 0 {
		fmt.Println("No secrets found matching the attribute filters.")
		return nil
	}

	// Sort secrets by name for consistent output
	sort.Slice(matchingSecrets, func(i, j int) bool {
		return matchingSecrets[i].Name < matchingSecrets[j].Name
	})

	if showUpdated {
		enrichSecretsWithVersionTimes(matchingSecrets, project)
	}

	// Determine which attributes to show
	var attributes []string
	if showAttributes != "" {
		// CLI parameter overrides everything
		attributes = ParseShowAttributes(showAttributes)
	} else {
		// Use config file settings
		attributes = GetListAttributes()
	}

	// Display filtered secrets with config attributes
	if len(attributes) > 0 {
		displaySecretsWithConfigAttributes(matchingSecrets, attributes, showLabels, showUpdated)
	} else if showLabels {
		displaySecretsWithLabels(matchingSecrets, showUpdated)
	} else {
		displaySecretsSimple(matchingSecrets, showUpdated)
	}

	return nil
}

// displaySecretsWithConfigAttributes displays secrets with configuration-based attributes
// Custom attributes are inserted after NAME, LABELS is shown only if showLabels is true
func displaySecretsWithConfigAttributes(secrets []SecretInfo, attributes []string, showLabels, showUpdated bool) {
	// Calculate column widths for built-in fields
	maxNameWidth := 4    // "NAME"
	maxLabelsWidth := 6  // "LABELS"
	maxCreatedWidth := 7 // "CREATED"
	maxUpdatedWidth := 7 // "UPDATED"
	attributeWidths := make([]int, len(attributes))

	// Initialize attribute widths with header names (display width)
	for i, attr := range attributes {
		attributeWidths[i] = displayWidth(strings.ToUpper(attr))
	}

	// Calculate widths based on content using terminal display width
	prefix := GetPrefix()
	for _, secret := range secrets {
		secretName := strings.TrimPrefix(extractSecretName(secret.Name), prefix)
		if w := displayWidth(secretName); w > maxNameWidth {
			maxNameWidth = w
		}

		// Calculate labels width
		if showLabels {
			labelsStr := formatLabels(secret.Labels)
			if w := displayWidth(labelsStr); w > maxLabelsWidth {
				maxLabelsWidth = w
			}
		}

		// Calculate created width
		createdStr := secret.CreateTime.UTC().Format(datetimeFormat)
		if w := displayWidth(createdStr); w > maxCreatedWidth {
			maxCreatedWidth = w
		}

		if showUpdated {
			if w := displayWidth(formatUpdateTime(secret.LatestVersionTime)); w > maxUpdatedWidth {
				maxUpdatedWidth = w
			}
		}

		cred := GetCredentialInfo(secretName) // secretName is already bare after TrimPrefix above
		for i, attr := range attributes {
			value := GetAttributeValue(cred, attr)
			if w := displayWidth(value); w > attributeWidths[i] {
				attributeWidths[i] = w
			}
		}
	}

	// Print header: NAME + custom attributes + built-in fields
	header := padRight("NAME", maxNameWidth)
	for i, attr := range attributes {
		header += "  " + padRight(strings.ToUpper(attr), attributeWidths[i])
	}
	if showLabels {
		header += "  " + padRight("LABELS", maxLabelsWidth)
	}
	header += "  " + padRight("CREATED (UTC)", maxCreatedWidth)
	if showUpdated {
		header += "  " + padRight("UPDATED (UTC)", maxUpdatedWidth)
	}
	fmt.Println(header)

	// Print separator
	separator := strings.Repeat("-", maxNameWidth)
	for _, width := range attributeWidths {
		separator += "  " + strings.Repeat("-", width)
	}
	if showLabels {
		separator += "  " + strings.Repeat("-", maxLabelsWidth)
	}
	separator += "  " + strings.Repeat("-", maxCreatedWidth)
	if showUpdated {
		separator += "  " + strings.Repeat("-", maxUpdatedWidth)
	}
	fmt.Println(separator)

	// Print secrets: NAME + custom attributes + built-in fields
	for _, secret := range secrets {
		secretName := strings.TrimPrefix(extractSecretName(secret.Name), prefix)

		cred := GetCredentialInfo(secretName) // bare name

		// Start with NAME
		row := padRight(secretName, maxNameWidth)

		// Add custom attributes
		for i, attr := range attributes {
			value := GetAttributeValue(cred, attr)
			row += "  " + padRight(value, attributeWidths[i])
		}

		// Add labels if requested
		createdStr := secret.CreateTime.UTC().Format(datetimeFormat)
		if showLabels {
			labelsStr := formatLabels(secret.Labels)
			row += "  " + padRight(labelsStr, maxLabelsWidth)
		}
		row += "  " + padRight(createdStr, maxCreatedWidth)
		if showUpdated {
			row += "  " + padRight(formatUpdateTime(secret.LatestVersionTime), maxUpdatedWidth)
		}

		fmt.Println(row)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("filter", "", "Filter expression to apply to Secret Manager labels")
	listCmd.Flags().String("attr-filter", "", "Filter by configuration file attributes (format: key=value,key2=value2)")
	listCmd.Flags().String("show", "", "Comma-separated list of attributes to display from configuration file (inserted after NAME, before built-in fields)")
	listCmd.Flags().String("show-attributes", "", "(Alias for --show) Comma-separated list of attributes to display from configuration file")
	listCmd.Flags().MarkHidden("show-attributes") // Hide from help but keep for compatibility
	listCmd.Flags().String("format", "", "Output format (e.g., table, json, yaml) - custom formats bypass attribute display")
	listCmd.Flags().Int("limit", 0, "Maximum number of secrets to list (0 for no limit)")
	listCmd.Flags().Bool("show-labels", false, "Show labels in output")
	listCmd.Flags().String("principal", "", "List secrets accessible by this principal (format: user:email@domain.com, group:group@domain.com, etc.)")
	listCmd.Flags().Bool("show-updated", false, "Show UPDATED column (fetches latest version time per secret; slower for large lists)")
}
