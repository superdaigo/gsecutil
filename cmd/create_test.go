package cmd

import (
	"reflect"
	"sort"
	"testing"
)

func TestMergeLabelsWithDefaults(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		userLabels  []string
		wantLabels  []string
		description string
	}{
		{
			name: "no default labels",
			config: &Config{
				Defaults: DefaultConfig{
					Labels: map[string]string{},
				},
			},
			userLabels:  []string{"custom=value"},
			wantLabels:  []string{"custom=value"},
			description: "Should return user labels when no defaults configured",
		},
		{
			name: "only default labels",
			config: &Config{
				Defaults: DefaultConfig{
					Labels: map[string]string{
						"managed_by": "gsecutil",
						"team":       "platform",
					},
				},
			},
			userLabels:  []string{},
			wantLabels:  []string{"managed_by=gsecutil", "team=platform"},
			description: "Should apply default labels when no user labels provided",
		},
		{
			name: "merge default and user labels",
			config: &Config{
				Defaults: DefaultConfig{
					Labels: map[string]string{
						"managed_by": "gsecutil",
						"team":       "platform",
					},
				},
			},
			userLabels:  []string{"environment=production"},
			wantLabels:  []string{"managed_by=gsecutil", "team=platform", "environment=production"},
			description: "Should merge default and user labels",
		},
		{
			name: "user labels override defaults",
			config: &Config{
				Defaults: DefaultConfig{
					Labels: map[string]string{
						"managed_by":  "gsecutil",
						"team":        "platform",
						"environment": "dev",
					},
				},
			},
			userLabels:  []string{"environment=production", "version=1.0"},
			wantLabels:  []string{"managed_by=gsecutil", "team=platform", "environment=production", "version=1.0"},
			description: "User-provided labels should override default labels",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global config for test
			globalConfig = tt.config

			// Call function
			got := mergeLabelsWithDefaults(tt.userLabels)

			// Sort both slices for comparison (map iteration order is not guaranteed)
			sort.Strings(got)
			sort.Strings(tt.wantLabels)

			if !reflect.DeepEqual(got, tt.wantLabels) {
				t.Errorf("mergeLabelsWithDefaults() = %v, want %v\nDescription: %s", got, tt.wantLabels, tt.description)
			}
		})
	}

	// Reset global config
	globalConfig = nil
}
