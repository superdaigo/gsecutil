package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestRootCommand tests the basic root command functionality
func TestRootCommand(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new root command for testing (to avoid modifying global state)
	rootCmd := &cobra.Command{
		Use:   "gsecutil",
		Short: "A Google Secret Manager utility CLI",
		Long: `gsecutil is a command-line utility that provides a simple wrapper
around the gcloud CLI for managing Google Secret Manager secrets.

It allows you to get, create, update, delete, list, and describe secrets
with simplified commands, and also provides the ability to copy secret
values directly to your clipboard.`,
	}

	// Set the output to our buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Test help output
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("Root command help should not return an error, got: %v", err)
	}

	output := buf.String()

	// Check that the help output contains expected content
	expectedStrings := []string{
		"gsecutil",
		"command-line utility",
		"gcloud CLI",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected help output to contain %q, but it didn't. Output: %s", expected, output)
		}
	}
}

// TestRootCommandFlags tests that the root command accepts the expected flags
func TestRootCommandFlags(t *testing.T) {
	// Create a new root command for testing
	rootCmd := &cobra.Command{
		Use:   "gsecutil",
		Short: "A Google Secret Manager utility CLI",
	}

	// Add the project flag (as done in the real root.go)
	rootCmd.PersistentFlags().StringP("project", "p", "", "Google Cloud project ID")

	// Test that the project flag exists
	projectFlag := rootCmd.PersistentFlags().Lookup("project")
	if projectFlag == nil {
		t.Error("Expected --project flag to exist")
	}

	if projectFlag.Shorthand != "p" {
		t.Errorf("Expected --project flag shorthand to be 'p', got %q", projectFlag.Shorthand)
	}

	if projectFlag.Usage != "Google Cloud project ID" {
		t.Errorf("Expected --project flag usage to be 'Google Cloud project ID', got %q", projectFlag.Usage)
	}
}

// TestRootCommandVersion tests the version flag functionality
func TestRootCommandVersion(t *testing.T) {
	// Create a buffer to capture output
	var buf bytes.Buffer

	// Create a new root command for testing
	rootCmd := &cobra.Command{
		Use:   "gsecutil",
		Short: "A Google Secret Manager utility CLI",
	}

	// Set version (simulate what's done in the real application)
	rootCmd.Version = "0.2.0"

	// Set the output to our buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	// Test version output
	rootCmd.SetArgs([]string{"--version"})
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("Root command version should not return an error, got: %v", err)
	}

	output := buf.String()

	// Check that version output contains the expected version
	if !strings.Contains(output, "0.2.0") {
		t.Errorf("Expected version output to contain '0.2.0', but got: %s", output)
	}
}

// TestCommandStructure tests that the root command has the expected structure
func TestCommandStructure(t *testing.T) {
	// Test that our actual root command has the expected properties
	if rootCmd.Use != "gsecutil" {
		t.Errorf("Expected root command Use to be 'gsecutil', got %q", rootCmd.Use)
	}

	if !strings.Contains(rootCmd.Short, "Google Secret Manager utility CLI") {
		t.Errorf("Expected root command Short to contain 'Google Secret Manager utility CLI', got %q", rootCmd.Short)
	}

	if !strings.Contains(rootCmd.Long, "command-line utility") {
		t.Errorf("Expected root command Long to contain 'command-line utility', got %q", rootCmd.Long)
	}

	// Test that the project flag is set up correctly
	projectFlag := rootCmd.PersistentFlags().Lookup("project")
	if projectFlag == nil {
		t.Error("Expected --project flag to be defined in root command")
	}

	// Just test that there are some commands (exact count depends on init order)
	if len(rootCmd.Commands()) == 0 {
		t.Error("Expected root command to have subcommands, but found none")
	}
}

// TestGlobalProjectFlag tests that the global project flag works correctly
func TestGlobalProjectFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Short flag",
			args:     []string{"-p", "test-project"},
			expected: "test-project",
		},
		{
			name:     "Long flag",
			args:     []string{"--project", "my-gcp-project"},
			expected: "my-gcp-project",
		},
		{
			name:     "No flag",
			args:     []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test command that checks the project flag
			var capturedProject string
			testCmd := &cobra.Command{
				Use: "test",
				Run: func(cmd *cobra.Command, args []string) {
					capturedProject, _ = cmd.Flags().GetString("project")
				},
			}

			// Create root command and add project flag
			rootCmd := &cobra.Command{Use: "root"}
			rootCmd.PersistentFlags().StringP("project", "p", "", "Google Cloud project ID")
			rootCmd.AddCommand(testCmd)

			// Execute with test args
			rootCmd.SetArgs(append(tt.args, "test"))
			err := rootCmd.Execute()

			if err != nil {
				t.Errorf("Command execution failed: %v", err)
			}

			if capturedProject != tt.expected {
				t.Errorf("Expected project flag value %q, got %q", tt.expected, capturedProject)
			}
		})
	}
}
