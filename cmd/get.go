package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get SECRET_NAME",
	Short: "Get a secret value from Google Secret Manager",
	Long: `Retrieve a secret value from Google Secret Manager.

By default, retrieves the latest (most recent) version of the secret.
You can specify a specific version number to access older versions.

Examples:
  gsecutil get my-secret                    # Get latest version
  gsecutil get my-secret --version 3        # Get specific version 3
  gsecutil get my-secret -v 1 --clipboard   # Get version 1 and copy to clipboard
  gsecutil get my-secret --show-metadata    # Show version info along with value`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userInputName := args[0]                           // What the user typed
		secretName := AddPrefixToSecretName(userInputName) // Add prefix if configured
		project, _ := cmd.Flags().GetString("project")
		project = GetProject(project) // Use configuration-based project resolution
		version, _ := cmd.Flags().GetString("version")
		clipboard, _ := cmd.Flags().GetBool("clipboard")
		showMetadata, _ := cmd.Flags().GetBool("show-metadata")

		// Determine version to use
		versionToUse := version
		if versionToUse == "" {
			versionToUse = "latest"
		}

		// Build gcloud command to get secret value
		gcloudArgs := []string{"secrets", "versions", "access", versionToUse, "--secret", secretName}
		if project != "" {
			gcloudArgs = append(gcloudArgs, "--project", project)
		}

		// Execute gcloud command to get secret value
		gcloudCmd := exec.Command("gcloud", gcloudArgs...)
		output, err := gcloudCmd.Output()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				return formatGcloudError(string(exitError.Stderr))
			}
			return fmt.Errorf("failed to execute gcloud command: %v", err)
		}

		secretValue := strings.TrimSpace(string(output))

		// Get metadata if requested
		var versionInfo *SecretVersionInfo
		if showMetadata {
			versionInfo, err = getSecretVersionInfo(secretName, versionToUse, project)
			if err != nil {
				// Don't fail if metadata fetch fails, just warn
				fmt.Printf("Warning: Failed to fetch version metadata: %v\n", err)
			}
		}

		// Display metadata first if requested
		if showMetadata && versionInfo != nil {
			fmt.Printf("Secret: %s\n", secretName)
			fmt.Printf("Version: %s\n", versionInfo.Name)
			fmt.Printf("State: %s\n", versionInfo.State)
			fmt.Printf("Created: %s\n", versionInfo.CreateTime.Format(time.RFC3339))
			if !versionInfo.DestroyTime.IsZero() {
				fmt.Printf("Destroy Time: %s\n", versionInfo.DestroyTime.Format(time.RFC3339))
			}
			fmt.Printf("ETag: %s\n", versionInfo.Etag)
			fmt.Println("---")
		}

		if clipboard {
			// Copy to clipboard
			if err := copyToClipboard(secretValue); err != nil {
				if showMetadata {
					fmt.Printf("Secret Value: %s\n", secretValue)
				} else {
					fmt.Printf("Secret value: %s\n", secretValue)
				}
				fmt.Printf("Warning: Failed to copy to clipboard: %v\n", err)
			} else {
				if showMetadata {
					fmt.Println("Secret value copied to clipboard")
				} else {
					fmt.Println("Secret value copied to clipboard")
				}
			}
		} else {
			if showMetadata {
				fmt.Printf("Secret Value: %s\n", secretValue)
			} else {
				fmt.Println(secretValue)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("version", "v", "", "Version of the secret to retrieve (default: latest)")
	getCmd.Flags().BoolP("clipboard", "c", false, "Copy secret value to clipboard")
	getCmd.Flags().BoolP("show-metadata", "m", false, "Show version metadata (version, created time, state)")
}
