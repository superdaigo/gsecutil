package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestValidateHeader tests CSV header validation with duplicate detection
func TestValidateHeader(t *testing.T) {
	tests := []struct {
		name        string
		header      []string
		expectError bool
		errorMsg    string
		nameIdx     int
		valueIdx    int
	}{
		{
			name:        "Valid header with name and value",
			header:      []string{"name", "value", "title", "label:env"},
			expectError: false,
			nameIdx:     0,
			valueIdx:    1,
		},
		{
			name:        "Valid header with name only",
			header:      []string{"name", "title", "owner"},
			expectError: false,
			nameIdx:     0,
			valueIdx:    -1,
		},
		{
			name:        "Missing name column",
			header:      []string{"value", "title"},
			expectError: true,
			errorMsg:    "CSV must have 'name' column",
		},
		{
			name:        "Duplicate column names (case-insensitive)",
			header:      []string{"name", "value", "owner", "Owner"},
			expectError: true,
			errorMsg:    "duplicate column names: 'owner'",
		},
		{
			name:        "Multiple duplicate columns",
			header:      []string{"name", "value", "owner", "Owner", "label:env", "Label:Env"},
			expectError: true,
			errorMsg:    "duplicate column names",
		},
		{
			name:        "Empty column name",
			header:      []string{"name", "value", ""},
			expectError: true,
			errorMsg:    "empty column name at position 3",
		},
		{
			name:        "Valid with label columns",
			header:      []string{"name", "value", "label:env", "label:team"},
			expectError: false,
			nameIdx:     0,
			valueIdx:    1,
		},
		{
			name:        "Case-insensitive name and value",
			header:      []string{"Name", "Value", "Title"},
			expectError: false,
			nameIdx:     0,
			valueIdx:    1,
		},
		{
			name:        "Duplicate label columns",
			header:      []string{"name", "label:env", "Label:ENV"},
			expectError: true,
			errorMsg:    "duplicate column names: 'label:env'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameIdx, valueIdx, err := validateHeader(tt.header)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing %q but got none", tt.errorMsg)
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q but got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				if nameIdx != tt.nameIdx {
					t.Errorf("Expected nameIdx %d but got %d", tt.nameIdx, nameIdx)
				}
				if valueIdx != tt.valueIdx {
					t.Errorf("Expected valueIdx %d but got %d", tt.valueIdx, valueIdx)
				}
			}
		})
	}
}

// TestExtractColumnsData tests extracting labels, title, and attributes from CSV rows
func TestExtractColumnsData(t *testing.T) {
	tests := []struct {
		name           string
		header         []string
		record         []string
		nameIdx        int
		valueIdx       int
		expectedLabels map[string]string
		expectedTitle  string
		expectedAttrs  map[string]string
	}{
		{
			name:     "Extract labels only",
			header:   []string{"name", "value", "label:env", "label:team"},
			record:   []string{"test-secret", "secretvalue", "production", "backend"},
			nameIdx:  0,
			valueIdx: 1,
			expectedLabels: map[string]string{
				"env":  "production",
				"team": "backend",
			},
			expectedTitle: "",
			expectedAttrs: map[string]string{},
		},
		{
			name:           "Extract title and attributes",
			header:         []string{"name", "value", "title", "owner", "description"},
			record:         []string{"test-secret", "secretvalue", "Test Secret", "alice", "Test description"},
			nameIdx:        0,
			valueIdx:       1,
			expectedLabels: map[string]string{},
			expectedTitle:  "Test Secret",
			expectedAttrs: map[string]string{
				"owner":       "alice",
				"description": "Test description",
			},
		},
		{
			name:     "Extract all types",
			header:   []string{"name", "value", "title", "label:env", "owner", "label:team"},
			record:   []string{"test-secret", "secretvalue", "Test Title", "production", "bob", "frontend"},
			nameIdx:  0,
			valueIdx: 1,
			expectedLabels: map[string]string{
				"env":  "production",
				"team": "frontend",
			},
			expectedTitle: "Test Title",
			expectedAttrs: map[string]string{
				"owner": "bob",
			},
		},
		{
			name:           "Skip empty values",
			header:         []string{"name", "value", "title", "owner", "description"},
			record:         []string{"test-secret", "secretvalue", "", "alice", ""},
			nameIdx:        0,
			valueIdx:       1,
			expectedLabels: map[string]string{},
			expectedTitle:  "",
			expectedAttrs: map[string]string{
				"owner": "alice",
			},
		},
		{
			name:           "Case-insensitive title detection",
			header:         []string{"name", "value", "Title", "Owner"},
			record:         []string{"test-secret", "secretvalue", "My Title", "alice"},
			nameIdx:        0,
			valueIdx:       1,
			expectedLabels: map[string]string{},
			expectedTitle:  "My Title",
			expectedAttrs: map[string]string{
				"Owner": "alice",
			},
		},
		{
			name:           "No value column",
			header:         []string{"name", "title", "owner"},
			record:         []string{"test-secret", "Test Secret", "alice"},
			nameIdx:        0,
			valueIdx:       -1,
			expectedLabels: map[string]string{},
			expectedTitle:  "Test Secret",
			expectedAttrs: map[string]string{
				"owner": "alice",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labels, title, attrs := extractColumnsData(tt.header, tt.record, tt.nameIdx, tt.valueIdx)

			if !reflect.DeepEqual(labels, tt.expectedLabels) {
				t.Errorf("Labels = %+v, expected %+v", labels, tt.expectedLabels)
			}

			if title != tt.expectedTitle {
				t.Errorf("Title = %q, expected %q", title, tt.expectedTitle)
			}

			if !reflect.DeepEqual(attrs, tt.expectedAttrs) {
				t.Errorf("Attributes = %+v, expected %+v", attrs, tt.expectedAttrs)
			}
		})
	}
}

