# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2025-09-20

### 🧪 Testing & Quality Assurance

#### Comprehensive Test Suite
- **Added 67+ test cases** covering core business logic and edge cases
- **Test coverage: 21.4%** focused on high-value code paths including:
  - Data parsing and validation logic (JSON unmarshaling, operations filtering)
  - Business rules and filtering algorithms (audit log filtering, secret list processing)
  - Helper functions (version extraction, label formatting, column width calculation)
  - Command structure and flag handling (root command functionality)
- **Four test files created**:
  - `cmd/auditlog_test.go` - Tests for audit log operations, filtering, and validation
  - `cmd/clipboard_test.go` - Tests for clipboard utilities, JSON parsing, and helper functions
  - `cmd/list_test.go` - Tests for secret list parsing, formatting, and display logic
  - `cmd/root_test.go` - Tests for root command functionality and flag handling

#### Development Quality Tools
- **Pre-commit hooks** configured with `.pre-commit-config.yaml`:
  - **File quality checks**: Remove trailing whitespace, fix line endings, check YAML/JSON syntax
  - **Go-specific checks**: `go fmt`, `go vet`, `go mod tidy`, `go test ./cmd`
  - **Advanced linting**: `golangci-lint` with staticcheck analysis
- **Automatic quality enforcement** on every commit
- **CI/CD compatibility** - hooks match GitHub Actions workflow requirements

### 🔧 Code Quality Improvements

#### Static Analysis Fixes
- **Fixed SA5011 staticcheck errors**: Resolved nil pointer dereference issues in test files
- **Enhanced error handling**: Added proper early returns after nil checks
- **Golangci-lint integration**: Now runs clean with `0 issues`
- **Go 1.24 compatibility**: Upgraded golangci-lint to latest version (2.4.0)

#### Documentation Enhancements
- **Updated README.md** with pre-commit setup instructions
- **Enhanced Contributing section** with detailed development workflow
- **Automated quality checks documentation** for contributors
- **Clear setup instructions** for development environment

### 🚀 Developer Experience

#### Automated Code Quality
- **Consistent formatting**: Automatic `go fmt` on every commit
- **Dependency management**: Automatic `go mod tidy`
- **Test execution**: Automatic test run before commits
- **Early problem detection**: Issues caught before code review

#### Improved Maintainability
- **Comprehensive test coverage** of critical business logic
- **Consistent code style** enforced automatically
- **Reduced manual testing** burden through automated tests
- **Better code documentation** through test examples

### 📋 Technical Details

#### Test Framework
- **Table-driven tests** for comprehensive scenario coverage
- **Helper functions** to reduce code duplication and improve readability
- **Proper test isolation** preventing test interdependencies
- **Edge case coverage** including error conditions and boundary values

#### Pre-commit Configuration
- **Multi-stage validation**: File checks → Go checks → Tests → Linting
- **Fast execution**: Only runs on changed files
- **Auto-fixing**: Some issues (like formatting) are fixed automatically
- **Blocking commits**: Prevents commits with quality issues

### 🎯 Impact

- **Improved code reliability** through comprehensive testing
- **Faster development cycles** with automated quality checks
- **Easier contributions** with clear development standards
- **Reduced bug risk** through static analysis and testing
- **Better maintainability** with consistent code quality

## [0.2.0] - 2025-09-20

### 🚀 Major Features Added

#### Enhanced `describe` Command
- **Comprehensive secret metadata display** including:
  - Default version information (version number, state, creation time)
  - Replication strategy (automatic multi-region or user-managed)
  - Labels and annotations (tags) with alphabetical sorting
  - Version aliases, expiration settings, and rotation configuration
  - Pub/Sub topic integrations
- **Improved data structure** with enhanced `SecretInfo` struct
- **Better help text** with detailed feature descriptions
- **Backward compatibility** maintained for JSON passthrough and existing flags

#### Enhanced `auditlog` Command (formerly `audit`)
- **Renamed command** from `audit` to `auditlog` for better clarity
- **Operations filtering** with `--operations` flag supporting:
  - Single operation: `--operations ACCESS`
  - Multiple operations: `--operations ACCESS,CREATE,DELETE`
  - Case-insensitive matching
  - Complete list of supported operations: ACCESS, CREATE, UPDATE, DELETE, GET_METADATA, LIST, UPDATE_METADATA, DESTROY_VERSION, DISABLE_VERSION, ENABLE_VERSION
- **Enhanced filtering system**:
  - Optional secret name parameter (show all logs when not specified)
  - User filtering with partial matching (`--user`)
  - Combined filtering (secret + user + operations)
  - Improved secret-related operation detection (filters out location listing noise)
- **Better user experience**:
  - Dynamic table headers showing active filters
  - Context-aware "no results" messages
  - Enhanced help text with comprehensive examples

#### Enhanced `list` Command
- **Labels support** with automatic display of secret labels
- **Flexible display options**:
  - `--no-labels` flag to hide labels column
  - Enhanced table formatting with dynamic column widths
  - Alphabetical sorting of labels for consistency
- **Format compatibility** maintained with existing gcloud passthrough
- **Clean presentation** showing "-" for secrets without labels

### 📖 Documentation Improvements

#### Comprehensive Audit Logging Documentation
- **New documentation file** `docs/audit-logging.md` with detailed setup instructions
- **Multiple setup methods** covered:
  - Google Cloud Console (beginner-friendly)
  - gcloud CLI commands
  - Terraform configuration
- **Cost considerations** and optimization tips
- **Troubleshooting guide** for common issues
- **Security best practices** for audit log management

#### Enhanced README and Help Text
- **Updated examples** throughout all commands
- **Better feature descriptions** with visual indicators
- **Comprehensive usage examples** for new filtering capabilities
- **Clear documentation** of version support in `get` command (already existed)

### 🔧 Technical Improvements
- **Improved error handling** across all enhanced commands
- **Better JSON parsing** with comprehensive struct definitions
- **Enhanced filtering algorithms** with partial matching support
- **Cleaner code organization** with helper functions
- **Maintained backward compatibility** for all existing functionality

### 📋 Version Support Clarification
- **`get` command version support** was already fully implemented but now better documented
- **Enhanced help text** and examples for version-specific secret retrieval
- **Complete backward compatibility** with existing version functionality

## [0.1.8] - 2025-09-04

### Fixed
- Fixed GitHub Actions CI workflow issues
- Updated deprecated actions and improved security scanning
- Enhanced release workflow with proper changelog generation

## Previous Versions

See Git history for changes in versions 0.1.7 and earlier.

---

## Legend

- 🚀 Major Features
- ✨ New Features
- 🔧 Improvements
- 🐛 Bug Fixes
- 📖 Documentation
- 🔒 Security
- ⚠️ Breaking Changes
- 🗑️ Deprecated
