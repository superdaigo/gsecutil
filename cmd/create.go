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
or interactively (prompt).

Version Management:
The free tier of Google Secret Manager allows up to 6 active secret versions.
Before creating a secret, this command will check if adding a new version would
exceed this limit. If so, it will ask if you want to disable old versions
to stay within the free tier, or proceed anyway (which may incur charges).
Use --force to bypass this check entirely.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		data, _ := cmd.Flags().GetString("data")
		dataFile, _ := cmd.Flags().GetString("data-file")
		labels, _ := cmd.Flags().GetStringSlice("labels")
		force, _ := cmd.Flags().GetBool("force")

		// Get secret value
		secretValue, err := getSecretInput(data, dataFile, "Enter secret value: ")
		if err != nil {
			return err
		}

		// Check if secret already exists for version management
		// For create command, we only need to check if secret exists for update scenarios
		gcloudCheckArgs := []string{"secrets", "describe", secretName, "--format", "value(name)"}
		if project != "" {
			gcloudCheckArgs = append(gcloudCheckArgs, "--project", project)
		}
		gcloudCheckCmd := exec.Command("gcloud", gcloudCheckArgs...)
		if gcloudCheckCmd.Run() == nil {
			// Secret exists, this is actually an update operation
			fmt.Printf("Secret '%s' already exists. This will create a new version.\n", secretName)
			// Perform version management check
			shouldContinue, err := manageVersionsForFreeTier(secretName, project, force)
			if err != nil {
				return err
			}
			if !shouldContinue {
				return fmt.Errorf("operation cancelled")
			}
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
	createCmd.Flags().BoolP("force", "f", false, "Force creation without version limit checks (may exceed free tier)")
}
