package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update SECRET_NAME",
	Short: "Update an existing secret in Google Secret Manager",
	Long: `Update an existing secret by creating a new version with new data.
You can provide the secret value via --data flag, from a file using --data-file,
or interactively (prompt).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretName := args[0]
		project, _ := cmd.Flags().GetString("project")
		data, _ := cmd.Flags().GetString("data")
		dataFile, _ := cmd.Flags().GetString("data-file")

		// Get secret value
		secretValue, err := getSecretInput(data, dataFile, "Enter new secret value: ")
		if err != nil {
			return err
		}

		// Build gcloud command to add new version
		gcloudArgs := []string{"secrets", "versions", "add", secretName}

		if project != "" {
			gcloudArgs = append(gcloudArgs, "--project", project)
		}

		gcloudArgs = append(gcloudArgs, "--data-file", "-")

		// Execute gcloud command
		gcloudCmd := exec.Command("gcloud", gcloudArgs...)
		gcloudCmd.Stdin = strings.NewReader(secretValue)

		output, err := gcloudCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("gcloud command failed: %s", string(output))
		}

		fmt.Printf("Secret '%s' updated successfully\n", secretName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("data", "d", "", "New secret data to store")
	updateCmd.Flags().String("data-file", "", "Path to file containing new secret data")
}
