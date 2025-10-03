package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe SECRET_NAME",
	Short: "Describe a secret in Google Secret Manager",
	Long: `Get comprehensive information about a secret including:
- Basic metadata (name, creation time, ETag)
- Default version information (version number, state, creation time)
- Replication strategy (automatic or user-managed)
- Labels (key-value pairs for organization)
- Tags/Annotations (additional metadata)
- Version aliases (if any)
- Expiration and rotation settings (if configured)
- Pub/Sub topics (if configured)

Use --show-versions to also display detailed information about all versions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured
		project, _ := cmd.Flags().GetString("project")
		format, _ := cmd.Flags().GetString("format")
		showVersions, _ := cmd.Flags().GetBool("show-versions")

		// If custom format is specified, use original behavior
		if format != "" {
			gcloudArgs := []string{"secrets", "describe", secretName, "--format", format}
			if project != "" {
				gcloudArgs = append(gcloudArgs, "--project", project)
			}

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

		// Enhanced describe with version information
		// Pass both the full secret name (with prefix) and user input name
		return describeSecretWithVersions(secretName, userInputName, project, showVersions)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().String("format", "", "Output format (e.g., json, yaml)")
	describeCmd.Flags().BoolP("show-versions", "v", false, "Show detailed version information including creation and update times")
}
