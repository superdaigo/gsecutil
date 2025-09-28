package cmd

import (
	"strings"
	"testing"
	"time"
)

// TestParseOperationFilter tests the parsing of operation filter
func TestParseOperationFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "Single valid operation",
			input:    "ACCESS",
			expected: []string{"ACCESS"},
		},
		{
			name:     "Multiple valid operations",
			input:    "ACCESS,CREATE,DELETE",
			expected: []string{"ACCESS", "CREATE", "DELETE"},
		},
		{
			name:     "Case insensitive operations",
			input:    "access,create,delete",
			expected: []string{"ACCESS", "CREATE", "DELETE"},
		},
		{
			name:     "Operations with spaces",
			input:    " ACCESS , CREATE , DELETE ",
			expected: []string{"ACCESS", "CREATE", "DELETE"},
		},
		{
			name:     "Mixed valid and invalid operations",
			input:    "ACCESS,INVALID,CREATE",
			expected: []string{"ACCESS", "CREATE"}, // Only valid ones
		},
		{
			name:     "All invalid operations",
			input:    "INVALID1,INVALID2",
			expected: []string{}, // None valid
		},
		{
			name:     "Empty elements",
			input:    "ACCESS,,CREATE",
			expected: []string{"ACCESS", "CREATE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOperationFilter(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseOperationFilter(%q) returned %d items, expected %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("parseOperationFilter(%q)[%d] = %q, expected %q", tt.input, i, result[i], expected)
				}
			}
		})
	}
}

// TestIsValidOperation tests operation validation
func TestIsValidOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		expected  bool
	}{
		{"Valid ACCESS", "ACCESS", true},
		{"Valid CREATE", "CREATE", true},
		{"Valid UPDATE", "UPDATE", true},
		{"Valid DELETE", "DELETE", true},
		{"Valid GET_METADATA", "GET_METADATA", true},
		{"Valid LIST", "LIST", true},
		{"Valid UPDATE_METADATA", "UPDATE_METADATA", true},
		{"Valid DESTROY_VERSION", "DESTROY_VERSION", true},
		{"Valid DISABLE_VERSION", "DISABLE_VERSION", true},
		{"Valid ENABLE_VERSION", "ENABLE_VERSION", true},
		{"Invalid operation", "INVALID", false},
		{"Empty string", "", false},
		{"Case sensitive - should be false for lowercase", "access", false},
		{"Partial match should be false", "ACC", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidOperation(tt.operation)
			if result != tt.expected {
				t.Errorf("isValidOperation(%q) = %v, expected %v", tt.operation, result, tt.expected)
			}
		})
	}
}

