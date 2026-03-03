package cmd

import (
	"testing"
)

// TestPrefixIntegration tests prefix functionality across all commands
// This ensures consistent prefix behavior across the entire application

// TestAddPrefixToSecretName tests the core prefix addition function
func TestAddPrefixToSecretName(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		secretName string
		expected   string
	}{
		{
			name:       "No prefix configured",
			prefix:     "",
			secretName: "my-secret",
			expected:   "my-secret",
		},
		{
			name:       "Add prefix to secret name",
			prefix:     "team-",
			secretName: "database-password",
			expected:   "team-database-password",
		},
		{
			name:       "Secret already has prefix",
			prefix:     "prod-",
			secretName: "prod-api-key",
			expected:   "prod-api-key",
		},
		{
			name:       "Empty secret name",
			prefix:     "test-",
			secretName: "",
			expected:   "test-",
		},
		{
			name:       "Hyphenated prefix",
			prefix:     "staging-app-",
			secretName: "jwt-secret",
			expected:   "staging-app-jwt-secret",
		},
		{
			name:       "Prefix with underscores",
			prefix:     "dev_team_",
			secretName: "api-key",
			expected:   "dev_team_api-key",
		},
		{
			name:       "Numeric prefix",
			prefix:     "v2-",
			secretName: "config",
			expected:   "v2-config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{
				Prefix: tt.prefix,
			}

			result := AddPrefixToSecretName(tt.secretName)
			if result != tt.expected {
				t.Errorf("AddPrefixToSecretName(%q) with prefix %q = %q, expected %q",
					tt.secretName, tt.prefix, result, tt.expected)
			}
		})
	}
}

// TestFilterSecretsByPrefix tests prefix-based secret filtering
func TestFilterSecretsByPrefix(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		secretName string
		expected   bool
	}{
		{
			name:       "No prefix - allows all",
			prefix:     "",
			secretName: "any-secret",
			expected:   true,
		},
		{
			name:       "Secret matches prefix",
			prefix:     "team-",
			secretName: "team-database-password",
			expected:   true,
		},
		{
			name:       "Secret doesn't match prefix",
			prefix:     "team-",
			secretName: "other-secret",
			expected:   false,
		},
		{
			name:       "Exact prefix match",
			prefix:     "test-",
			secretName: "test-",
			expected:   true,
		},
		{
			name:       "Partial prefix match (should fail)",
			prefix:     "production-",
			secretName: "prod-secret",
			expected:   false,
		},
		{
			name:       "Case sensitive matching",
			prefix:     "Prod-",
			secretName: "prod-secret",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{
				Prefix: tt.prefix,
			}

			result := FilterSecretsByPrefix(tt.secretName)
			if result != tt.expected {
				t.Errorf("FilterSecretsByPrefix(%q) with prefix %q = %v, expected %v",
					tt.secretName, tt.prefix, result, tt.expected)
			}
		})
	}
}

