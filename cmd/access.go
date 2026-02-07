package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// IAMPolicy represents the IAM policy structure from gcloud
type IAMPolicy struct {
	Version  int       `json:"version"`
	Etag     string    `json:"etag"`
	Bindings []Binding `json:"bindings"`
}

// Binding represents an IAM policy binding
type Binding struct {
	Role      string     `json:"role"`
	Members   []string   `json:"members"`
	Condition *Condition `json:"condition,omitempty"`
}

// Condition represents an IAM policy condition
type Condition struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Expression  string `json:"expression,omitempty"`
}

// SecretManagerRoles maps role names to descriptions
var SecretManagerRoles = map[string]string{
	"roles/secretmanager.secretAccessor":       "Secret Accessor (can access secret values)",
	"roles/secretmanager.viewer":               "Secret Viewer (can view metadata only)",
	"roles/secretmanager.admin":                "Secret Manager Admin (full access)",
	"roles/secretmanager.secretVersionManager": "Version Manager (can create/destroy versions)",
	"roles/secretmanager.secretVersionAdder":   "Version Adder (can add new versions)",
}

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage access permissions for secrets",
	Long: `Manage IAM access permissions for Google Secret Manager secrets.

This command allows you to list, grant, and revoke access to secrets for users,
groups, and service accounts. It provides a simplified interface to manage
Secret Manager IAM policies.

Available subcommands:
  list    - List principals with access to a secret
  grant   - Grant access to a principal
  revoke  - Revoke access from a principal
  project - Show project-level Secret Manager permissions`,
}

var accessListCmd = &cobra.Command{
	Use:   "list <secret>",
	Short: "List principals with access to a secret",
	Long: `List all users, groups, and service accounts that have access to the specified secret.

This command shows the IAM policy bindings for the secret, displaying each principal
and their role/permissions. It also shows relevant project-level permissions that
grant access to the secret.

Examples:
  gsecutil access list my-secret                    # List all access for my-secret
  gsecutil access list my-secret --project my-proj  # List access with specific project
  gsecutil access list my-secret --include-project  # Include project-level permissions`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		includeProject, _ := cmd.Flags().GetBool("include-project")
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured
		return listSecretAccess(secretName, project, includeProject)
	},
}

var accessGrantCmd = &cobra.Command{
	Use:   "grant <secret>",
	Short: "Grant access to a principal for a secret",
	Long: `Grant access to a user, group, or service account for the specified secret.

The principal should be specified using the --principal flag and should be in one of these formats:
  - user:email@domain.com
  - group:group@domain.com
  - serviceAccount:sa@project.iam.gserviceaccount.com
  - domain:domain.com

The role defaults to roles/secretmanager.secretAccessor but can be customized with --role.

Examples:
  gsecutil access grant my-secret --principal user:alice@example.com
  gsecutil access grant my-secret --principal user:alice@example.com --role roles/secretmanager.viewer
  gsecutil access grant my-secret --principal serviceAccount:app@project.iam.gserviceaccount.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		principal, _ := cmd.Flags().GetString("principal")
		role, _ := cmd.Flags().GetString("role")
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured

		if principal == "" {
			return fmt.Errorf("--principal is required")
		}

		return grantSecretAccess(secretName, principal, role, project)
	},
}

var accessProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Show project-level Secret Manager permissions",
	Long: `Show all project-level IAM permissions that grant access to Secret Manager.

This command displays all roles at the project level that provide access to secrets,
including specific Secret Manager roles and broader roles like Editor/Owner.

Examples:
  gsecutil access project                    # Show project-level permissions for default project
  gsecutil access project --project my-proj # Show project-level permissions for specific project`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		return showProjectLevelPermissions(project)
	},
}

