package cmd

import (
	"fmt"
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
		data, _ := cmd.Flags().GetString("data")
		dataFile, _ := cmd.Flags().GetString("data-file")
		labels, _ := cmd.Flags().GetStringSlice("labels")

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
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("data", "d", "", "Secret data to store")
	createCmd.Flags().String("data-file", "", "Path to file containing secret data")
	createCmd.Flags().StringSlice("labels", []string{}, "Labels to apply to the secret (format: key=value)")
}
