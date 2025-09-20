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

// SecretInfo represents comprehensive secret metadata
type SecretInfo struct {
	Name        string            `json:"name"`
	CreateTime  time.Time         `json:"createTime"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Etag        string            `json:"etag"`
	Replication struct {
		Automatic   interface{} `json:"automatic,omitempty"`
		UserManaged interface{} `json:"userManaged,omitempty"`
	} `json:"replication"`
	VersionAliases map[string]string `json:"versionAliases"`
	ExpireTime     *time.Time        `json:"expireTime,omitempty"`
	Ttl            string            `json:"ttl,omitempty"`
	Rotation       struct {
		NextRotationTime *time.Time `json:"nextRotationTime,omitempty"`
		RotationPeriod   string     `json:"rotationPeriod,omitempty"`
	} `json:"rotation,omitempty"`
	Topics []struct {
		Name string `json:"name"`
	} `json:"topics,omitempty"`
}

// describeSecretWithVersions provides enhanced secret description with comprehensive information
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
		return fmt.Errorf("failed to execute gcloud command: %v", err)
	}

	var secretInfo SecretInfo
	if err := json.Unmarshal(output, &secretInfo); err != nil {
		return fmt.Errorf("failed to parse secret metadata: %w", err)
	}

	// Get default version information
	defaultVersion, err := getDefaultVersionInfo(secretName, project)
	if err != nil {
		// Don't fail the whole command if we can't get version info
		fmt.Printf("Warning: Could not retrieve default version info: %v\n", err)
	}

	return displayEnhancedSecretInfo(secretInfo, defaultVersion, secretName, project, showVersions)
}

// getDefaultVersionInfo retrieves information about the default (latest enabled) version
func getDefaultVersionInfo(secretName, project string) (*SecretVersionInfo, error) {
	return getSecretVersionInfo(secretName, "latest", project)
}

// displayEnhancedSecretInfo displays comprehensive secret information
func displayEnhancedSecretInfo(secretInfo SecretInfo, defaultVersion *SecretVersionInfo, secretName, project string, showVersions bool) error {
	// Basic information
	fmt.Printf("Name: %s\n", secretInfo.Name)
	fmt.Printf("Created: %s\n", secretInfo.CreateTime.Format(time.RFC3339))
	fmt.Printf("ETag: %s\n", secretInfo.Etag)

	// Default version information
	if defaultVersion != nil {
		versionNumber := extractVersionNumber(defaultVersion.Name)
		fmt.Printf("Default Version: %s\n", versionNumber)
		fmt.Printf("Default Version State: %s\n", defaultVersion.State)
		fmt.Printf("Default Version Created: %s\n", defaultVersion.CreateTime.Format(time.RFC3339))
		if !defaultVersion.DestroyTime.IsZero() {
			fmt.Printf("Default Version Destroy Time: %s\n", defaultVersion.DestroyTime.Format(time.RFC3339))
		}
	}

	// Replication strategy
	fmt.Printf("Replication: %s\n", getReplicationStrategy(secretInfo.Replication))

	// Labels
	if len(secretInfo.Labels) > 0 {
		fmt.Println("Labels:")
		keys := make([]string, 0, len(secretInfo.Labels))
		for key := range secretInfo.Labels {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("  %s: %s\n", key, secretInfo.Labels[key])
		}
	} else {
		fmt.Println("Labels: None")
	}

	// Annotations (Tags)
	if len(secretInfo.Annotations) > 0 {
		fmt.Println("Tags (Annotations):")
		keys := make([]string, 0, len(secretInfo.Annotations))
		for key := range secretInfo.Annotations {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("  %s: %s\n", key, secretInfo.Annotations[key])
		}
	} else {
		fmt.Println("Tags (Annotations): None")
	}

	// Version aliases
	if len(secretInfo.VersionAliases) > 0 {
		fmt.Println("Version Aliases:")
		keys := make([]string, 0, len(secretInfo.VersionAliases))
		for key := range secretInfo.VersionAliases {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Printf("  %s: %s\n", key, secretInfo.VersionAliases[key])
		}
	}

	// Expiration information
	if secretInfo.ExpireTime != nil {
		fmt.Printf("Expires: %s\n", secretInfo.ExpireTime.Format(time.RFC3339))
	}
	if secretInfo.Ttl != "" {
		fmt.Printf("TTL: %s\n", secretInfo.Ttl)
	}

	// Rotation information
	if secretInfo.Rotation.NextRotationTime != nil {
		fmt.Printf("Next Rotation: %s\n", secretInfo.Rotation.NextRotationTime.Format(time.RFC3339))
	}
	if secretInfo.Rotation.RotationPeriod != "" {
		fmt.Printf("Rotation Period: %s\n", secretInfo.Rotation.RotationPeriod)
	}

	// Pub/Sub topics
	if len(secretInfo.Topics) > 0 {
		fmt.Println("Pub/Sub Topics:")
		for _, topic := range secretInfo.Topics {
			fmt.Printf("  %s\n", topic.Name)
		}
	}

	if showVersions {
		fmt.Println("\n--- All Versions ---")
		return listSecretVersions(secretName, project)
	}

	return nil
}

// extractVersionNumber extracts version number from full version name
func extractVersionNumber(versionName string) string {
	if parts := strings.Split(versionName, "/"); len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return versionName
}

// getReplicationStrategy returns a human-readable replication strategy
func getReplicationStrategy(replication struct {
	Automatic   interface{} `json:"automatic,omitempty"`
	UserManaged interface{} `json:"userManaged,omitempty"`
}) string {
	if replication.Automatic != nil {
		return "Automatic (multi-region)"
	}
	if replication.UserManaged != nil {
		return "User-managed"
	}
	return "Unknown"
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
