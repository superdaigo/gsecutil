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
	Timestamp    time.Time `json:"timestamp"`
	Severity     string    `json:"severity"`
	LogName      string    `json:"logName"`
	Resource     struct {
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

var auditCmd = &cobra.Command{
	Use:   "audit SECRET_NAME",
	Short: "Show audit log for secret access",
	Long: `Show audit log entries for a specific secret, including who accessed it,
when they accessed it, and what operation was performed.

This command queries Google Cloud Logging for audit events related to the specified secret.
It shows operations like secret access, creation, updates, and deletions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretName := args[0]
		project, _ := cmd.Flags().GetString("project")
		days, _ := cmd.Flags().GetInt("days")
		limit, _ := cmd.Flags().GetInt("limit")
		format, _ := cmd.Flags().GetString("format")

		// Build the filter for Secret Manager audit logs
		filter := fmt.Sprintf(`protoPayload.serviceName="secretmanager.googleapis.com"
AND (
  protoPayload.resourceName:"%s" 
  OR protoPayload.request.name:"%s"
  OR protoPayload.response.name:"%s"
)
AND timestamp>="%s"`, 
			secretName, secretName, secretName, 
			time.Now().AddDate(0, 0, -days).Format(time.RFC3339))

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
				return fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
			}
			return fmt.Errorf("failed to execute gcloud command: %v", err)
		}

		// Parse JSON response
		var logEntries []AuditLogEntry
		if err := json.Unmarshal(output, &logEntries); err != nil {
			return fmt.Errorf("failed to parse log entries: %w", err)
		}

		if len(logEntries) == 0 {
			fmt.Printf("No audit log entries found for secret '%s' in the last %d days.\n", secretName, days)
			fmt.Println("Note: Audit logs may take some time to appear, and require Cloud Audit Logs to be enabled.")
			return nil
		}

		// Display results based on format
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(logEntries, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON output: %w", err)
			}
			fmt.Println(string(jsonOutput))
			return nil
		}

		// Default table format
		fmt.Printf("Audit log entries for secret '%s' (last %d days):\n\n", secretName, days)
		fmt.Printf("%-20s %-30s %-40s %-30s\n", "TIMESTAMP", "OPERATION", "USER", "RESOURCE")
		fmt.Println(strings.Repeat("-", 120))

		for _, entry := range logEntries {
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

		fmt.Printf("\nTotal entries: %d\n", len(logEntries))
		return nil
	},
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
	rootCmd.AddCommand(auditCmd)
	auditCmd.Flags().IntP("days", "d", 7, "Number of days to look back for audit logs")
	auditCmd.Flags().IntP("limit", "l", 50, "Maximum number of log entries to retrieve")
	auditCmd.Flags().String("format", "", "Output format (table, json)")
}
