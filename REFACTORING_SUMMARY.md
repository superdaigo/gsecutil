# List Command Refactoring Summary

## Overview
Refactored the `gsecutil list` command following **Option 4 (Hybrid Approach)** to fix inconsistencies and improve clarity while maintaining functionality.

## Changes Made

### 1. Flag Name Changes

#### Removed Flags
- `--no-labels` → Replaced with `--show-labels` (labels now **hidden by default**)
- `--filter-attributes` → Replaced with clearer `--attr-filter`

#### New/Updated Flags
- `--show-labels` - Show labels in output (labels are **hidden by default**)
- `--attr-filter` - Filter by configuration file attributes
- `--show` - Primary flag for showing custom attributes
- `--show-attributes` - Hidden alias for backward compatibility

### 2. Column Header Standardization

**Before:**
- `CREATE_TIME (UTC)` in simple mode
- `CREATED (UTC)` in labels mode
- `UPDATE_TIME (UTC)` when showing updated times

**After (Consistent):**
- `CREATED (UTC)` - Always
- `UPDATED (UTC)` - When using `--show-updated`

### 3. Updated Files

#### Code Files
- `cmd/list.go` - Updated flag definitions, function signatures, column headers
- `cmd/list_test.go` - Updated test expectations for new column headers

#### Documentation Files
- `docs/commands.md` - Updated flag names and examples
- `docs/configuration.md` - Updated all references to old flags
- `examples/list-output-example.md` - Updated all command examples
- `examples/describe-output-example.md` - Updated usage pattern examples

### 4. Backward Compatibility

The `--show-attributes` flag is kept as a hidden alias:
```go
listCmd.Flags().String("show-attributes", "", "(Alias for --show) ...")
listCmd.Flags().MarkHidden("show-attributes")
```

This allows existing scripts to continue working during the transition period.

## Migration Guide

### For Users

**Old Command** → **New Command**

```bash
# Show labels (they're hidden by default now)
gsecutil list --no-labels
→ gsecutil list  # Labels hidden by default
→ gsecutil list --show-labels  # To show labels

# Filter by attributes
gsecutil list --filter-attributes "environment=prod"
→ gsecutil list --attr-filter "environment=prod"

# Show custom attributes
gsecutil list --show-attributes "title,owner"
→ gsecutil list --show "title,owner"
# (--show-attributes still works but is hidden)
```

### Benefits

1. **Cleaner Default Output**: Labels are hidden by default, making list output more concise
2. **Ergonomic Flags**: Use `--show-labels` to show labels (no need for `=false` syntax)
3. **Clearer Intent**: `--attr-filter` is more explicit than `--filter-attributes`
4. **Shorter Syntax**: `--show` instead of `--show-attributes`
5. **Standardized Output**: Column headers are consistent across all display modes

## Testing

All changes have been tested:
- ✅ Code builds successfully
- ✅ All unit tests pass
- ✅ `make fmt` and `make vet` pass
- ✅ Help output displays correctly with new flags
- ✅ Backward compatibility maintained via hidden alias

## Future Considerations

For v2.0, consider:
- Removing the hidden `--show-attributes` alias
- Adding `--columns` flag for more flexible column selection
- Potentially moving `--principal` to a separate `access` subcommand
