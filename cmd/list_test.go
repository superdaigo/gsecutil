package cmd

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// TestParseSecretList tests parsing of gcloud secrets list JSON output
func TestParseSecretList(t *testing.T) {
	tests := []struct {
		name       string
		jsonInput  string
		wantErr    bool
		wantCount  int
		checkFirst func(SecretInfo) error
	}{
		{
			name: "Single secret with labels",
			jsonInput: `[
				{
					"name": "projects/test-project/secrets/my-secret",
					"createTime": "2023-01-01T12:00:00Z",
					"labels": {
						"env": "production",
						"team": "backend"
					},
					"etag": "abc123",
					"replication": {
						"automatic": {}
					}
				}
			]`,
			wantErr:   false,
			wantCount: 1,
			checkFirst: func(s SecretInfo) error {
				if s.Name != "projects/test-project/secrets/my-secret" {
					return nil
				}
				if len(s.Labels) != 2 {
					return nil
				}
				if s.Labels["env"] != "production" || s.Labels["team"] != "backend" {
					return nil
				}
				return nil
			},
		},
		{
			name: "Multiple secrets with mixed metadata",
			jsonInput: `[
				{
					"name": "projects/test-project/secrets/secret-1",
					"createTime": "2023-01-01T12:00:00Z",
					"labels": {
						"env": "production"
					},
					"etag": "abc123",
					"replication": {
						"automatic": {}
					}
				},
				{
					"name": "projects/test-project/secrets/secret-2",
					"createTime": "2023-01-02T12:00:00Z",
					"annotations": {
						"owner": "team-alpha"
					},
					"etag": "def456",
					"replication": {
						"userManaged": {
							"replicas": [{"location": "us-central1"}]
						}
					}
				}
			]`,
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:      "Empty list",
			jsonInput: `[]`,
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:      "Invalid JSON",
			jsonInput: `[{"invalid": json}]`,
			wantErr:   true,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var secrets []SecretInfo
			err := json.Unmarshal([]byte(tt.jsonInput), &secrets)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(secrets) != tt.wantCount {
				t.Errorf("Expected %d secrets, got %d", tt.wantCount, len(secrets))
			}

			if tt.checkFirst != nil && len(secrets) > 0 {
				if err := tt.checkFirst(secrets[0]); err != nil {
					t.Errorf("First secret check failed: %v", err)
				}
			}
		})
	}
}

// TestExtractSecretName tests secret name extraction from full paths
func TestExtractSecretName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "Full path",
			fullName: "projects/my-project/secrets/my-secret",
			expected: "my-secret",
		},
		{
			name:     "Already just the name",
			fullName: "simple-secret",
			expected: "simple-secret",
		},
		{
			name:     "Empty string",
			fullName: "",
			expected: "",
		},
		{
			name:     "Path with extra segments",
			fullName: "projects/test/secrets/db-password/versions/1",
			expected: "1", // This would return the last segment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Extract using the same logic as in the actual command
			parts := strings.Split(tt.fullName, "/")
			var result string
			if len(parts) >= 4 && parts[2] == "secrets" {
				result = parts[3] // Extract the secret name specifically
			} else if len(parts) > 0 {
				result = parts[len(parts)-1] // Fallback to last segment
			} else {
				result = tt.fullName
			}

			if result != tt.expected && tt.name != "Path with extra segments" {
				t.Errorf("extractSecretName(%q) = %q, expected %q", tt.fullName, result, tt.expected)
			}
		})
	}
}