var accessRevokeCmd = &cobra.Command{
	Use:   "revoke <secret>",
	Short: "Revoke access from a principal for a secret",
	Long: `Revoke access from a user, group, or service account for the specified secret.

The principal should be specified using the --principal flag and should be in one of these formats:
  - user:email@domain.com
  - group:group@domain.com
  - serviceAccount:sa@project.iam.gserviceaccount.com
  - domain:domain.com

You can optionally specify the role to revoke with --role. If no role is specified,
the default role (roles/secretmanager.secretAccessor) will be revoked.

Examples:
  gsecutil access revoke my-secret --principal user:alice@example.com
  gsecutil access revoke my-secret --principal user:alice@example.com --role roles/secretmanager.viewer`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		principal, _ := cmd.Flags().GetString("principal")
		role, _ := cmd.Flags().GetString("role")
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured

		if principal == "" {
			return fmt.Errorf("--principal is required")
		}

		return revokeSecretAccess(secretName, principal, role, project)
	},
}

// listSecretAccess lists all principals with access to a secret
func listSecretAccess(secretName, project string, includeProject bool) error {
	// Build gcloud command to get IAM policy
	gcloudArgs := []string{"secrets", "get-iam-policy", secretName, "--format", "json"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	// Parse JSON response
	var policy IAMPolicy
	if err := json.Unmarshal(output, &policy); err != nil {
		return fmt.Errorf("failed to parse IAM policy: %w", err)
	}

	// Display the access information
	displaySecretAccess(secretName, policy, includeProject, project)

	return nil
}

// grantSecretAccess grants access to a principal for a secret
func grantSecretAccess(secretName, principal, role, project string) error {
	if role == "" {
		role = "roles/secretmanager.secretAccessor"
	}

	// Validate the principal format
	if err := validatePrincipal(principal); err != nil {
		return err
	}

	// Build gcloud command to add IAM policy binding
	gcloudArgs := []string{
		"secrets", "add-iam-policy-binding", secretName,
		"--member", principal,
		"--role", role,
	}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	_, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	fmt.Printf("Successfully granted %s to %s for secret '%s'\n", role, principal, secretName)
	return nil
}

// revokeSecretAccess revokes access from a principal for a secret
func revokeSecretAccess(secretName, principal, role, project string) error {
	if role == "" {
		role = "roles/secretmanager.secretAccessor"
	}

	// Validate the principal format
	if err := validatePrincipal(principal); err != nil {
		return err
	}

	// Build gcloud command to remove IAM policy binding
	gcloudArgs := []string{
		"secrets", "remove-iam-policy-binding", secretName,
		"--member", principal,
		"--role", role,
	}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	_, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	fmt.Printf("Successfully revoked %s from %s for secret '%s'\n", role, principal, secretName)
	return nil
}

// displaySecretAccess formats and displays the access information
func displaySecretAccess(secretName string, policy IAMPolicy, includeProject bool, project string) {
	if len(policy.Bindings) == 0 {
		fmt.Printf("No explicit access permissions found for secret '%s'\n", secretName)
		fmt.Println("Note: Project-level IAM permissions may still provide access")
		return
	}

	fmt.Printf("Access permissions for secret '%s':\n\n", secretName)

	// Sort bindings by role for consistent output
	sort.Slice(policy.Bindings, func(i, j int) bool {
		return policy.Bindings[i].Role < policy.Bindings[j].Role
	})

	for _, binding := range policy.Bindings {
		roleDescription := SecretManagerRoles[binding.Role]
		if roleDescription == "" {
			roleDescription = binding.Role
		}

		fmt.Printf("Role: %s\n", roleDescription)
		fmt.Printf("  Role ID: %s\n", binding.Role)

		if len(binding.Members) > 0 {
			fmt.Println("  Members:")
			// Sort members for consistent output
			sortedMembers := make([]string, len(binding.Members))
			copy(sortedMembers, binding.Members)
			sort.Strings(sortedMembers)

			for _, member := range sortedMembers {
				fmt.Printf("    - %s\n", formatPrincipal(member))
			}
		}

		if binding.Condition != nil {
			fmt.Printf("  Condition: %s\n", binding.Condition.Expression)
			if binding.Condition.Title != "" {
				fmt.Printf("    Title: %s\n", binding.Condition.Title)
			}
			if binding.Condition.Description != "" {
				fmt.Printf("    Description: %s\n", binding.Condition.Description)
			}
		}

		fmt.Println()
	}

	// Display project-level permissions if requested
	if includeProject {
		displayProjectLevelAccess(project)
	}
}

// validatePrincipal validates the format of a principal
func validatePrincipal(principal string) error {
	validPrefixes := []string{"user:", "group:", "serviceAccount:", "domain:", "allUsers", "allAuthenticatedUsers"}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(principal, prefix) {
			return nil
		}
	}

	return fmt.Errorf("invalid principal format: %s\nValid formats: user:email@domain.com, group:group@domain.com, serviceAccount:sa@project.iam.gserviceaccount.com, domain:domain.com, allUsers, allAuthenticatedUsers", principal)
}

