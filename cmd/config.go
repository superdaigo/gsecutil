package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the gsecutil configuration file structure
type Config struct {
	Project     string           `yaml:"project,omitempty"`
	Prefix      string           `yaml:"prefix,omitempty"`
	List        ListConfig       `yaml:"list,omitempty"`
	Credentials []CredentialInfo `yaml:"credentials,omitempty"`
	Defaults    DefaultConfig    `yaml:"defaults,omitempty"`
}

// ListConfig contains configuration for the list command
type ListConfig struct {
	Attributes []string `yaml:"attributes,omitempty"`
}

// CredentialInfo contains metadata for a specific credential
type CredentialInfo struct {
	Name       string                 `yaml:"name"`
	Title      string                 `yaml:"title,omitempty"`
	Attributes map[string]interface{} `yaml:",inline"`
}

// DefaultConfig contains default settings for new secrets
type DefaultConfig struct {
	Labels map[string]string `yaml:"labels,omitempty"`
}

var (
	globalConfig   *Config
	configFilePath string
)

// LoadConfig loads configuration from the specified file or default location
func LoadConfig(customPath string) (*Config, error) {
	var configPath string

	if customPath != "" {
		configPath = customPath
	} else {
		configPath = getDefaultConfigPath()
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return empty config if no file exists
		return &Config{}, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", configPath, err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file '%s': %w", configPath, err)
	}

	// Store the config path for reference
	configFilePath = configPath

	return &config, nil
}

// getDefaultConfigPath returns the default configuration file path.
// Priority order:
// 1. Current directory: gsecutil.conf
// 2. Home directory: ~/.config/gsecutil/gsecutil.conf (or %USERPROFILE%\.config\gsecutil\gsecutil.conf on Windows)
func getDefaultConfigPath() string {
	// 1. Check current directory for gsecutil.conf
	cwd, err := os.Getwd()
	if err == nil {
		if path := filepath.Join(cwd, "gsecutil.conf"); fileExists(path) {
			return path
		}
	}

	// 2. Default to home directory config
	homeDir, _ := os.UserHomeDir()
	var configDir string

	switch runtime.GOOS {
	case "windows":
		// Use %USERPROFILE%\.config\gsecutil\gsecutil.conf on Windows
		configDir = filepath.Join(homeDir, ".config", "gsecutil")
	default:
		// Use $HOME/.config/gsecutil/gsecutil.conf on Unix-like systems
		configDir = filepath.Join(homeDir, ".config", "gsecutil")
	}

	return filepath.Join(configDir, "gsecutil.conf")
}

// fileExists checks if a file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && !info.IsDir()
}

// GetConfig returns the global configuration, loading it if not already loaded
func GetConfig() *Config {
	if globalConfig == nil {
		config, err := LoadConfig("")
		if err != nil {
			// On error, return empty config and continue
			// This ensures gsecutil works even with invalid config files
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			return &Config{}
		}
		globalConfig = config
	}
	return globalConfig
}

// SetCustomConfigPath loads configuration from a custom path
func SetCustomConfigPath(path string) error {
	config, err := LoadConfig(path)
	if err != nil {
		return err
	}
	globalConfig = config
	return nil
}

// GetProject returns the project ID from config, with priority order:
// 1. CLI parameter (passed as argument)
// 2. Configuration file
// 3. Environment variable GSECUTIL_PROJECT
// 4. gcloud default project
func GetProject(cliProject string) string {
	// 1. CLI parameter has highest priority
	if cliProject != "" {
		return cliProject
	}

	// 2. Configuration file
	config := GetConfig()
	if config.Project != "" {
		return config.Project
	}

	// 3. Environment variable
	if envProject := os.Getenv("GSECUTIL_PROJECT"); envProject != "" {
		return envProject
	}

	// 4. gcloud default (will be handled by gcloud CLI itself)
	return ""
}