// TestUpdateConfigWithMetadata tests updating config with metadata from CSV
func TestUpdateConfigWithMetadata(t *testing.T) {
	tests := []struct {
		name            string
		initialConfig   *Config
		secretName      string
		title           string
		attributes      map[string]string
		expectedCredLen int
		expectedTitle   string
		expectedAttrs   map[string]interface{}
	}{
		{
			name: "Add new credential to empty config",
			initialConfig: &Config{
				Credentials: []CredentialInfo{},
			},
			secretName: "new-secret",
			title:      "New Secret",
			attributes: map[string]string{
				"owner": "alice",
			},
			expectedCredLen: 1,
			expectedTitle:   "New Secret",
			expectedAttrs: map[string]interface{}{
				"owner": "alice",
			},
		},
		{
			name: "Update existing credential",
			initialConfig: &Config{
				Credentials: []CredentialInfo{
					{
						Name:  "existing-secret",
						Title: "Old Title",
						Attributes: map[string]interface{}{
							"owner": "bob",
						},
					},
				},
			},
			secretName: "existing-secret",
			title:      "Updated Title",
			attributes: map[string]string{
				"owner":       "alice",
				"description": "New description",
			},
			expectedCredLen: 1,
			expectedTitle:   "Updated Title",
			expectedAttrs: map[string]interface{}{
				"owner":       "alice",
				"description": "New description",
			},
		},
		{
			name: "Add credential to existing config",
			initialConfig: &Config{
				Credentials: []CredentialInfo{
					{Name: "existing-secret", Title: "Existing"},
				},
			},
			secretName:      "new-secret",
			title:           "New Secret",
			attributes:      map[string]string{},
			expectedCredLen: 2,
			expectedTitle:   "New Secret",
			expectedAttrs:   map[string]interface{}{},
		},
		{
			name: "Empty title should update",
			initialConfig: &Config{
				Credentials: []CredentialInfo{},
			},
			secretName:      "test-secret",
			title:           "",
			attributes:      map[string]string{"owner": "alice"},
			expectedCredLen: 1,
			expectedTitle:   "",
			expectedAttrs: map[string]interface{}{
				"owner": "alice",
			},
		},
		{
			name: "Attributes are lowercased",
			initialConfig: &Config{
				Credentials: []CredentialInfo{},
			},
			secretName: "test-secret",
			title:      "Test",
			attributes: map[string]string{
				"Owner":       "alice",
				"Environment": "production",
			},
			expectedCredLen: 1,
			expectedTitle:   "Test",
			expectedAttrs: map[string]interface{}{
				"owner":       "alice",
				"environment": "production",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.initialConfig

			updateConfigWithMetadata(config, tt.secretName, tt.title, tt.attributes)

			if len(config.Credentials) != tt.expectedCredLen {
				t.Errorf("Expected %d credentials but got %d", tt.expectedCredLen, len(config.Credentials))
				return
			}

			// Find the credential
			var cred *CredentialInfo
			for i := range config.Credentials {
				if config.Credentials[i].Name == tt.secretName {
					cred = &config.Credentials[i]
					break
				}
			}

			if cred == nil {
				t.Errorf("Credential %q not found in config", tt.secretName)
				return
			}

			if cred.Title != tt.expectedTitle {
				t.Errorf("Title = %q, expected %q", cred.Title, tt.expectedTitle)
			}

			if !reflect.DeepEqual(cred.Attributes, tt.expectedAttrs) {
				t.Errorf("Attributes = %+v, expected %+v", cred.Attributes, tt.expectedAttrs)
			}
		})
	}
}

