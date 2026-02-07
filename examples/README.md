# Examples

This directory contains example files demonstrating `gsecutil` configuration and usage.

## Configuration Files

Configuration files define default settings, secret metadata, and custom attributes.

### `gsecutil.conf`
Complete configuration example with all available options. Includes:
- Default GCP project settings
- Secret titles and descriptions
- Custom attributes (owner, cost_center, etc.)
- Label definitions

Use as reference for full configuration syntax.

### `gsecutil-minimal.conf`
Minimal configuration with only essential settings. Good starting point for new users.

### `gsecutil-advanced.conf`
Advanced configuration demonstrating complex metadata structures and label combinations.

### `gsecutil-windows.conf`
Configuration example with Windows-specific file paths and settings.

## CSV Sample Files

The `csv/` directory contains sample CSV files for import/export operations.

### `csv/sample-import.csv`
Complete example with all columns including secret values. Demonstrates:
- Basic secret creation
- Labels (`label:env`, `label:team`)
- Custom attributes (owner, description)
- Title metadata

Use with:
```bash
gsecutil import examples/csv/sample-import.csv --dry-run
```

### `csv/sample-multiline.csv`
Example of multi-line values (Excel format) including:
- SSL certificates with line breaks
- JSON configuration files
- Multi-line text values

Use with:
```bash
gsecutil import examples/csv/sample-multiline.csv --dry-run
```

### `csv/sample-metadata-only.csv`
Metadata-only CSV without secret values. Use with `--update-config` to update configuration file:
```bash
gsecutil import examples/csv/sample-metadata-only.csv --update-config
```

## Output Examples

### `describe-output-example.md`
Example output from the `describe` command showing secret metadata, version information, and replication settings.

### `list-output-example.md`
Example output from the `list` command demonstrating various formatting options and filtering capabilities.

## CSV Format Reference

CSV files support the following columns:

### Required
- `name` - Secret name (required)
- `value` - Secret value (required for creation, optional for metadata-only updates)

### Optional
- `title` - Secret title (stored in config)
- `label:<key>` - Labels (e.g., `label:env`, `label:team`)
- Any other columns - Custom attributes stored in config

### Example

```csv
name,value,title,owner,label:env,label:team
my-secret,secret123,My Secret,alice,production,backend
```

## Common Usage Examples

### Import Operations

Import with dry-run (preview changes without applying):
```bash
gsecutil import examples/csv/sample-import.csv --dry-run
```

Import and create secrets:
```bash
gsecutil import examples/csv/sample-import.csv
```

Import and update configuration:
```bash
gsecutil import examples/csv/sample-import.csv --update-config
```

Import with upsert (create or update):
```bash
gsecutil import examples/csv/sample-import.csv --upsert
```

### Export Operations

Export to CSV (metadata only):
```bash
gsecutil export -o output.csv
```

Export with secret values:
```bash
gsecutil export --with-values -o backup.csv
```

Export filtered secrets:
```bash
gsecutil export --filter env=production -o prod-secrets.csv
```

### Configuration

Specify custom configuration file:
```bash
gsecutil --config examples/gsecutil.conf list
```

Use minimal configuration:
```bash
gsecutil --config examples/gsecutil-minimal.conf describe my-secret
```

## Additional Resources

For detailed documentation on CSV operations, see:
- [CSV Operations Guide](../docs/csv-operations.md) - Comprehensive guide to import/export functionality
- [Commands Reference](../docs/commands.md) - Complete command documentation
- [Configuration Guide](../docs/configuration.md) - Detailed configuration file documentation
