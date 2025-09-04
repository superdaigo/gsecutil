package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets in Google Secret Manager",
	Long: `List all secrets in the specified Google Cloud project.
You can filter results and control the output format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		filter, _ := cmd.Flags().GetString("filter")
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")

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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("filter", "", "Filter expression to apply to the list")
	listCmd.Flags().String("format", "", "Output format (e.g., table, json, yaml)")
	listCmd.Flags().Int("limit", 0, "Maximum number of secrets to list (0 for no limit)")
}