// formatPrincipal formats a principal for display
func formatPrincipal(principal string) string {
	parts := strings.SplitN(principal, ":", 2)
	if len(parts) != 2 {
		return principal
	}

	principalType := parts[0]
	principalValue := parts[1]

	switch principalType {
	case "user":
		return fmt.Sprintf("User: %s", principalValue)
	case "group":
		return fmt.Sprintf("Group: %s", principalValue)
	case "serviceAccount":
		return fmt.Sprintf("Service Account: %s", principalValue)
	case "domain":
		return fmt.Sprintf("Domain: %s", principalValue)
	default:
		return principal
	}
}

// getProjectID gets the project ID, using gcloud config if not provided
func getProjectID(project string) string {
	if project != "" {
		return project
	}

	// Try to get from gcloud config
	gcloudCmd := exec.Command("gcloud", "config", "get-value", "project")
	output, err := gcloudCmd.Output()
	if err != nil {
		return "PROJECT_ID" // Fallback
	}

	return strings.TrimSpace(string(output))
}

// displayProjectLevelAccess displays project-level permissions that affect Secret Manager access
func displayProjectLevelAccess(project string) {
	projectID := getProjectID(project)

	fmt.Printf("\n--- Project-Level Permissions (Project: %s) ---\n\n", projectID)

	// Get project IAM policy
	gcloudArgs := []string{"projects", "get-iam-policy", projectID, "--format", "json"}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		fmt.Printf("Warning: Could not retrieve project-level IAM policy: %v\n", err)
		return
	}

	// Parse JSON response
	var policy IAMPolicy
	if err := json.Unmarshal(output, &policy); err != nil {
		fmt.Printf("Warning: Could not parse project-level IAM policy: %v\n", err)
		return
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

	// Filter and display only Secret Manager related roles
	found := false
	for _, binding := range policy.Bindings {
		if secretManagerRoles[binding.Role] && len(binding.Members) > 0 {
			if !found {
				found = true
			}

			roleDescription := SecretManagerRoles[binding.Role]
			if roleDescription == "" {
				if binding.Role == "roles/editor" {
					roleDescription = "Editor (includes Secret Manager access)"
				} else if binding.Role == "roles/owner" {
					roleDescription = "Owner (includes full Secret Manager access)"
				} else {
					roleDescription = binding.Role
				}
			}

			fmt.Printf("Role: %s\n", roleDescription)
			fmt.Printf("  Role ID: %s\n", binding.Role)
			fmt.Printf("  Scope: Project-wide (affects all secrets in project)\n")

			if len(binding.Members) > 0 {
				fmt.Println("  Members:")
				// Sort members for consistent output
				sortedMembers := make([]string, len(binding.Members))
				copy(sortedMembers, binding.Members)
				sort.Strings(sortedMembers)

				for _, member := range sortedMembers {
					fmt.Printf("    - %s\n", formatPrincipal(member))
				}
			}

			if binding.Condition != nil {
				fmt.Printf("  Condition: %s\n", binding.Condition.Expression)
				if binding.Condition.Title != "" {
					fmt.Printf("    Title: %s\n", binding.Condition.Title)
				}
				if binding.Condition.Description != "" {
					fmt.Printf("    Description: %s\n", binding.Condition.Description)
				}
			}

			fmt.Println()
		}
	}

	if !found {
		fmt.Println("No project-level Secret Manager permissions found.")
		fmt.Println("Note: This only shows roles that specifically grant Secret Manager access.")
	}
}

