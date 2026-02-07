package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [output-file]",
	Short: "Export secrets to CSV file",
	Long: `Export secrets and their metadata to a CSV file.

The CSV includes secret names, values, labels, and metadata from the configuration file.
If no output file is specified, output is written to stdout.

The exported CSV can be edited in Excel or other spreadsheet applications and
re-imported using the 'import' command.`,
	Example: `  gsecutil export secrets.csv
  gsecutil export secrets.csv --with-values
  gsecutil export > secrets.csv
  gsecutil export --filter "labels.env=prod" secrets.csv`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

var (
	exportWithValues bool
	exportFilter     string
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().BoolVar(&exportWithValues, "with-values", false, "Include secret values in export (use with caution)")
	exportCmd.Flags().StringVar(&exportFilter, "filter", "", "Filter secrets by label")
}

func runExport(cmd *cobra.Command, args []string) error {
	project, _ := cmd.Flags().GetString("project")
	project = GetProject(project)

	// Get list of secrets
	secrets, err := fetchSecretsForExport(project, exportFilter)
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets found to export")
		return nil
	}

	// Prepare CSV data
	records := prepareCsvRecords(secrets, exportWithValues, project)

	// Write to file or stdout
	var writer *csv.Writer
	if len(args) > 0 {
		file, err := os.Create(args[0])
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		writer = csv.NewWriter(file)
		defer writer.Flush()
	} else {
		writer = csv.NewWriter(os.Stdout)
		defer writer.Flush()
	}

	// Write CSV
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV: %w", err)
		}
	}

	if len(args) > 0 {
		fmt.Printf("Exported %d secrets to %s\n", len(secrets), args[0])
	}

	return nil
}

func fetchSecretsForExport(project, filter string) ([]SecretInfo, error) {
	secrets, err := fetchSecrets(project, filter, 0)
	if err != nil {
		return nil, err
	}

	// Sort by name
	sortSecrets(secrets)

	return secrets, nil
}

func prepareCsvRecords(secrets []SecretInfo, withValues bool, project string) [][]string {
	// Collect all unique label keys and config attributes
	labelKeys := make(map[string]bool)
	configAttrs := make(map[string]bool)

	for _, secret := range secrets {
		name := extractSecretName(secret.Name)

		// Collect label keys
		for key := range secret.Labels {
			labelKeys[key] = true
		}

		// Collect config attributes
		if credInfo := GetCredentialInfo(name); credInfo != nil {
			for key := range credInfo.Attributes {
				configAttrs[key] = true
			}
		}
	}

	// Sort keys for consistent column order
	labelKeysSorted := make([]string, 0, len(labelKeys))
	for key := range labelKeys {
		labelKeysSorted = append(labelKeysSorted, key)
	}
	sort.Strings(labelKeysSorted)

	configAttrsSorted := make([]string, 0, len(configAttrs))
	for key := range configAttrs {
		configAttrsSorted = append(configAttrsSorted, key)
	}
	sort.Strings(configAttrsSorted)

	// Build header
	header := []string{"name"}
	if withValues {
		header = append(header, "value")
	}
	header = append(header, "title")
	for _, key := range labelKeysSorted {
		header = append(header, "label:"+key)
	}
	for _, key := range configAttrsSorted {
		header = append(header, key)
	}

	records := [][]string{header}

	// Build data rows
	for _, secret := range secrets {
		name := extractSecretName(secret.Name)
		row := []string{name}

		// Add value if requested
		if withValues {
			value := getSecretValue(name, project)
			row = append(row, value)
		}

		// Add title from config
		credInfo := GetCredentialInfo(name)
		if credInfo != nil && credInfo.Title != "" {
			row = append(row, credInfo.Title)
		} else {
			row = append(row, "")
		}

		// Add labels
		for _, key := range labelKeysSorted {
			if value, exists := secret.Labels[key]; exists {
				row = append(row, value)
			} else {
				row = append(row, "")
			}
		}

		// Add config attributes
		for _, key := range configAttrsSorted {
			if credInfo != nil {
				if value, exists := credInfo.Attributes[key]; exists {
					row = append(row, fmt.Sprintf("%v", value))
				} else {
					row = append(row, "")
				}
			} else {
				row = append(row, "")
			}
		}

		records = append(records, row)
	}

	return records
}

