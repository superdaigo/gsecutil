package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// extractSecretName extracts the secret name from the full resource name
// Full name format: "projects/PROJECT_ID/secrets/SECRET_NAME"
func extractSecretName(fullName string) string {
	parts := strings.Split(fullName, "/")
	if len(parts) >= 4 {
		return parts[3] // Return just the secret name
	}
	return fullName // Fallback to full name if parsing fails
}

// fetchSecrets retrieves secrets list from Google Secret Manager
func fetchSecrets(project, filter string, limit int) ([]SecretInfo, error) {
	gcloudArgs := []string{"secrets", "list", "--format", "json"}

	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	if filter != "" {
		gcloudArgs = append(gcloudArgs, "--filter", filter)
	}

	if limit > 0 {
		gcloudArgs = append(gcloudArgs, "--limit", fmt.Sprintf("%d", limit))
	}

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, formatGcloudError(string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	var secrets []SecretInfo
	if err := json.Unmarshal(output, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets list: %w", err)
	}

	return secrets, nil
}

// sortSecrets sorts secrets by name
func sortSecrets(secrets []SecretInfo) {
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].Name < secrets[j].Name
	})
}

// getSecretValue retrieves the latest version value of a secret
func getSecretValue(secretName, project string) string {
	gcloudArgs := []string{"secrets", "versions", "access", "latest", "--secret", secretName}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		return "(error retrieving value)"
	}

	return strings.TrimSpace(string(output))
}