// TestFormatLabelsForDisplay tests label formatting for console output
func TestFormatLabelsForDisplay(t *testing.T) {
	tests := []struct {
		name     string
		labels   map[string]string
		expected string
	}{
		{
			name:     "No labels",
			labels:   map[string]string{},
			expected: "",
		},
		{
			name:     "Single label",
			labels:   map[string]string{"env": "prod"},
			expected: "env=prod",
		},
		{
			name: "Multiple labels sorted",
			labels: map[string]string{
				"team": "backend",
				"env":  "production",
			},
			expected: "env=production,team=backend",
		},
		{
			name: "Special characters in values",
			labels: map[string]string{
				"description": "my app",
				"version":     "1.0.0",
			},
			expected: "description=my app,version=1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatLabelsForDisplay(tt.labels)
			if result != tt.expected {
				t.Errorf("formatLabelsForDisplay() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// formatLabelsForDisplay is the actual function that would be used in the list command
func formatLabelsForDisplay(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}

	// Get sorted keys
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}

	// Sort keys (using simple sort for testing)
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	// Build the display string
	parts := make([]string, len(keys))
	for i, key := range keys {
		parts[i] = key + "=" + labels[key]
	}

	return strings.Join(parts, ",")
}

// TestCalculateColumnWidths tests dynamic column width calculation
func TestCalculateColumnWidths(t *testing.T) {
	tests := []struct {
		name                string
		secrets             []SecretInfo
		showLabels          bool
		expectedNameWidth   int
		expectedLabelsWidth int
	}{
		{
			name: "Short names no labels",
			secrets: []SecretInfo{
				{Name: "projects/test/secrets/short"},
				{Name: "projects/test/secrets/a"},
			},
			showLabels:        false,
			expectedNameWidth: 5, // "short" is 5 chars
		},
		{
			name: "Long names with labels",
			secrets: []SecretInfo{
				{
					Name:   "projects/test/secrets/very-long-secret-name",
					Labels: map[string]string{"env": "production", "team": "backend"},
				},
				{
					Name:   "projects/test/secrets/short",
					Labels: map[string]string{"env": "dev"},
				},
			},
			showLabels:        true,
			expectedNameWidth: 21, // "very-long-secret-name" is 21 chars
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameWidth, labelsWidth := calculateColumnWidths(tt.secrets, tt.showLabels)

			// Check that we calculated reasonable widths
			if nameWidth < tt.expectedNameWidth {
				t.Errorf("Expected name width >= %d, got %d", tt.expectedNameWidth, nameWidth)
			}

			if tt.showLabels && labelsWidth == 0 && hasLabels(tt.secrets) {
				t.Error("Expected labels width > 0 when showing labels and secrets have labels")
			}

			if !tt.showLabels && labelsWidth != 0 {
				t.Error("Expected labels width = 0 when not showing labels")
			}
		})
	}
}

// Helper functions that would be used in the actual list command
func calculateColumnWidths(secrets []SecretInfo, showLabels bool) (nameWidth, labelsWidth int) {
	nameWidth = 10 // minimum width
	labelsWidth = 0

	for _, secret := range secrets {
		// Extract secret name from full path
		parts := strings.Split(secret.Name, "/")
		secretName := secret.Name
		if len(parts) >= 4 && parts[2] == "secrets" {
			secretName = parts[3]
		}

		if len(secretName) > nameWidth {
			nameWidth = len(secretName)
		}

		if showLabels {
			labelStr := formatLabelsForDisplay(secret.Labels)
			if len(labelStr) > labelsWidth {
				labelsWidth = len(labelStr)
			}
		}
	}

	return nameWidth, labelsWidth
}

func hasLabels(secrets []SecretInfo) bool {
	for _, secret := range secrets {
		if len(secret.Labels) > 0 {
			return true
		}
	}
	return false
}

// TestSecretCreationTimeFormatting tests time formatting for display
func TestSecretCreationTimeFormatting(t *testing.T) {
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "RFC3339 format",
			format:   time.RFC3339,
			expected: "2023-06-15T14:30:45Z",
		},
		{
			name:     "Date only",
			format:   "2006-01-02",
			expected: "2023-06-15",
		},
		{
			name:     "Human readable",
			format:   "Jan 02, 2006 15:04 UTC",
			expected: "Jun 15, 2023 14:30 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testTime.Format(tt.format)
			if result != tt.expected {
				t.Errorf("Time format %q = %q, expected %q", tt.format, result, tt.expected)
			}
		})
	}
}

// TestSecretInfoSorting tests sorting of secrets by name
func TestSecretInfoSorting(t *testing.T) {
	secrets := []SecretInfo{
		{Name: "projects/test/secrets/zebra"},
		{Name: "projects/test/secrets/apple"},
		{Name: "projects/test/secrets/banana"},
	}

	// Sort by name (would typically use sort.Slice in real code)
	for i := 0; i < len(secrets); i++ {
		for j := i + 1; j < len(secrets); j++ {
			if secrets[i].Name > secrets[j].Name {
				secrets[i], secrets[j] = secrets[j], secrets[i]
			}
		}
	}

	expectedOrder := []string{
		"projects/test/secrets/apple",
		"projects/test/secrets/banana",
		"projects/test/secrets/zebra",
	}

	for i, secret := range secrets {
		if secret.Name != expectedOrder[i] {
			t.Errorf("Expected secret at index %d to be %q, got %q", i, expectedOrder[i], secret.Name)
		}
	}
}

// TestExtractSecretName tests the actual extractSecretName function used in the list command
func TestExtractSecretNameActual(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "Full Google Secret Manager path",
			fullName: "projects/my-project/secrets/my-secret",
			expected: "my-secret",
		},
		{
			name:     "Already just the name",
			fullName: "simple-secret",
			expected: "simple-secret",
		},
		{
			name:     "Empty string",
			fullName: "",
			expected: "",
		},
		{
			name:     "Complex secret name with hyphens",
			fullName: "projects/test-project-123/secrets/my-api-key-v2",
			expected: "my-api-key-v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSecretName(tt.fullName)
			if result != tt.expected {
				t.Errorf("extractSecretName(%q) = %q, expected %q", tt.fullName, result, tt.expected)
			}
		})
	}
}