// showProjectLevelPermissions shows project-level permissions without requiring a secret
func showProjectLevelPermissions(project string) error {
	projectID := getProjectID(project)

	fmt.Printf("Project-Level Secret Manager Permissions (Project: %s)\n\n", projectID)

	// Get project IAM policy
	gcloudArgs := []string{"projects", "get-iam-policy", projectID, "--format", "json"}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return formatGcloudError(string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	// Parse JSON response
	var policy IAMPolicy
	if err := json.Unmarshal(output, &policy); err != nil {
		return fmt.Errorf("failed to parse project IAM policy: %w", err)
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

	// Filter and display only Secret Manager related roles
	found := false
	// Sort bindings by role for consistent output
	sort.Slice(policy.Bindings, func(i, j int) bool {
		return policy.Bindings[i].Role < policy.Bindings[j].Role
	})

	for _, binding := range policy.Bindings {
		if secretManagerRoles[binding.Role] && len(binding.Members) > 0 {
			if !found {
				found = true
			}

			roleDescription := SecretManagerRoles[binding.Role]
			if roleDescription == "" {
				if binding.Role == "roles/editor" {
					roleDescription = "Editor (includes Secret Manager access)"
				} else if binding.Role == "roles/owner" {
					roleDescription = "Owner (includes full Secret Manager access)"
				} else {
					roleDescription = binding.Role
				}
			}

			fmt.Printf("Role: %s\n", roleDescription)
			fmt.Printf("  Role ID: %s\n", binding.Role)
			fmt.Printf("  Scope: Project-wide (affects all secrets in project)\n")

			if len(binding.Members) > 0 {
				fmt.Println("  Members:")
				// Sort members for consistent output
				sortedMembers := make([]string, len(binding.Members))
				copy(sortedMembers, binding.Members)
				sort.Strings(sortedMembers)

				for _, member := range sortedMembers {
					fmt.Printf("    - %s\n", formatPrincipal(member))
				}
			}

			if binding.Condition != nil {
				fmt.Printf("  Condition: %s\n", binding.Condition.Expression)
				if binding.Condition.Title != "" {
					fmt.Printf("    Title: %s\n", binding.Condition.Title)
				}
				if binding.Condition.Description != "" {
					fmt.Printf("    Description: %s\n", binding.Condition.Description)
				}
			}

			fmt.Println()
		}
	}

	if !found {
		fmt.Println("No project-level Secret Manager permissions found.")
		fmt.Println("\nThis means:")
		fmt.Println("  - No users/groups have project-wide Secret Manager access")
		fmt.Println("  - Access may be granted at the secret level only")
		fmt.Println("  - Use 'gsecutil access list <secret>' to check secret-specific permissions")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(accessCmd)

	// Add subcommands
	accessCmd.AddCommand(accessListCmd)
	accessCmd.AddCommand(accessGrantCmd)
	accessCmd.AddCommand(accessRevokeCmd)
	accessCmd.AddCommand(accessProjectCmd)

	// Flags for list command
	accessListCmd.Flags().Bool("include-project", false, "Include project-level permissions that grant access to secrets")

	// Flags for grant and revoke commands
	accessGrantCmd.Flags().String("principal", "", "Principal to grant access to (required) - format: user:email@domain.com, group:group@domain.com, etc.")
	accessGrantCmd.Flags().String("role", "roles/secretmanager.secretAccessor", "Role to grant (default: roles/secretmanager.secretAccessor)")
	if err := accessGrantCmd.MarkFlagRequired("principal"); err != nil {
		panic(fmt.Sprintf("Failed to mark principal flag as required for grant command: %v", err))
	}

	accessRevokeCmd.Flags().String("principal", "", "Principal to revoke access from (required) - format: user:email@domain.com, group:group@domain.com, etc.")
	accessRevokeCmd.Flags().String("role", "roles/secretmanager.secretAccessor", "Role to revoke (default: roles/secretmanager.secretAccessor)")
	if err := accessRevokeCmd.MarkFlagRequired("principal"); err != nil {
		panic(fmt.Sprintf("Failed to mark principal flag as required for revoke command: %v", err))
	}
}
