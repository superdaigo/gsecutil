# Testing GitHub Actions Workflows

This document provides instructions for testing the GitHub Actions workflows in gsecutil.

## Workflows Overview

1. **CI Workflow** (`ci.yml`) - Runs on push/PR to main/develop branches
2. **Release Workflow** (`release.yml`) - Runs when version tags are pushed
3. **Draft Release Workflow** (`draft-release.yml`) - Runs on main branch pushes or manual trigger

## Testing Each Workflow

### 1. Testing CI Workflow ✅

**Triggers:** Push to main/develop branches or pull requests

**Status:** Should have already triggered when you pushed the workflows to main.

**Check status:**
- Go to: https://github.com/superdaigo/gsecutil/actions
- Look for the "CI" workflow run

**What it tests:**
- Go formatting (`go fmt`)
- Static analysis (`go vet`)
- Linting (`golangci-lint`)
- Tests (`go test`)
- Cross-platform builds
- Security scanning
- Dependency vulnerability checks

### 2. Testing Draft Release Workflow

**Option A: Manual Trigger**
1. Go to https://github.com/superdaigo/gsecutil/actions
2. Click on "Draft Release" workflow
3. Click "Run workflow" button
4. Click "Run workflow" to confirm

**Option B: Push to Main**
The workflow will automatically trigger on the next push to main (unless the commit message contains `[skip draft]`).

**What it does:**
- Auto-increments version number
- Builds all platform binaries
- Creates a draft release with binaries
- Provides testing checklist
- Awaits manual review and publishing

### 3. Testing Release Workflow (Full Release)

**Using the Release Helper Script:**

```bash
# Test with a small version number first
./scripts/release.sh 0.1.0
```

**Or manually:**

```bash
# Create and push a version tag
git tag v0.1.0
git push origin v0.1.0
```

**What it does:**
- Builds all platform binaries
- Creates release notes with changelog
- Uploads binaries with checksums
- Publishes the release automatically
- Provides download links and installation instructions

## Test Plan

### Phase 1: Verify CI Workflow ✅
- [x] Workflow triggered on main branch push
- [ ] All jobs complete successfully (test, lint, build, security)
- [ ] Artifacts are created and stored

### Phase 2: Test Draft Release
```bash
# Trigger manually or wait for next push to main
# Check that draft release is created with binaries
```

### Phase 3: Test Full Release
```bash
# Use release script for controlled testing
./scripts/release.sh 0.1.0

# Verify:
# - Release is created
# - All binaries are uploaded
# - Checksums are generated
# - Download links work
# - Installation instructions are correct
```

### Phase 4: Test Release Assets
```bash
# Download and test each binary
curl -LO https://github.com/superdaigo/gsecutil/releases/download/v0.1.0/gsecutil-linux-amd64
chmod +x gsecutil-linux-amd64
./gsecutil-linux-amd64 --version
./gsecutil-linux-amd64 --help

# Verify checksum
curl -LO https://github.com/superdaigo/gsecutil/releases/download/v0.1.0/gsecutil-linux-amd64.sha256
sha256sum -c gsecutil-linux-amd64.sha256
```

## Troubleshooting

### Common Issues

1. **Workflow not triggering:**
   - Check branch names (main vs master)
   - Ensure workflows are in `.github/workflows/`
   - Check GitHub Actions permissions

2. **Build failures:**
   - Check Go version compatibility
   - Verify dependencies can be downloaded
   - Check for linting/formatting issues

3. **Release creation fails:**
   - Verify repository has releases enabled
   - Check tag format (must start with 'v')
   - Ensure no existing tag with same name

4. **Binary upload fails:**
   - Check file size limits
   - Verify network connectivity
   - Check GitHub API rate limits

### Checking Workflow Status

```bash
# Using GitHub CLI (if installed)
gh run list

# Or visit the Actions page
echo "https://github.com/superdaigo/gsecutil/actions"
```

### Manual Workflow Trigger

You can manually trigger workflows from the GitHub Actions UI:
1. Go to Actions tab
2. Select the workflow
3. Click "Run workflow"
4. Choose branch and click "Run workflow"

## Test Results Template

```markdown
## CI Workflow Test Results
- [ ] Test job completed successfully  
- [ ] Lint job completed successfully
- [ ] Build job completed successfully
- [ ] Security scan completed successfully
- [ ] Dependency check completed successfully
- [ ] Artifacts uploaded

## Draft Release Test Results
- [ ] Draft release created
- [ ] All platform binaries included
- [ ] Checksums generated
- [ ] Release notes generated
- [ ] Testing checklist included

## Release Workflow Test Results  
- [ ] Release created on tag push
- [ ] All binaries uploaded correctly
- [ ] Checksums match
- [ ] Download links work
- [ ] Installation instructions accurate
- [ ] Version number correct in binaries
```

## Next Steps After Testing

1. **If tests pass:** The workflows are ready for production use
2. **If tests fail:** Fix issues and re-test
3. **Create production release:** Use `./scripts/release.sh 1.0.0` for first stable release

## Monitoring

- **GitHub Actions:** https://github.com/superdaigo/gsecutil/actions
- **Releases:** https://github.com/superdaigo/gsecutil/releases
- **Repository:** https://github.com/superdaigo/gsecutil