// TestReadCsvFile tests reading and parsing CSV files
func TestReadCsvFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "gsecutil-csv-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name            string
		csvContent      string
		expectedHeader  []string
		expectedRecords int
		expectError     bool
	}{
		{
			name: "Simple CSV",
			csvContent: `name,value,title
secret1,value1,Title 1
secret2,value2,Title 2`,
			expectedHeader:  []string{"name", "value", "title"},
			expectedRecords: 2,
			expectError:     false,
		},
		{
			name: "CSV with multi-line values (Excel format)",
			csvContent: `name,value,title
secret1,"line1
line2
line3",Multi-line Secret`,
			expectedHeader:  []string{"name", "value", "title"},
			expectedRecords: 1,
			expectError:     false,
		},
		{
			name: "CSV with labels and attributes",
			csvContent: `name,value,title,label:env,owner
secret1,value1,Title 1,production,alice
secret2,value2,Title 2,staging,bob`,
			expectedHeader:  []string{"name", "value", "title", "label:env", "owner"},
			expectedRecords: 2,
			expectError:     false,
		},
		{
			name:            "Empty CSV file",
			csvContent:      "",
			expectedHeader:  nil,
			expectedRecords: 0,
			expectError:     true,
		},
		{
			name:            "CSV with header only",
			csvContent:      `name,value,title`,
			expectedHeader:  []string{"name", "value", "title"},
			expectedRecords: 0,
			expectError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test CSV file
			csvFile := filepath.Join(tempDir, tt.name+".csv")
			if err := os.WriteFile(csvFile, []byte(tt.csvContent), 0644); err != nil {
				t.Fatalf("Failed to write test CSV file: %v", err)
			}

			records, header, err := readCsvFile(csvFile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(header, tt.expectedHeader) {
				t.Errorf("Header = %+v, expected %+v", header, tt.expectedHeader)
			}

			if len(records) != tt.expectedRecords {
				t.Errorf("Records count = %d, expected %d", len(records), tt.expectedRecords)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestPrepareCsvRecords tests CSV record preparation
func TestPrepareCsvRecords(t *testing.T) {
	tests := []struct {
		name         string
		secrets      []SecretInfo
		withValues   bool
		expectedCols []string // Expected header columns
	}{
		{
			name: "Simple secrets without values",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/s1"},
				{Name: "projects/p/secrets/s2"},
			},
			withValues:   false,
			expectedCols: []string{"name", "title"},
		},
		{
			name: "Secrets with labels",
			secrets: []SecretInfo{
				{
					Name: "projects/p/secrets/s1",
					Labels: map[string]string{
						"env": "production",
					},
				},
			},
			withValues:   false,
			expectedCols: []string{"name", "title", "label:env"},
		},
		{
			name: "With values flag",
			secrets: []SecretInfo{
				{Name: "projects/p/secrets/s1"},
			},
			withValues:   true,
			expectedCols: []string{"name", "value", "title"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock global config
			originalConfig := globalConfig
			defer func() { globalConfig = originalConfig }()
			globalConfig = &Config{Credentials: []CredentialInfo{}}

			records := prepareCsvRecords(tt.secrets, tt.withValues, "test-project")

			if len(records) == 0 {
				t.Error("Expected at least header row")
				return
			}

			header := records[0]
			for i, expectedCol := range tt.expectedCols {
				if i >= len(header) || header[i] != expectedCol {
					t.Errorf("Expected column %d to be %q but got %q", i, expectedCol, header[i])
				}
			}
		})
	}
}

// TestLoadOrCreateConfig tests config loading/creation
func TestLoadOrCreateConfig(t *testing.T) {
	// Save original config
	originalConfig := globalConfig
	defer func() { globalConfig = originalConfig }()

	// Reset global config
	globalConfig = nil

	config, err := loadOrCreateConfig()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if config == nil {
		t.Error("Expected config to be created")
	}
	// Config should either load existing or create new one with empty credentials
	// Just verify it returns successfully
}

// TestSaveConfig tests config file saving
func TestSaveConfig(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "gsecutil-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set HOME to temp directory to control config path
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Reset global config
	originalConfig := globalConfig
	defer func() { globalConfig = originalConfig }()
	globalConfig = nil

	// Create test config
	config := &Config{
		Project: "test-project",
		Prefix:  "test-",
		Credentials: []CredentialInfo{
			{
				Name:  "test-secret",
				Title: "Test Secret",
			},
		},
	}

	// Save config
	err = saveConfig(config)
	if err != nil {
		t.Errorf("Failed to save config: %v", err)
		return
	}

	// Verify file was created at default path
	expectedPath := filepath.Join(tempDir, ".config", "gsecutil", "gsecutil.conf")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at expected path: %s", expectedPath)
	}
}
