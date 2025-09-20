package cmd

import (
	"encoding/json"
	"testing"
	"time"
)

// TestExtractVersionNumber tests version number extraction from full version names
func TestExtractVersionNumber(t *testing.T) {
	tests := []struct {
		name        string
		versionName string
		expected    string
	}{
		{
			name:        "Full version path",
			versionName: "projects/my-project/secrets/my-secret/versions/1",
			expected:    "1",
		},
		{
			name:        "Latest version",
			versionName: "projects/my-project/secrets/my-secret/versions/latest",
			expected:    "latest",
		},
		{
			name:        "Numeric version",
			versionName: "projects/my-project/secrets/my-secret/versions/42",
			expected:    "42",
		},
		{
			name:        "Simple version name",
			versionName: "1",
			expected:    "1",
		},
		{
			name:        "Empty string",
			versionName: "",
			expected:    "",
		},
		{
			name:        "Invalid format but has slashes",
			versionName: "some/path/5",
			expected:    "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVersionNumber(tt.versionName)
			if result != tt.expected {
				t.Errorf("extractVersionNumber(%q) = %q, expected %q", tt.versionName, result, tt.expected)
			}
		})
	}
}

// TestGetReplicationStrategy tests replication strategy detection
func TestGetReplicationStrategy(t *testing.T) {
	tests := []struct {
		name        string
		replication struct {
			Automatic   interface{} `json:"automatic,omitempty"`
			UserManaged interface{} `json:"userManaged,omitempty"`
		}
		expected string
	}{
		{
			name: "Automatic replication",
			replication: struct {
				Automatic   interface{} `json:"automatic,omitempty"`
				UserManaged interface{} `json:"userManaged,omitempty"`
			}{
				Automatic: map[string]interface{}{},
			},
			expected: "Automatic (multi-region)",
		},
		{
			name: "User-managed replication",
			replication: struct {
				Automatic   interface{} `json:"automatic,omitempty"`
				UserManaged interface{} `json:"userManaged,omitempty"`
			}{
				UserManaged: map[string]interface{}{
					"replicas": []interface{}{},
				},
			},
			expected: "User-managed",
		},
		{
			name: "No replication specified",
			replication: struct {
				Automatic   interface{} `json:"automatic,omitempty"`
				UserManaged interface{} `json:"userManaged,omitempty"`
			}{},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getReplicationStrategy(tt.replication)
			if result != tt.expected {
				t.Errorf("getReplicationStrategy() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestSecretInfoJSONParsing tests that SecretInfo can properly parse JSON responses
func TestSecretInfoJSONParsing(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected SecretInfo
		wantErr  bool
	}{
		{
			name: "Basic secret with automatic replication",
			jsonData: `{
				"name": "projects/test-project/secrets/my-secret",
				"createTime": "2023-01-01T12:00:00Z",
				"labels": {
					"env": "test",
					"team": "dev"
				},
				"etag": "abc123",
				"replication": {
					"automatic": {}
				}
			}`,
			expected: SecretInfo{
				Name:       "projects/test-project/secrets/my-secret",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Labels: map[string]string{
					"env":  "test",
					"team": "dev",
				},
				Etag: "abc123",
			},
			wantErr: false,
		},
		{
			name: "Secret with annotations and version aliases",
			jsonData: `{
				"name": "projects/test-project/secrets/my-secret",
				"createTime": "2023-01-01T12:00:00Z",
				"annotations": {
					"owner": "team-alpha",
					"purpose": "database-password"
				},
				"versionAliases": {
					"current": "3",
					"previous": "2"
				},
				"etag": "def456",
				"replication": {
					"userManaged": {
						"replicas": [
							{"location": "us-central1"},
							{"location": "us-east1"}
						]
					}
				}
			}`,
			expected: SecretInfo{
				Name:       "projects/test-project/secrets/my-secret",
				CreateTime: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				Annotations: map[string]string{
					"owner":   "team-alpha",
					"purpose": "database-password",
				},
				VersionAliases: map[string]string{
					"current":  "3",
					"previous": "2",
				},
				Etag: "def456",
			},
			wantErr: false,
		},
		{
			name: "Minimal secret info",
			jsonData: `{
				"name": "projects/test-project/secrets/minimal-secret",
				"createTime": "2023-06-15T08:30:00Z",
				"etag": "xyz789",
				"replication": {
					"automatic": {}
				}
			}`,
			expected: SecretInfo{
				Name:       "projects/test-project/secrets/minimal-secret",
				CreateTime: time.Date(2023, 6, 15, 8, 30, 0, 0, time.UTC),
				Etag:       "xyz789",
			},
			wantErr: false,
		},
		{
			name:     "Invalid JSON",
			jsonData: `{"invalid": json}`,
			wantErr:  true,
		},
		{
			name:     "Empty JSON object",
			jsonData: `{}`,
			expected: SecretInfo{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var secretInfo SecretInfo
			err := json.Unmarshal([]byte(tt.jsonData), &secretInfo)
			
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

			// Compare basic fields
			if secretInfo.Name != tt.expected.Name {
				t.Errorf("Name = %q, expected %q", secretInfo.Name, tt.expected.Name)
			}

			if !secretInfo.CreateTime.Equal(tt.expected.CreateTime) {
				t.Errorf("CreateTime = %v, expected %v", secretInfo.CreateTime, tt.expected.CreateTime)
			}

			if secretInfo.Etag != tt.expected.Etag {
				t.Errorf("Etag = %q, expected %q", secretInfo.Etag, tt.expected.Etag)
			}

			// Compare labels
			if len(secretInfo.Labels) != len(tt.expected.Labels) {
				t.Errorf("Labels length = %d, expected %d", len(secretInfo.Labels), len(tt.expected.Labels))
			} else {
				for key, expectedValue := range tt.expected.Labels {
					if actualValue, exists := secretInfo.Labels[key]; !exists || actualValue != expectedValue {
						t.Errorf("Labels[%q] = %q, expected %q", key, actualValue, expectedValue)
					}
				}
			}

			// Compare annotations
			if len(secretInfo.Annotations) != len(tt.expected.Annotations) {
				t.Errorf("Annotations length = %d, expected %d", len(secretInfo.Annotations), len(tt.expected.Annotations))
			} else {
				for key, expectedValue := range tt.expected.Annotations {
					if actualValue, exists := secretInfo.Annotations[key]; !exists || actualValue != expectedValue {
						t.Errorf("Annotations[%q] = %q, expected %q", key, actualValue, expectedValue)
					}
				}
			}

			// Compare version aliases
			if len(secretInfo.VersionAliases) != len(tt.expected.VersionAliases) {
				t.Errorf("VersionAliases length = %d, expected %d", len(secretInfo.VersionAliases), len(tt.expected.VersionAliases))
			} else {
				for key, expectedValue := range tt.expected.VersionAliases {
					if actualValue, exists := secretInfo.VersionAliases[key]; !exists || actualValue != expectedValue {
						t.Errorf("VersionAliases[%q] = %q, expected %q", key, actualValue, expectedValue)
					}
				}
			}
		})
	}
}

// TestSecretVersionInfoJSONParsing tests SecretVersionInfo JSON parsing
func TestSecretVersionInfoJSONParsing(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected SecretVersionInfo
		wantErr  bool
	}{
		{
			name: "Complete version info",
			jsonData: `{
				"name": "projects/test-project/secrets/my-secret/versions/1",
				"createTime": "2023-01-01T12:00:00Z",
				"destroyTime": "2023-12-31T23:59:59Z",
				"state": "ENABLED",
				"etag": "version-abc123"
			}`,
			expected: SecretVersionInfo{
				Name:        "projects/test-project/secrets/my-secret/versions/1",
				CreateTime:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
				DestroyTime: time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
				State:       "ENABLED",
				Etag:        "version-abc123",
			},
			wantErr: false,
		},
		{
			name: "Version info without destroy time",
			jsonData: `{
				"name": "projects/test-project/secrets/my-secret/versions/2",
				"createTime": "2023-06-15T14:30:00Z",
				"state": "ENABLED",
				"etag": "version-def456"
			}`,
			expected: SecretVersionInfo{
				Name:       "projects/test-project/secrets/my-secret/versions/2",
				CreateTime: time.Date(2023, 6, 15, 14, 30, 0, 0, time.UTC),
				State:      "ENABLED",
				Etag:       "version-def456",
			},
			wantErr: false,
		},
		{
			name:     "Invalid JSON",
			jsonData: `{"invalid": json}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var versionInfo SecretVersionInfo
			err := json.Unmarshal([]byte(tt.jsonData), &versionInfo)
			
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

			if versionInfo.Name != tt.expected.Name {
				t.Errorf("Name = %q, expected %q", versionInfo.Name, tt.expected.Name)
			}

			if !versionInfo.CreateTime.Equal(tt.expected.CreateTime) {
				t.Errorf("CreateTime = %v, expected %v", versionInfo.CreateTime, tt.expected.CreateTime)
			}

			if !versionInfo.DestroyTime.Equal(tt.expected.DestroyTime) {
				t.Errorf("DestroyTime = %v, expected %v", versionInfo.DestroyTime, tt.expected.DestroyTime)
			}

			if versionInfo.State != tt.expected.State {
				t.Errorf("State = %q, expected %q", versionInfo.State, tt.expected.State)
			}

			if versionInfo.Etag != tt.expected.Etag {
				t.Errorf("Etag = %q, expected %q", versionInfo.Etag, tt.expected.Etag)
			}
		})
	}
}

// TestGetSecretInput tests secret input handling from various sources
func TestGetSecretInput(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		dataFile string
		prompt   string
		expected string
		wantErr  bool
	}{
		{
			name:     "Direct data input",
			data:     "my-secret-value",
			dataFile: "",
			prompt:   "Enter secret: ",
			expected: "my-secret-value",
			wantErr:  false,
		},
		{
			name:     "Empty data should return empty",
			data:     "",
			dataFile: "",
			prompt:   "Enter secret: ",
			expected: "", // Will fail in interactive mode, but we test the logic
			wantErr:  true, // Will error in test environment due to no terminal
		},
		{
			name:     "Data takes precedence over file",
			data:     "direct-data",
			dataFile: "/some/file",
			prompt:   "Enter secret: ",
			expected: "direct-data",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip interactive tests as they can't work in automated testing
			if tt.data == "" && tt.dataFile == "" {
				t.Skip("Skipping interactive test")
				return
			}

			result, err := getSecretInput(tt.data, tt.dataFile, tt.prompt)
			
			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("getSecretInput() = %q, expected %q", result, tt.expected)
			}
		})
	}
}