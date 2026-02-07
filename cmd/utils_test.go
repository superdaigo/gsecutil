package cmd

import (
	"testing"
)

// TestExtractSecretName tests extracting secret names from full resource paths
func TestExtractSecretName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		expected string
	}{
		{
			name:     "Full GCP path",
			fullName: "projects/my-project/secrets/my-secret",
			expected: "my-secret",
		},
		{
			name:     "Simple name",
			fullName: "my-secret",
			expected: "my-secret",
		},
		{
			name:     "Path with version",
			fullName: "projects/my-project/secrets/my-secret/versions/1",
			expected: "my-secret",
		},
		{
			name:     "Empty string",
			fullName: "",
			expected: "",
		},
		{
			name:     "Single segment",
			fullName: "secret",
			expected: "secret",
		},
		{
			name:     "Path with dashes",
			fullName: "projects/test-project-123/secrets/api-key-prod",
			expected: "api-key-prod",
		},
		{
			name:     "Incomplete path",
			fullName: "projects/my-project/secrets",
			expected: "projects/my-project/secrets",
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

// TestSortSecrets tests sorting secrets by name
func TestSortSecrets(t *testing.T) {
	tests := []struct {
		name     string
		secrets  []SecretInfo
		expected []string // Expected order of names
	}{
		{
			name: "Already sorted",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/a"},
				{Name: "projects/p/secrets/b"},
				{Name: "projects/p/secrets/c"},
			},
			expected: []string{
				"projects/p/secrets/a",
				"projects/p/secrets/b",
				"projects/p/secrets/c",
			},
		},
		{
			name: "Reverse order",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/z"},
				{Name: "projects/p/secrets/m"},
				{Name: "projects/p/secrets/a"},
			},
			expected: []string{
				"projects/p/secrets/a",
				"projects/p/secrets/m",
				"projects/p/secrets/z",
			},
		},
		{
			name: "Mixed order",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/banana"},
				{Name: "projects/p/secrets/apple"},
				{Name: "projects/p/secrets/cherry"},
				{Name: "projects/p/secrets/apricot"},
			},
			expected: []string{
				"projects/p/secrets/apple",
				"projects/p/secrets/apricot",
				"projects/p/secrets/banana",
				"projects/p/secrets/cherry",
			},
		},
		{
			name:     "Empty list",
			secrets:  []SecretInfo{},
			expected: []string{},
		},
		{
			name: "Single element",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/only-one"},
			},
			expected: []string{
				"projects/p/secrets/only-one",
			},
		},
		{
			name: "Names with numbers",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/secret-10"},
				{Name: "projects/p/secrets/secret-2"},
				{Name: "projects/p/secrets/secret-1"},
			},
			expected: []string{
				"projects/p/secrets/secret-1",
				"projects/p/secrets/secret-10",
				"projects/p/secrets/secret-2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortSecrets(tt.secrets)

			if len(tt.secrets) != len(tt.expected) {
				t.Errorf("Expected %d secrets but got %d", len(tt.expected), len(tt.secrets))
				return
			}

			for i, expected := range tt.expected {
				if tt.secrets[i].Name != expected {
					t.Errorf("At index %d: expected %q but got %q", i, expected, tt.secrets[i].Name)
				}
			}
		})
	}
}