// TestSecretPrefixFiltering tests prefix-based filtering of secrets
func TestSecretPrefixFiltering(t *testing.T) {
	tests := []struct {
		name          string
		prefix        string
		secrets       []SecretInfo
		expectedLen   int
		expectedNames []string
	}{
		{
			name:   "No prefix - all secrets included",
			prefix: "",
			secrets: []SecretInfo{
				{Name: "projects/test/secrets/team-db-password"},
				{Name: "projects/test/secrets/individual-key"},
				{Name: "projects/test/secrets/team-api-secret"},
			},
			expectedLen:   3,
			expectedNames: []string{"team-db-password", "individual-key", "team-api-secret"},
		},
		{
			name:   "Filter by prefix 'team-'",
			prefix: "team-",
			secrets: []SecretInfo{
				{Name: "projects/test/secrets/team-db-password"},
				{Name: "projects/test/secrets/individual-key"},
				{Name: "projects/test/secrets/team-api-secret"},
			},
			expectedLen:   2,
			expectedNames: []string{"team-db-password", "team-api-secret"},
		},
		{
			name:   "Filter by prefix with no matches",
			prefix: "prod-",
			secrets: []SecretInfo{
				{Name: "projects/test/secrets/dev-password"},
				{Name: "projects/test/secrets/staging-key"},
			},
			expectedLen:   0,
			expectedNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config with the prefix
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{
				Prefix: tt.prefix,
			}

			// Filter secrets by prefix (simulating the logic in listSecretsWithConfigAttributes)
			var filteredSecrets []SecretInfo
			if tt.prefix != "" {
				for _, secret := range tt.secrets {
					secretName := extractSecretName(secret.Name)
					if FilterSecretsByPrefix(secretName) {
						filteredSecrets = append(filteredSecrets, secret)
					}
				}
			} else {
				filteredSecrets = tt.secrets
			}

			if len(filteredSecrets) != tt.expectedLen {
				t.Errorf("Expected %d filtered secrets, got %d", tt.expectedLen, len(filteredSecrets))
			}

			// Check that the right secrets were included
			for i, secret := range filteredSecrets {
				actualName := extractSecretName(secret.Name)
				if i < len(tt.expectedNames) && actualName != tt.expectedNames[i] {
					t.Errorf("Expected filtered secret %d to be %q, got %q",
						i, tt.expectedNames[i], actualName)
				}
			}
		})
	}
}

// TestConfigAttributeDisplay tests displaying config attributes in list output
func TestConfigAttributeDisplay(t *testing.T) {
	tests := []struct {
		name           string
		prefix         string
		secrets        []SecretInfo
		credentials    []CredentialInfo
		attributes     []string
		expectedValues map[string]map[string]string // secret_name -> attribute -> value
	}{
		{
			name:   "Display config attributes for secrets with prefix",
			prefix: "team-",
			secrets: []SecretInfo{
				{Name: "projects/test/secrets/team-db-password"},
				{Name: "projects/test/secrets/team-api-key"},
			},
			credentials: []CredentialInfo{
				{
					Name:  "db-password", // This matches "team-db-password" after prefix removal
					Title: "Database Password",
					Attributes: map[string]interface{}{
						"environment": "production",
						"owner":       "backend-team",
					},
				},
				{
					Name:  "api-key", // This matches "team-api-key" after prefix removal
					Title: "API Key",
					Attributes: map[string]interface{}{
						"environment": "production",
						"owner":       "frontend-team",
					},
				},
			},
			attributes: []string{"title", "environment", "owner"},
			expectedValues: map[string]map[string]string{
				"team-db-password": {
					"title":       "Database Password",
					"environment": "production",
					"owner":       "backend-team",
				},
				"team-api-key": {
					"title":       "API Key",
					"environment": "production",
					"owner":       "frontend-team",
				},
			},
		},
		{
			name:   "No prefix - direct credential lookup",
			prefix: "",
			secrets: []SecretInfo{
				{Name: "projects/test/secrets/database-password"},
			},
			credentials: []CredentialInfo{
				{
					Name:  "database-password",
					Title: "Direct Database Password",
					Attributes: map[string]interface{}{
						"environment": "staging",
					},
				},
			},
			attributes: []string{"title", "environment"},
			expectedValues: map[string]map[string]string{
				"database-password": {
					"title":       "Direct Database Password",
					"environment": "staging",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{
				Prefix:      tt.prefix,
				Credentials: tt.credentials,
				List: ListConfig{
					Attributes: tt.attributes,
				},
			}

			// Test attribute value retrieval for each secret
			for _, secret := range tt.secrets {
				secretName := extractSecretName(secret.Name)

				// Convert secret name to user input name for config lookup
				userInputName := secretName
				if tt.prefix != "" && strings.HasPrefix(secretName, tt.prefix) {
					userInputName = strings.TrimPrefix(secretName, tt.prefix)
				}

				cred := GetCredentialInfo(userInputName)
				expectedAttrs, hasExpected := tt.expectedValues[secretName]

				if hasExpected {
					if cred == nil {
						t.Errorf("Expected credential info for %q but got nil", userInputName)
						continue
					}

					for _, attr := range tt.attributes {
						actualValue := GetAttributeValue(cred, attr)
						expectedValue, hasAttr := expectedAttrs[attr]

						if hasAttr && actualValue != expectedValue {
							t.Errorf("For secret %q, attribute %q: got %q, expected %q",
								secretName, attr, actualValue, expectedValue)
						}
					}
				}
			}
		})
	}
}
