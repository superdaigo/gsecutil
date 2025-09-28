# gsecutil v1.0.0 Release Summary

## 🎉 Major Release - Complete Access Management Suite

**Release Date**: 2025-09-28
**Version**: 1.0.0
**Previous Version**: 0.3.0

---

## 🚀 What's New in v1.0.0

### **Complete IAM Access Management System**

This release introduces a comprehensive access management suite that makes gsecutil a complete solution for Google Secret Manager administration:

#### **New `access` Command Suite**
- `gsecutil access list <secret>` - List all principals with access to a secret
- `gsecutil access grant <secret> --principal <user/group/serviceAccount>` - Grant access
- `gsecutil access revoke <secret> --principal <user/group/serviceAccount>` - Revoke access
- `gsecutil access project` - Show project-level Secret Manager permissions

#### **Enhanced `list` Command**
- `gsecutil list --principal <user/group/serviceAccount>` - Show secrets accessible by a principal

#### **Parameter Consistency**
- Unified `--principal` parameter across all commands (auditlog, access, list)
- Consistent `--operation` parameter naming

---

## 🔧 Key Features

### **Multi-Level Permission Analysis**
- **Secret-level IAM policies**: Direct secret permissions
- **Project-level IAM policies**: Editor/Owner roles that provide Secret Manager access
- **Combined analysis**: Automatically checks both levels

### **IAM Condition Awareness**
- **Full CEL expression support**: Shows conditional access policies
- **Condition details**: Displays titles, descriptions, and expressions
- **Time-based conditions**: Understands access restrictions by time/date
- **Resource-based conditions**: Shows resource path limitations

### **Comprehensive Principal Support**
- **Users**: `user:alice@example.com`
- **Groups**: `group:team@example.com`
- **Service Accounts**: `serviceAccount:app@project.iam.gserviceaccount.com`
- **Domains**: `domain:example.com`
- **Special principals**: `allUsers`, `allAuthenticatedUsers`

---

## 📊 Usage Examples

### Access Management Workflow
```bash
# See who has access to a secret (both levels)
gsecutil access list my-secret --include-project

# Grant access to a user
gsecutil access grant my-secret --principal user:alice@example.com

# Grant viewer access to a service account
gsecutil access grant my-secret \
  --principal serviceAccount:app@project.iam.gserviceaccount.com \
  --role roles/secretmanager.viewer

# Revoke access
gsecutil access revoke my-secret --principal user:bob@example.com

# See all project-level permissions
gsecutil access project
```

### Principal-Based Discovery
```bash
# See what secrets a user can access
gsecutil list --principal user:alice@example.com

# Audit what a service account accessed
gsecutil auditlog --principal app@project.iam.gserviceaccount.com
```

---

## 🛠️ Technical Improvements

### **Enhanced IAM Policy Parsing**
- Complete IAM policy JSON parsing including conditions
- Project-level IAM policy queries
- Cross-platform gcloud integration

### **Robust Error Handling**
- Clear validation for principal formats
- Graceful handling of missing policies
- Network resilience for gcloud operations

### **Consistent API Design**
- Unified parameter naming across all commands
- Predictable command structure and patterns
- Comprehensive help text and examples

---

## ⚠️ Breaking Changes

### **Parameter Name Changes**
- `auditlog --user` → `auditlog --principal`
- `auditlog --operations` → `auditlog --operation`

**Migration is simple:**
```bash
# Old (still works but deprecated)
gsecutil auditlog --user alice --operations ACCESS

# New (recommended)
gsecutil auditlog --principal alice --operation ACCESS
```

---

## 📦 Installation

### Pre-built Binaries
Download from GitHub releases:
- Linux (x64/ARM64): `gsecutil-linux-amd64-v1.0.0`, `gsecutil-linux-arm64-v1.0.0`
- macOS (Intel/Apple Silicon): `gsecutil-darwin-amd64-v1.0.0`, `gsecutil-darwin-arm64-v1.0.0`
- Windows (x64): `gsecutil-windows-amd64-v1.0.0.exe`

### Go Install
```bash
go install github.com/superdaigo/gsecutil@v1.0.0
```

### Build from Source
```bash
git clone https://github.com/superdaigo/gsecutil.git
cd gsecutil
make build-all
```

---

## ✅ Quality Assurance

### **Testing**
- ✅ All existing tests pass
- ✅ New access management functions tested
- ✅ Parameter consistency validated
- ✅ Cross-platform builds successful

### **Code Quality**
- ✅ `go fmt` - Code formatting
- ✅ `go vet` - Static analysis
- ✅ Comprehensive error handling
- ✅ Consistent patterns throughout

### **Documentation**
- ✅ Updated README with v1.0.0 features
- ✅ Comprehensive CHANGELOG entry
- ✅ Help text updated for all commands
- ✅ Usage examples for all new features

---

## 🎯 Why v1.0.0?

This release represents a **complete, production-ready** Secret Manager administration tool:

✅ **Complete Feature Set** - All essential Secret Manager operations
✅ **Enterprise-Ready** - IAM conditions and complex access policies
✅ **Consistent API** - Unified parameters and predictable patterns
✅ **Production Stability** - Robust error handling and edge cases
✅ **Comprehensive Documentation** - Complete usage examples and guides

---

## 🔜 Next Steps

1. **Commit and tag the release**:
   ```bash
   git add .
   git commit -m "Release v1.0.0: Complete Access Management Suite"
   git tag -a v1.0.0 -m "v1.0.0: Major release with complete IAM access management"
   git push origin main --tags
   ```

2. **Create GitHub release** with release notes and binaries

3. **Update documentation** and examples

---

**🎉 Congratulations on reaching v1.0.0!**

This release transforms gsecutil from a simple secret management tool into a comprehensive Secret Manager administration platform with enterprise-grade access management capabilities.
