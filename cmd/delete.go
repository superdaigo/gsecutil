package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete SECRET_NAME",
	Short: "Delete a secret from Google Secret Manager",
	Long: `Delete a secret from Google Secret Manager.
This operation is irreversible and will permanently remove the secret
and all of its versions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretName := args[0]
		project, _ := cmd.Flags().GetString("project")
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf("Are you sure you want to delete secret '%s'? This action is irreversible. (y/N): ", secretName)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read confirmation input: %w", err)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Delete operation cancelled.")
				return nil
			}
		}

		// Build gcloud command
		gcloudArgs := []string{"secrets", "delete", secretName, "--quiet"}

		if project != "" {
			gcloudArgs = append(gcloudArgs, "--project", project)
		}

		// Execute gcloud command
		gcloudCmd := exec.Command("gcloud", gcloudArgs...)
		output, err := gcloudCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("gcloud command failed: %s", string(output))
		}

		fmt.Printf("Secret '%s' deleted successfully\n", secretName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation prompt")
}
