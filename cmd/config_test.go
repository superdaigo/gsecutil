package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestAddPrefixToSecretName tests prefix addition to secret names
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

// TestGetCredentialInfo tests credential lookup functionality
func TestGetCredentialInfo(t *testing.T) {
	tests := []struct {
		name        string
		credentials []CredentialInfo
		lookupName  string
		expected    *CredentialInfo
	}{
		{
			name: "Find existing credential",
			credentials: []CredentialInfo{
				{
					Name:  "database-password",
					Title: "Production Database Password",
					Attributes: map[string]interface{}{
						"environment": "production",
						"owner":       "backend-team",
					},
				},
				{
					Name:  "api-key",
					Title: "External API Key",
				},
			},
			lookupName: "database-password",
			expected: &CredentialInfo{
				Name:  "database-password",
				Title: "Production Database Password",
				Attributes: map[string]interface{}{
					"environment": "production",
					"owner":       "backend-team",
				},
			},
		},
		{
			name: "Credential not found",
			credentials: []CredentialInfo{
				{Name: "api-key", Title: "API Key"},
			},
			lookupName: "nonexistent-secret",
			expected:   nil,
		},
		{
			name:        "Empty credentials list",
			credentials: []CredentialInfo{},
			lookupName:  "any-secret",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &Config{
				Credentials: tt.credentials,
			}

			result := GetCredentialInfo(tt.lookupName)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil but got %+v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("Expected %+v but got nil", tt.expected)
				return
			}

			if result.Name != tt.expected.Name || result.Title != tt.expected.Title {
				t.Errorf("GetCredentialInfo(%q) = %+v, expected %+v", tt.lookupName, result, tt.expected)
			}

			// Compare attributes if they exist
			if !reflect.DeepEqual(result.Attributes, tt.expected.Attributes) {
				t.Errorf("GetCredentialInfo(%q) attributes = %+v, expected %+v",
					tt.lookupName, result.Attributes, tt.expected.Attributes)
			}
		})
	}
}

// TestGetAttributeValue tests attribute value retrieval
func TestGetAttributeValue(t *testing.T) {
	tests := []struct {
		name      string
		cred      *CredentialInfo
		attribute string
		expected  string
	}{
		{
			name: "Get title from credential",
			cred: &CredentialInfo{
				Name:  "test-secret",
				Title: "Test Secret Title",
			},
			attribute: "title",
			expected:  "Test Secret Title",
		},
		{
			name: "Get name from credential",
			cred: &CredentialInfo{
				Name:  "my-secret",
				Title: "My Secret",
			},
			attribute: "name",
			expected:  "my-secret",
		},
		{
			name: "Get attribute from attributes map",
			cred: &CredentialInfo{
				Name: "api-secret",
				Attributes: map[string]interface{}{
					"environment": "production",
					"owner":       "backend-team",
					"rotation":    30,
				},
			},
			attribute: "environment",
			expected:  "production",
		},
		{
			name: "Get numeric attribute",
			cred: &CredentialInfo{
				Name: "api-secret",
				Attributes: map[string]interface{}{
					"rotation_days": 30,
				},
			},
			attribute: "rotation_days",
			expected:  "30",
		},
		{
			name: "Get unknown attribute",
			cred: &CredentialInfo{
				Name:  "test-secret",
				Title: "Test",
			},
			attribute: "nonexistent",
			expected:  "(unknown)",
		},
		{
			name: "Get title from credential with no title",
			cred: &CredentialInfo{
				Name: "no-title-secret",
			},
			attribute: "title",
			expected:  "(no title)",
		},
		{
			name:      "Nil credential",
			cred:      nil,
			attribute: "title",
			expected:  "(unknown)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAttributeValue(tt.cred, tt.attribute)
			if result != tt.expected {
				t.Errorf("GetAttributeValue(%+v, %q) = %q, expected %q",
					tt.cred, tt.attribute, result, tt.expected)
			}
		})
	}
}

// TestParseFilterAttributes tests parsing of filter attribute strings
func TestParseFilterAttributes(t *testing.T) {
	tests := []struct {
		name      string
		filterStr string
		expected  map[string]string
		wantErr   bool
	}{
		{
			name:      "Empty filter string",
			filterStr: "",
			expected:  map[string]string{},
			wantErr:   false,
		},
		{
			name:      "Single filter",
			filterStr: "environment=production",
			expected:  map[string]string{"environment": "production"},
			wantErr:   false,
		},
		{
			name:      "Multiple filters",
			filterStr: "environment=production,owner=backend-team",
			expected: map[string]string{
				"environment": "production",
				"owner":       "backend-team",
			},
			wantErr: false,
		},
		{
			name:      "Filters with spaces",
			filterStr: " environment = production , owner = backend-team ",
			expected: map[string]string{
				"environment": "production",
				"owner":       "backend-team",
			},
			wantErr: false,
		},
		{
			name:      "Invalid format - no equals",
			filterStr: "environment:production",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Invalid format - empty key",
			filterStr: "=production",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Invalid format - empty value",
			filterStr: "environment=",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "Mixed valid and invalid",
			filterStr: "environment=production,invalid",
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseFilterAttributes(tt.filterStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseFilterAttributes(%q) = %+v, expected %+v",
					tt.filterStr, result, tt.expected)
			}
		})
	}
}

