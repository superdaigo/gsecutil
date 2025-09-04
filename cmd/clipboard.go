package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"golang.org/x/term"
)

// SecretVersionInfo represents version metadata from Google Secret Manager
type SecretVersionInfo struct {
	Name        string    `json:"name"`
	CreateTime  time.Time `json:"createTime"`
	DestroyTime time.Time `json:"destroyTime"`
	State       string    `json:"state"`
	Etag        string    `json:"etag"`
}

// copyToClipboard copies the given text to the system clipboard
func copyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// getSecretInput handles getting secret value from various sources
func getSecretInput(data, dataFile, prompt string) (string, error) {
	if data != "" {
		return data, nil
	}

	if dataFile != "" {
		content, err := os.ReadFile(dataFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", dataFile, err)
		}
		return string(content), nil
	}

	// Interactive prompt
	fmt.Print(prompt)
	byteValue, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read secret value: %w", err)
	}
	fmt.Println() // Add newline after password input
	return string(byteValue), nil
}

// getSecretVersionInfo retrieves version metadata for a secret
func getSecretVersionInfo(secretName, version, project string) (*SecretVersionInfo, error) {
	// Build gcloud command to get version metadata
	gcloudArgs := []string{"secrets", "versions", "describe", version, "--secret", secretName, "--format", "json"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	// Execute gcloud command
	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	// Parse JSON response
	var versionInfo SecretVersionInfo
	if err := json.Unmarshal(output, &versionInfo); err != nil {
		return nil, fmt.Errorf("failed to parse version metadata: %w", err)
	}

	return &versionInfo, nil
}

// SecretInfo represents basic secret metadata
type SecretInfo struct {
	Name       string            `json:"name"`
	CreateTime time.Time         `json:"createTime"`
	Labels     map[string]string `json:"labels"`
	Etag       string            `json:"etag"`
}

// describeSecretWithVersions provides enhanced secret description with version information
func describeSecretWithVersions(secretName, project string, showVersions bool) error {
	// Get basic secret information
	gcloudArgs := []string{"secrets", "describe", secretName, "--format", "json"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	var secretInfo SecretInfo
	if err := json.Unmarshal(output, &secretInfo); err != nil {
		return fmt.Errorf("failed to parse secret metadata: %w", err)
	}

	// Display basic secret information
	fmt.Printf("Name: %s\n", secretInfo.Name)
	fmt.Printf("Created: %s\n", secretInfo.CreateTime.Format(time.RFC3339))
	fmt.Printf("ETag: %s\n", secretInfo.Etag)

	if len(secretInfo.Labels) > 0 {
		fmt.Println("Labels:")
		for key, value := range secretInfo.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	if showVersions {
		fmt.Println("\n--- Versions ---")
		return listSecretVersions(secretName, project)
	}

	return nil
}

// listSecretVersions lists all versions of a secret with their metadata
func listSecretVersions(secretName, project string) error {
	// Get version list
	gcloudArgs := []string{"secrets", "versions", "list", secretName, "--format", "json"}
	if project != "" {
		gcloudArgs = append(gcloudArgs, "--project", project)
	}

	gcloudCmd := exec.Command("gcloud", gcloudArgs...)
	output, err := gcloudCmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("gcloud command failed: %s", string(exitError.Stderr))
		}
		return fmt.Errorf("failed to execute gcloud command: %w", err)
	}

	var versions []SecretVersionInfo
	if err := json.Unmarshal(output, &versions); err != nil {
		return fmt.Errorf("failed to parse version list: %w", err)
	}

	if len(versions) == 0 {
		fmt.Println("No versions found.")
		return nil
	}

	// Sort versions by creation time (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreateTime.After(versions[j].CreateTime)
	})

	// Display versions
	for i, version := range versions {
		if i > 0 {
			fmt.Println()
		}

		// Extract version number from name (e.g., "projects/.../versions/1" -> "1")
		versionNumber := version.Name
		if parts := strings.Split(version.Name, "/"); len(parts) > 0 {
			versionNumber = parts[len(parts)-1]
		}

		fmt.Printf("Version: %s\n", versionNumber)
		fmt.Printf("  State: %s\n", version.State)
		fmt.Printf("  Created: %s\n", version.CreateTime.Format(time.RFC3339))
		if !version.DestroyTime.IsZero() {
			fmt.Printf("  Destroy Time: %s\n", version.DestroyTime.Format(time.RFC3339))
		}
		fmt.Printf("  ETag: %s\n", version.Etag)
	}

	return nil
}