// TestIsSecretRelatedOperation tests secret-related operation detection
func TestIsSecretRelatedOperation(t *testing.T) {
	tests := []struct {
		name     string
		entry    AuditLogEntry
		expected bool
	}{
		{
			name: "AccessSecretVersion operation",
			entry: AuditLogEntry{
				ProtoPayload: struct {
					Type               string `json:"@type"`
					AuthenticationInfo struct {
						PrincipalEmail string `json:"principalEmail"`
					} `json:"authenticationInfo"`
					MethodName   string `json:"methodName"`
					ResourceName string `json:"resourceName"`
					Request      struct {
						Name string `json:"name"`
					} `json:"request"`
					Response struct {
						Name string `json:"name"`
					} `json:"response"`
				}{
					MethodName: "google.cloud.secretmanager.v1.SecretManagerService.AccessSecretVersion",
				},
			},
			expected: true,
		},
		{
			name: "CreateSecret operation",
			entry: AuditLogEntry{
				ProtoPayload: struct {
					Type               string `json:"@type"`
					AuthenticationInfo struct {
						PrincipalEmail string `json:"principalEmail"`
					} `json:"authenticationInfo"`
					MethodName   string `json:"methodName"`
					ResourceName string `json:"resourceName"`
					Request      struct {
						Name string `json:"name"`
					} `json:"request"`
					Response struct {
						Name string `json:"name"`
					} `json:"response"`
				}{
					MethodName: "google.cloud.secretmanager.v1.SecretManagerService.CreateSecret",
				},
			},
			expected: true,
		},
		{
			name: "ListSecrets with location (should be false)",
			entry: AuditLogEntry{
				ProtoPayload: struct {
					Type               string `json:"@type"`
					AuthenticationInfo struct {
						PrincipalEmail string `json:"principalEmail"`
					} `json:"authenticationInfo"`
					MethodName   string `json:"methodName"`
					ResourceName string `json:"resourceName"`
					Request      struct {
						Name string `json:"name"`
					} `json:"request"`
					Response struct {
						Name string `json:"name"`
					} `json:"response"`
				}{
					MethodName:   "google.cloud.secretmanager.v1.SecretManagerService.ListSecrets",
					ResourceName: "projects/test-project/locations/us-central1",
				},
			},
			expected: false,
		},
		{
			name: "ListSecrets project-level (should be true)",
			entry: AuditLogEntry{
				ProtoPayload: struct {
					Type               string `json:"@type"`
					AuthenticationInfo struct {
						PrincipalEmail string `json:"principalEmail"`
					} `json:"authenticationInfo"`
					MethodName   string `json:"methodName"`
					ResourceName string `json:"resourceName"`
					Request      struct {
						Name string `json:"name"`
					} `json:"request"`
					Response struct {
						Name string `json:"name"`
					} `json:"response"`
				}{
					MethodName:   "google.cloud.secretmanager.v1.SecretManagerService.ListSecrets",
					ResourceName: "projects/test-project",
				},
			},
			expected: true,
		},
		{
			name: "Unrelated operation",
			entry: AuditLogEntry{
				ProtoPayload: struct {
					Type               string `json:"@type"`
					AuthenticationInfo struct {
						PrincipalEmail string `json:"principalEmail"`
					} `json:"authenticationInfo"`
					MethodName   string `json:"methodName"`
					ResourceName string `json:"resourceName"`
					Request      struct {
						Name string `json:"name"`
					} `json:"request"`
					Response struct {
						Name string `json:"name"`
					} `json:"response"`
				}{
					MethodName: "google.cloud.other.v1.OtherService.SomeOperation",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSecretRelatedOperation(tt.entry)
			if result != tt.expected {
				t.Errorf("isSecretRelatedOperation() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetOperationName tests operation name extraction
func TestGetOperationName(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		expected   string
	}{
		{
			name:       "AccessSecretVersion",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.AccessSecretVersion",
			expected:   "ACCESS",
		},
		{
			name:       "CreateSecret",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.CreateSecret",
			expected:   "CREATE",
		},
		{
			name:       "AddSecretVersion",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.AddSecretVersion",
			expected:   "UPDATE",
		},
		{
			name:       "DeleteSecret",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.DeleteSecret",
			expected:   "DELETE",
		},
		{
			name:       "GetSecret",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.GetSecret",
			expected:   "GET_METADATA",
		},
		{
			name:       "ListSecrets",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.ListSecrets",
			expected:   "LIST",
		},
		{
			name:       "UpdateSecret",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.UpdateSecret",
			expected:   "UPDATE_METADATA",
		},
		{
			name:       "DestroySecretVersion",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.DestroySecretVersion",
			expected:   "DESTROY_VERSION",
		},
		{
			name:       "DisableSecretVersion",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.DisableSecretVersion",
			expected:   "DISABLE_VERSION",
		},
		{
			name:       "EnableSecretVersion",
			methodName: "google.cloud.secretmanager.v1.SecretManagerService.EnableSecretVersion",
			expected:   "ENABLE_VERSION",
		},
		{
			name:       "Unknown method",
			methodName: "google.cloud.other.v1.Service.UnknownMethod",
			expected:   "UnknownMethod",
		},
		{
			name:       "Simple method name",
			methodName: "SimpleMethod",
			expected:   "SimpleMethod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getOperationName(tt.methodName)
			if result != tt.expected {
				t.Errorf("getOperationName(%q) = %q, expected %q", tt.methodName, result, tt.expected)
			}
		})
	}
}

// TestFilterLogEntries tests the log entry filtering logic
func TestFilterLogEntries(t *testing.T) {
	// Create test entries
	entries := []AuditLogEntry{
		{
			Timestamp: time.Now(),
			ProtoPayload: struct {
				Type               string `json:"@type"`
				AuthenticationInfo struct {
					PrincipalEmail string `json:"principalEmail"`
				} `json:"authenticationInfo"`
				MethodName   string `json:"methodName"`
				ResourceName string `json:"resourceName"`
				Request      struct {
					Name string `json:"name"`
				} `json:"request"`
				Response struct {
					Name string `json:"name"`
				} `json:"response"`
			}{
				MethodName:   "google.cloud.secretmanager.v1.SecretManagerService.AccessSecretVersion",
				ResourceName: "projects/test/secrets/my-secret/versions/1",
				AuthenticationInfo: struct {
					PrincipalEmail string `json:"principalEmail"`
				}{
					PrincipalEmail: "user@example.com",
				},
			},
		},
		{
			Timestamp: time.Now(),
			ProtoPayload: struct {
				Type               string `json:"@type"`
				AuthenticationInfo struct {
					PrincipalEmail string `json:"principalEmail"`
				} `json:"authenticationInfo"`
				MethodName   string `json:"methodName"`
				ResourceName string `json:"resourceName"`
				Request      struct {
					Name string `json:"name"`
				} `json:"request"`
				Response struct {
					Name string `json:"name"`
				} `json:"response"`
			}{
				MethodName:   "google.cloud.secretmanager.v1.SecretManagerService.CreateSecret",
				ResourceName: "projects/test/secrets/other-secret",
				AuthenticationInfo: struct {
					PrincipalEmail string `json:"principalEmail"`
				}{
					PrincipalEmail: "admin@example.com",
				},
			},
		},
		{
			Timestamp: time.Now(),
			ProtoPayload: struct {
				Type               string `json:"@type"`
				AuthenticationInfo struct {
					PrincipalEmail string `json:"principalEmail"`
				} `json:"authenticationInfo"`
				MethodName   string `json:"methodName"`
				ResourceName string `json:"resourceName"`
				Request      struct {
					Name string `json:"name"`
				} `json:"request"`
				Response struct {
					Name string `json:"name"`
				} `json:"response"`
			}{
				MethodName:   "google.cloud.secretmanager.v1.SecretManagerService.ListSecrets",
				ResourceName: "projects/test/locations/us-central1", // Should be filtered out
			},
		},
	}

	tests := []struct {
		name           string
		secretName     string
		userFilter     string
		operations     []string
		expectedCount  int
		expectedSecret string
	}{
		{
			name:          "No filters - should return secret-related entries only",
			secretName:    "",
			userFilter:    "",
			operations:    nil,
			expectedCount: 2, // Two secret operations, location listing filtered out
		},
		{
			name:          "Filter by secret name",
			secretName:    "my-secret",
			userFilter:    "",
			operations:    nil,
			expectedCount: 1, // Only the ACCESS operation on my-secret
		},
		{
			name:          "Filter by user",
			secretName:    "",
			userFilter:    "user@example",
			operations:    nil,
			expectedCount: 1, // Only the user@example.com entry
		},
		{
			name:          "Filter by operation",
			secretName:    "",
			userFilter:    "",
			operations:    []string{"CREATE"},
			expectedCount: 1, // Only the CREATE operation
		},
		{
			name:          "Filter by multiple operations",
			secretName:    "",
			userFilter:    "",
			operations:    []string{"ACCESS", "CREATE"},
			expectedCount: 2, // Both ACCESS and CREATE operations
		},
		{
			name:          "Combined filters",
			secretName:    "my-secret",
			userFilter:    "user",
			operations:    []string{"ACCESS"},
			expectedCount: 1, // Should match the ACCESS operation on my-secret by user
		},
		{
			name:          "No matches",
			secretName:    "nonexistent",
			userFilter:    "",
			operations:    nil,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterLogEntries(entries, tt.secretName, tt.userFilter, tt.operations)
			if len(result) != tt.expectedCount {
				t.Errorf("filterLogEntries() returned %d entries, expected %d", len(result), tt.expectedCount)
			}
		})
	}
}

// TestBuildLogFilter tests log filter construction
func TestBuildLogFilter(t *testing.T) {
	tests := []struct {
		name        string
		secretName  string
		userFilter  string
		days        int
		contains    []string
		notContains []string
	}{
		{
			name:       "Basic filter",
			secretName: "",
			userFilter: "",
			days:       7,
			contains:   []string{`protoPayload.serviceName="secretmanager.googleapis.com"`},
		},
		{
			name:       "With secret name",
			secretName: "my-secret",
			userFilter: "",
			days:       7,
			contains:   []string{`protoPayload.resourceName:"my-secret"`},
		},
		{
			name:       "With user filter",
			secretName: "",
			userFilter: "user@example.com",
			days:       7,
			contains:   []string{`protoPayload.authenticationInfo.principalEmail:"user@example.com"`},
		},
		{
			name:       "Combined filters",
			secretName: "my-secret",
			userFilter: "user@example.com",
			days:       30,
			contains: []string{
				`protoPayload.serviceName="secretmanager.googleapis.com"`,
				`protoPayload.resourceName:"my-secret"`,
				`protoPayload.authenticationInfo.principalEmail:"user@example.com"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildLogFilter(tt.secretName, tt.userFilter, tt.days)

			// Check that required strings are present
			for _, required := range tt.contains {
				if !strings.Contains(result, required) {
					t.Errorf("buildLogFilter() result doesn't contain %q", required)
				}
			}

			// Check that forbidden strings are not present
			for _, forbidden := range tt.notContains {
				if strings.Contains(result, forbidden) {
					t.Errorf("buildLogFilter() result contains forbidden string %q", forbidden)
				}
			}

			// Check that timestamp filter is present
			if !strings.Contains(result, "timestamp>=") {
				t.Error("buildLogFilter() result doesn't contain timestamp filter")
			}
		})
	}
}