// TestParseShowAttributes tests parsing of show attributes string
func TestParseShowAttributes(t *testing.T) {
	tests := []struct {
		name          string
		attributesStr string
		expected      []string
	}{
		{
			name:          "Empty string",
			attributesStr: "",
			expected:      nil,
		},
		{
			name:          "Single attribute",
			attributesStr: "title",
			expected:      []string{"title"},
		},
		{
			name:          "Multiple attributes",
			attributesStr: "title,owner,environment",
			expected:      []string{"title", "owner", "environment"},
		},
		{
			name:          "Attributes with spaces",
			attributesStr: " title , owner , environment ",
			expected:      []string{"title", "owner", "environment"},
		},
		{
			name:          "Empty attributes filtered out",
			attributesStr: "title,,owner,",
			expected:      []string{"title", "owner"},
		},
		{
			name:          "Single attribute with spaces",
			attributesStr: "  title  ",
			expected:      []string{"title"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseShowAttributes(tt.attributesStr)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseShowAttributes(%q) = %+v, expected %+v",
					tt.attributesStr, result, tt.expected)
			}
		})
	}
}

// TestGetListAttributes tests default list attributes configuration
func TestGetListAttributes(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected []string
	}{
		{
			name: "No credentials - empty list",
			config: Config{
				Credentials: []CredentialInfo{},
			},
			expected: nil,
		},
		{
			name: "Has credentials but no list config - default title",
			config: Config{
				Credentials: []CredentialInfo{
					{Name: "test-secret"},
				},
			},
			expected: []string{"title"},
		},
		{
			name: "Has credentials and list config",
			config: Config{
				List: ListConfig{
					Attributes: []string{"title", "owner", "environment"},
				},
				Credentials: []CredentialInfo{
					{Name: "test-secret"},
				},
			},
			expected: []string{"title", "owner", "environment"},
		},
		{
			name: "Empty list config with credentials",
			config: Config{
				List: ListConfig{
					Attributes: []string{},
				},
				Credentials: []CredentialInfo{
					{Name: "test-secret"},
				},
			},
			expected: []string{"title"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &tt.config

			result := GetListAttributes()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetListAttributes() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}

// TestConfigFileLoading tests loading configuration from files
func TestConfigFileLoading(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "gsecutil-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `# Test configuration
project: "test-project-123"
prefix: "team-"

list:
  attributes:
    - title
    - owner
    - environment

credentials:
  - name: "database-password"
    title: "Production Database Password"
    environment: "production"
    owner: "backend-team"

  - name: "api-key"
    title: "External API Key"
    environment: "production"
    owner: "api-team"
`

	configFile := filepath.Join(tempDir, "test-config.conf")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test loading the config
	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config values
	if config.Project != "test-project-123" {
		t.Errorf("Expected project 'test-project-123', got %q", config.Project)
	}

	if config.Prefix != "team-" {
		t.Errorf("Expected prefix 'team-', got %q", config.Prefix)
	}

	expectedAttrs := []string{"title", "owner", "environment"}
	if !reflect.DeepEqual(config.List.Attributes, expectedAttrs) {
		t.Errorf("Expected list attributes %+v, got %+v", expectedAttrs, config.List.Attributes)
	}

	if len(config.Credentials) != 2 {
		t.Errorf("Expected 2 credentials, got %d", len(config.Credentials))
	}

	// Verify first credential
	cred1 := config.Credentials[0]
	if cred1.Name != "database-password" {
		t.Errorf("Expected credential name 'database-password', got %q", cred1.Name)
	}
	if cred1.Title != "Production Database Password" {
		t.Errorf("Expected credential title 'Production Database Password', got %q", cred1.Title)
	}
}

// TestConfigFileLoadingErrors tests error handling in config loading
func TestConfigFileLoadingErrors(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		expectError bool
	}{
		{
			name:        "Nonexistent file",
			configFile:  "/path/to/nonexistent/config.conf",
			expectError: false, // Should return empty config, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfig(tt.configFile)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && config == nil {
				t.Error("Expected config but got nil")
			}
		})
	}
}

// TestGetProject tests project resolution with different priorities
func TestGetProject(t *testing.T) {
	tests := []struct {
		name       string
		cliProject string
		config     Config
		envVar     string
		expected   string
	}{
		{
			name:       "CLI parameter takes precedence",
			cliProject: "cli-project",
			config:     Config{Project: "config-project"},
			envVar:     "env-project",
			expected:   "cli-project",
		},
		{
			name:       "Config file used when no CLI parameter",
			cliProject: "",
			config:     Config{Project: "config-project"},
			envVar:     "env-project",
			expected:   "config-project",
		},
		{
			name:       "Environment variable used when no CLI or config",
			cliProject: "",
			config:     Config{},
			envVar:     "env-project",
			expected:   "env-project",
		},
		{
			name:       "Empty when no project specified anywhere",
			cliProject: "",
			config:     Config{},
			envVar:     "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()

			globalConfig = &tt.config

			// Set environment variable
			if tt.envVar != "" {
				os.Setenv("GOOGLE_CLOUD_PROJECT", tt.envVar)
				defer os.Unsetenv("GOOGLE_CLOUD_PROJECT")
			}

			result := GetProject(tt.cliProject)
			if result != tt.expected {
				t.Errorf("GetProject(%q) = %q, expected %q", tt.cliProject, result, tt.expected)
			}
		})
	}
}