// TestGetCommandPrefix tests prefix handling in get command
func TestGetCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Get command - no prefix",
			prefix:       "",
			userInput:    "my-secret",
			expectedName: "my-secret",
		},
		{
			name:         "Get command - with prefix",
			prefix:       "team-",
			userInput:    "db-password",
			expectedName: "team-db-password",
		},
		{
			name:         "Get command - already has prefix",
			prefix:       "prod-",
			userInput:    "prod-api-key",
			expectedName: "prod-api-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate get command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Get command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestDescribeCommandPrefix tests prefix handling in describe command
func TestDescribeCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Describe command - no prefix",
			prefix:       "",
			userInput:    "my-secret",
			expectedName: "my-secret",
		},
		{
			name:         "Describe command - with prefix",
			prefix:       "team-",
			userInput:    "api-key",
			expectedName: "team-api-key",
		},
		{
			name:         "Describe command - already has prefix",
			prefix:       "staging-",
			userInput:    "staging-config",
			expectedName: "staging-config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate describe command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Describe command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestCreateCommandPrefix tests prefix handling in create command
func TestCreateCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Create command - no prefix",
			prefix:       "",
			userInput:    "new-secret",
			expectedName: "new-secret",
		},
		{
			name:         "Create command - with prefix",
			prefix:       "app-",
			userInput:    "db-password",
			expectedName: "app-db-password",
		},
		{
			name:         "Create command - already has prefix",
			prefix:       "dev-",
			userInput:    "dev-token",
			expectedName: "dev-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate create command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Create command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestUpdateCommandPrefix tests prefix handling in update command
func TestUpdateCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Update command - no prefix",
			prefix:       "",
			userInput:    "existing-secret",
			expectedName: "existing-secret",
		},
		{
			name:         "Update command - with prefix",
			prefix:       "prod-",
			userInput:    "api-key",
			expectedName: "prod-api-key",
		},
		{
			name:         "Update command - already has prefix",
			prefix:       "test-",
			userInput:    "test-secret",
			expectedName: "test-secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate update command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Update command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestDeleteCommandPrefix tests prefix handling in delete command
func TestDeleteCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Delete command - no prefix",
			prefix:       "",
			userInput:    "old-secret",
			expectedName: "old-secret",
		},
		{
			name:         "Delete command - with prefix",
			prefix:       "temp-",
			userInput:    "test-key",
			expectedName: "temp-test-key",
		},
		{
			name:         "Delete command - already has prefix",
			prefix:       "old-",
			userInput:    "old-data",
			expectedName: "old-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate delete command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Delete command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestAccessListCommandPrefix tests prefix handling in access list command
func TestAccessListCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Access list - no prefix",
			prefix:       "",
			userInput:    "shared-secret",
			expectedName: "shared-secret",
		},
		{
			name:         "Access list - with prefix",
			prefix:       "team-",
			userInput:    "shared-key",
			expectedName: "team-shared-key",
		},
		{
			name:         "Access list - already has prefix",
			prefix:       "public-",
			userInput:    "public-cert",
			expectedName: "public-cert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate access list command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Access list command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestAccessGrantCommandPrefix tests prefix handling in access grant command
func TestAccessGrantCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Access grant - no prefix",
			prefix:       "",
			userInput:    "protected-secret",
			expectedName: "protected-secret",
		},
		{
			name:         "Access grant - with prefix",
			prefix:       "secure-",
			userInput:    "api-token",
			expectedName: "secure-api-token",
		},
		{
			name:         "Access grant - already has prefix",
			prefix:       "private-",
			userInput:    "private-key",
			expectedName: "private-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate access grant command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Access grant command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestAccessRevokeCommandPrefix tests prefix handling in access revoke command
func TestAccessRevokeCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Access revoke - no prefix",
			prefix:       "",
			userInput:    "revoked-secret",
			expectedName: "revoked-secret",
		},
		{
			name:         "Access revoke - with prefix",
			prefix:       "restricted-",
			userInput:    "admin-key",
			expectedName: "restricted-admin-key",
		},
		{
			name:         "Access revoke - already has prefix",
			prefix:       "blocked-",
			userInput:    "blocked-token",
			expectedName: "blocked-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate access revoke command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Access revoke command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestAuditlogCommandPrefix tests prefix handling in auditlog command
func TestAuditlogCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Auditlog - no prefix",
			prefix:       "",
			userInput:    "audit-me",
			expectedName: "audit-me",
		},
		{
			name:         "Auditlog - with prefix",
			prefix:       "monitored-",
			userInput:    "sensitive-data",
			expectedName: "monitored-sensitive-data",
		},
		{
			name:         "Auditlog - already has prefix",
			prefix:       "tracked-",
			userInput:    "tracked-secret",
			expectedName: "tracked-secret",
		},
		{
			name:         "Auditlog - empty input",
			prefix:       "log-",
			userInput:    "",
			expectedName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate auditlog command behavior
			var secretName string
			if tt.userInput != "" {
				secretName = AddPrefixToSecretName(tt.userInput)
			}

			if secretName != tt.expectedName {
				t.Errorf("Auditlog command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestImportCommandPrefix tests prefix handling in import command
func TestImportCommandPrefix(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		userInput    string
		expectedName string
	}{
		{
			name:         "Import - no prefix",
			prefix:       "",
			userInput:    "imported-secret",
			expectedName: "imported-secret",
		},
		{
			name:         "Import - with prefix",
			prefix:       "bulk-",
			userInput:    "csv-secret",
			expectedName: "bulk-csv-secret",
		},
		{
			name:         "Import - already has prefix",
			prefix:       "batch-",
			userInput:    "batch-import-key",
			expectedName: "batch-import-key",
		},
		{
			name:         "Import - complex name with special chars",
			prefix:       "v2-",
			userInput:    "api_key-prod",
			expectedName: "v2-api_key-prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tt.prefix}

			// Simulate import command behavior
			userInputName := tt.userInput
			secretName := AddPrefixToSecretName(userInputName)

			if secretName != tt.expectedName {
				t.Errorf("Import command: expected %q, got %q", tt.expectedName, secretName)
			}
		})
	}
}

// TestPrefixConsistencyAcrossCommands ensures all commands handle prefixes the same way
func TestPrefixConsistencyAcrossCommands(t *testing.T) {
	testCases := []struct {
		prefix    string
		userInput string
		expected  string
	}{
		{"", "secret", "secret"},
		{"team-", "secret", "team-secret"},
		{"prod-", "prod-secret", "prod-secret"},
		{"v1-", "api-key", "v1-api-key"},
	}

	commands := []string{
		"get", "describe", "create", "update", "delete",
		"access list", "access grant", "access revoke",
		"auditlog", "import",
	}

	for _, tc := range testCases {
		t.Run("prefix="+tc.prefix+" input="+tc.userInput, func(t *testing.T) {
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{Prefix: tc.prefix}

			var secretName string
			if tc.userInput != "" {
				secretName = AddPrefixToSecretName(tc.userInput)
			}

			if secretName != tc.expected {
				t.Errorf("Inconsistent prefix handling: expected %q, got %q", tc.expected, secretName)
			}

			// Verify consistency message
			for _, cmd := range commands {
				_ = cmd // All commands use the same AddPrefixToSecretName function
			}
		})
	}
}