// GetPrefix returns the secret name prefix from configuration
func GetPrefix() string {
	config := GetConfig()
	return config.Prefix
}

// GetCredentialInfo returns metadata for a specific credential
func GetCredentialInfo(name string) *CredentialInfo {
	config := GetConfig()

	for _, cred := range config.Credentials {
		if cred.Name == name {
			return &cred
		}
	}

	return nil
}

// GetListAttributes returns the list of attributes to show in list command
func GetListAttributes() []string {
	config := GetConfig()

	// If no credentials defined, don't show attributes by default
	if len(config.Credentials) == 0 {
		return nil
	}

	// If list.attributes is configured, use it
	if len(config.List.Attributes) > 0 {
		return config.List.Attributes
	}

	// Default: show title if credentials exist
	return []string{"title"}
}

// HasCredentialsConfig returns true if configuration file has credentials section
func HasCredentialsConfig() bool {
	config := GetConfig()
	return len(config.Credentials) > 0
}

// FilterSecretsByPrefix filters secret names based on configured prefix
func FilterSecretsByPrefix(secretName string) bool {
	prefix := GetPrefix()
	if prefix == "" {
		return true // No prefix filtering
	}
	return strings.HasPrefix(secretName, prefix)
}

// AddPrefixToSecretName adds the configured prefix to a secret name if not already present
func AddPrefixToSecretName(secretName string) string {
	prefix := GetPrefix()
	if prefix == "" {
		return secretName
	}

	// If already has prefix, return as-is
	if strings.HasPrefix(secretName, prefix) {
		return secretName
	}

	return prefix + secretName
}

// FilterCredentialsByAttributes filters credentials based on attribute values
func FilterCredentialsByAttributes(filters map[string]string) []CredentialInfo {
	config := GetConfig()
	var filtered []CredentialInfo

	for _, cred := range config.Credentials {
		matches := true

		for key, value := range filters {
			// Check if credential has this attribute and it matches
			if credValue, exists := cred.Attributes[key]; exists {
				// Convert to string for comparison
				credValueStr := fmt.Sprintf("%v", credValue)
				if credValueStr != value {
					matches = false
					break
				}
			} else if key == "title" && cred.Title != value {
				// Special case for title field
				matches = false
				break
			} else if key == "name" && cred.Name != value {
				// Special case for name field
				matches = false
				break
			} else {
				// Attribute doesn't exist in credential
				matches = false
				break
			}
		}

		if matches {
			filtered = append(filtered, cred)
		}
	}

	return filtered
}

// GetAttributeValue returns the value of a specific attribute for a credential
func GetAttributeValue(cred *CredentialInfo, attribute string) string {
	if cred == nil {
		return "(unknown)"
	}

	// Handle special fields
	switch attribute {
	case "name":
		return cred.Name
	case "title":
		if cred.Title != "" {
			return cred.Title
		}
		return "(no title)"
	}

	// Check in attributes map
	if value, exists := cred.Attributes[attribute]; exists {
		return fmt.Sprintf("%v", value)
	}

	return "(unknown)"
}

// ParseFilterAttributes parses filter string like "env=prod,owner=backend"
func ParseFilterAttributes(filterStr string) (map[string]string, error) {
	filters := make(map[string]string)

	if filterStr == "" {
		return filters, nil
	}

	pairs := strings.Split(filterStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid filter format: %s (expected key=value)", pair)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" || value == "" {
			return nil, fmt.Errorf("invalid filter format: %s (key and value cannot be empty)", pair)
		}

		filters[key] = value
	}

	return filters, nil
}

// ParseShowAttributes parses comma-separated attribute names
func ParseShowAttributes(attributesStr string) []string {
	if attributesStr == "" {
		return nil
	}

	var attributes []string
	for _, attr := range strings.Split(attributesStr, ",") {
		attr = strings.TrimSpace(attr)
		if attr != "" {
			attributes = append(attributes, attr)
		}
	}

	return attributes
}
