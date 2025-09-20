package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// AuditLogEntry represents an audit log entry from Google Cloud Logging
type AuditLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"`
	LogName   string    `json:"logName"`
	Resource  struct {
		Type   string `json:"type"`
		Labels struct {
			ProjectID string `json:"project_id"`
		} `json:"labels"`
	} `json:"resource"`
	ProtoPayload struct {
		Type               string `json:"@type"`
		AuthenticationInfo struct {
			PrincipalEmail string `json:"principalEmail"`
		} `json:"authenticationInfo"`
		MethodName   string `json:"methodName"`
		ResourceName string `json:"resourceName"`
		Request      struct {
			Name string `json:"name"`
		} `json:"request"`
		Response struct {
			Name string `json:"name"`
		} `json:"response"`
	} `json:"protoPayload"`
}

var auditlogCmd = &cobra.Command{
	Use:   "auditlog [SECRET_NAME]",
	Short: "Show audit log for secret access",
	Long: `Show audit log entries for secrets, including who accessed them,
when they accessed them, and what operations were performed.

When called without a secret name, shows all Secret Manager audit logs.
When called with a secret name, shows logs for secrets matching that name (supports partial matching).

This command queries Google Cloud Logging for audit events related to Secret Manager.
It shows operations like secret access, creation, updates, and deletions.

Note: This command requires Data Access audit logs to be enabled for the Secret Manager API.
See the audit-logging documentation for setup instructions.

Examples:
  gsecutil auditlog                    # Show all Secret Manager audit logs
  gsecutil auditlog my-secret          # Show logs for secrets containing "my-secret"
  gsecutil auditlog --user john        # Show logs for user containing "john"
  gsecutil auditlog --operations ACCESS,CREATE    # Show only ACCESS and CREATE operations
  gsecutil auditlog db --user admin --operations UPDATE    # Specific filters combined`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get command arguments and flags
		var secretName string
		if len(args) > 0 {
			secretName = args[0]
		}

		project, _ := cmd.Flags().GetString("project")
		days, _ := cmd.Flags().GetInt("days")
		limit, _ := cmd.Flags().GetInt("limit")
		format, _ := cmd.Flags().GetString("format")
		userFilter, _ := cmd.Flags().GetString("user")
		operationsFilter, _ := cmd.Flags().GetString("operations")

		return runAuditLogQuery(project, secretName, userFilter, operationsFilter, days, limit, format)
	},
}

// runAuditLogQuery executes the audit log query with filtering
func runAuditLogQuery(project, secretName, userFilter, operationsFilter string, days, limit int, format string) error {
	// Parse operations filter
	operations := parseOperationsFilter(operationsFilter)

	// Build the filter for Secret Manager audit logs
	filter := buildLogFilter(secretName, userFilter, days)

	// Execute gcloud logging command
	logEntries, err := executeLogQuery(project, filter, limit)
	if err != nil {
		return err
	}

	// Filter entries if needed (for partial matching that gcloud filter can't handle well)
	filteredEntries := filterLogEntries(logEntries, secretName, userFilter, operations)

	if len(filteredEntries) == 0 {
		printNoResultsMessage(secretName, userFilter, operationsFilter, days)
		return nil
	}

	// Display results
	return displayLogEntries(filteredEntries, secretName, userFilter, operationsFilter, days, format)
}

// buildLogFilter constructs the gcloud logging filter query
func buildLogFilter(secretName, userFilter string, days int) string {
	// Base filter for Secret Manager service
	filter := `protoPayload.serviceName="secretmanager.googleapis.com"`

	// Add time constraint
	filter += fmt.Sprintf(` AND timestamp>="%s"`, time.Now().AddDate(0, 0, -days).Format(time.RFC3339))

	// Add secret name filter if provided
	if secretName != "" {
		// Use partial matching by searching for the secret name in resource paths
		filter += fmt.Sprintf(` AND (
  protoPayload.resourceName:"%s" 
  OR protoPayload.request.name:"%s"
  OR protoPayload.response.name:"%s"
)`, secretName, secretName, secretName)
	}

	// Add user filter if provided (basic filter, we'll do more precise filtering in post-processing)
	if userFilter != "" {
		filter += fmt.Sprintf(` AND protoPayload.authenticationInfo.principalEmail:"%s"`, userFilter)
	}

	return filter
}

// executeLogQuery runs the gcloud logging read command
func executeLogQuery(project, filter string, limit int) ([]AuditLogEntry, error) {
	// Build gcloud logging command
	gcloudArgs := []string{"logging", "read", filter, "--format", "json"}

	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	if limit > 0 {
		gcloudArgs = append(gcloudArgs, "--limit", fmt.Sprintf("%d", limit))
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	// Parse JSON response
	var logEntries []AuditLogEntry
	if err := json.Unmarshal(output, &logEntries); err != nil {
		return nil, fmt.Errorf("failed to parse log entries: %w", err)
	}

	return logEntries, nil
}

// parseOperationsFilter parses the comma-separated operations filter into a slice
func parseOperationsFilter(operationsFilter string) []string {
	if operationsFilter == "" {
		return nil
	}

	// Split by comma and trim whitespace
	parts := strings.Split(operationsFilter, ",")
	var operations []string
	for _, part := range parts {
		operation := strings.TrimSpace(strings.ToUpper(part))
		if operation != "" {
			operations = append(operations, operation)
		}
	}

	return operations
}

// isValidOperation checks if an operation name is valid
func isValidOperation(operation string) bool {
	validOps := []string{
		"ACCESS", "CREATE", "UPDATE", "DELETE", "GET_METADATA",
		"LIST", "UPDATE_METADATA", "DESTROY_VERSION",
		"DISABLE_VERSION", "ENABLE_VERSION",
	}

	for _, validOp := range validOps {
		if operation == validOp {
			return true
		}
	}
	return false
}

// filterLogEntries performs post-processing filtering for partial matches and secret relevance
func filterLogEntries(entries []AuditLogEntry, secretName, userFilter string, operations []string) []AuditLogEntry {
	var filtered []AuditLogEntry

	for _, entry := range entries {
		// First, check if this is a secret-related operation (not just location listing)
		if !isSecretRelatedOperation(entry) {
			continue
		}

		// Check secret name partial match if specified
		if secretName != "" {
			match := false

			// Check all possible resource name fields
			resourceNames := []string{
				entry.ProtoPayload.ResourceName,
				entry.ProtoPayload.Request.Name,
				entry.ProtoPayload.Response.Name,
			}

			for _, resourceName := range resourceNames {
				if resourceName != "" && strings.Contains(strings.ToLower(resourceName), strings.ToLower(secretName)) {
					match = true
					break
				}
			}

			if !match {
				continue
			}
		}

		// Check user partial match if specified
		if userFilter != "" {
			user := entry.ProtoPayload.AuthenticationInfo.PrincipalEmail
			if user == "" || !strings.Contains(strings.ToLower(user), strings.ToLower(userFilter)) {
				continue
			}
		}

		// Check operations filter if specified
		if len(operations) > 0 {
			operationName := getOperationName(entry.ProtoPayload.MethodName)
			match := false
			for _, allowedOp := range operations {
				if operationName == allowedOp {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		filtered = append(filtered, entry)
	}

	return filtered
}

// isSecretRelatedOperation determines if an audit log entry is actually about secret operations
// rather than general service operations like listing locations
func isSecretRelatedOperation(entry AuditLogEntry) bool {
	methodName := entry.ProtoPayload.MethodName
	resourceName := entry.ProtoPayload.ResourceName
	requestName := entry.ProtoPayload.Request.Name
	responseName := entry.ProtoPayload.Response.Name

	// Check if it's a specific secret operation (not general service operations)
	if strings.Contains(methodName, "AccessSecretVersion") ||
		strings.Contains(methodName, "CreateSecret") ||
		strings.Contains(methodName, "AddSecretVersion") ||
		strings.Contains(methodName, "DeleteSecret") ||
		strings.Contains(methodName, "GetSecret") ||
		strings.Contains(methodName, "UpdateSecret") ||
		strings.Contains(methodName, "DestroySecretVersion") ||
		strings.Contains(methodName, "DisableSecretVersion") ||
		strings.Contains(methodName, "EnableSecretVersion") {
		return true
	}

	// For ListSecrets operations, check if they're targeting specific secrets
	// rather than just listing locations
	if strings.Contains(methodName, "ListSecrets") {
		// Check if any of the resource/request names contain "/secrets/" (actual secret)
		// rather than just "/locations/" (location listing)
		if strings.Contains(resourceName, "/secrets/") ||
			strings.Contains(requestName, "/secrets/") ||
			strings.Contains(responseName, "/secrets/") {
			return true
		}

		// Also include project-level ListSecrets (but exclude location-specific ones)
		if strings.Contains(resourceName, "/locations/") {
			return false // Exclude location-specific listing
		}

		// Include project-level secret listing
		if resourceName != "" && !strings.Contains(resourceName, "/locations/") {
			return true
		}
	}

	return false
}

// printNoResultsMessage displays appropriate message when no results are found
func printNoResultsMessage(secretName, userFilter, operationsFilter string, days int) {
	filters := []string{}
	if secretName != "" {
		filters = append(filters, fmt.Sprintf("secrets matching '%s'", secretName))
	}
	if userFilter != "" {
		filters = append(filters, fmt.Sprintf("user matching '%s'", userFilter))
	}
	if operationsFilter != "" {
		filters = append(filters, fmt.Sprintf("operations '%s'", operationsFilter))
	}

	if len(filters) > 0 {
		fmt.Printf("No audit log entries found for %s in the last %d days.\n", strings.Join(filters, " and "), days)
	} else {
		fmt.Printf("No Secret Manager audit log entries found in the last %d days.\n", days)
	}
	fmt.Println("Note: Audit logs may take some time to appear, and require Cloud Audit Logs to be enabled.")
}

// displayLogEntries formats and displays the log entries
func displayLogEntries(entries []AuditLogEntry, secretName, userFilter, operationsFilter string, days int, format string) error {
	// Display results based on format
	if format == "json" {
		jsonOutput, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON output: %w", err)
		}
		fmt.Println(string(jsonOutput))
		return nil
	}

	// Default table format
	printTableHeader(secretName, userFilter, operationsFilter, days)

	for _, entry := range entries {
		timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
		operation := getOperationName(entry.ProtoPayload.MethodName)
		user := entry.ProtoPayload.AuthenticationInfo.PrincipalEmail
		if user == "" {
			user = "system"
		}

		// Extract resource name
		resourceName := entry.ProtoPayload.ResourceName
		if resourceName == "" {
			resourceName = entry.ProtoPayload.Request.Name
		}
		if resourceName == "" {
			resourceName = entry.ProtoPayload.Response.Name
		}

		// Shorten resource name for display
		if len(resourceName) > 30 {
			parts := strings.Split(resourceName, "/")
			if len(parts) > 2 {
				resourceName = ".../" + strings.Join(parts[len(parts)-2:], "/")
			}
		}

		fmt.Printf("%-20s %-30s %-40s %-30s\n", timestamp, operation, user, resourceName)
	}

	fmt.Printf("\nTotal entries: %d\n", len(entries))
	return nil
}

// printTableHeader prints the appropriate table header based on filters
func printTableHeader(secretName, userFilter, operationsFilter string, days int) {
	filters := []string{}
	if secretName != "" {
		filters = append(filters, fmt.Sprintf("secret '%s'", secretName))
	}
	if userFilter != "" {
		filters = append(filters, fmt.Sprintf("user '%s'", userFilter))
	}
	if operationsFilter != "" {
		filters = append(filters, fmt.Sprintf("operations '%s'", operationsFilter))
	}

	if len(filters) > 0 {
		fmt.Printf("Secret Manager audit logs matching %s (last %d days):\n\n", strings.Join(filters, " and "), days)
	} else {
		fmt.Printf("Secret Manager audit logs (last %d days):\n\n", days)
	}

	fmt.Printf("%-20s %-30s %-40s %-30s\n", "TIMESTAMP", "OPERATION", "USER", "RESOURCE")
	fmt.Println(strings.Repeat("-", 120))
}

// getOperationName converts gcloud method names to human-readable operation names
func getOperationName(methodName string) string {
	switch {
	case strings.Contains(methodName, "AccessSecretVersion"):
		return "ACCESS"
	case strings.Contains(methodName, "CreateSecret"):
		return "CREATE"
	case strings.Contains(methodName, "AddSecretVersion"):
		return "UPDATE"
	case strings.Contains(methodName, "DeleteSecret"):
		return "DELETE"
	case strings.Contains(methodName, "GetSecret"):
		return "GET_METADATA"
	case strings.Contains(methodName, "ListSecrets"):
		return "LIST"
	case strings.Contains(methodName, "UpdateSecret"):
		return "UPDATE_METADATA"
	case strings.Contains(methodName, "DestroySecretVersion"):
		return "DESTROY_VERSION"
	case strings.Contains(methodName, "DisableSecretVersion"):
		return "DISABLE_VERSION"
	case strings.Contains(methodName, "EnableSecretVersion"):
		return "ENABLE_VERSION"
	default:
		// Return the method name without the service prefix
		parts := strings.Split(methodName, ".")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return methodName
	}
}

func init() {
	rootCmd.AddCommand(auditlogCmd)
	auditlogCmd.Flags().IntP("days", "d", 7, "Number of days to look back for audit logs")
	auditlogCmd.Flags().IntP("limit", "l", 50, "Maximum number of log entries to retrieve")
	auditlogCmd.Flags().String("format", "", "Output format (table, json)")
	auditlogCmd.Flags().StringP("user", "u", "", "Filter by username (supports partial matching)")
	auditlogCmd.Flags().StringP("operations", "o", "", "Filter by operations (comma-separated): ACCESS,CREATE,UPDATE,DELETE,GET_METADATA,LIST,UPDATE_METADATA,DESTROY_VERSION,DISABLE_VERSION,ENABLE_VERSION")
}
