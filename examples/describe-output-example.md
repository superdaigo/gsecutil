# gsecutil describe Command with Configuration

This document shows how the `gsecutil describe` command output would look when using configuration files with metadata.

## Configuration File Example

```yaml
# gsecutil.conf
project: "team-secrets-prod"
prefix: "team-"

list:
  attributes:
    - title
    - owner
    - environment

credentials:
  - name: "team-db-prod"
    title: "Production Database"
    description: "PostgreSQL master database connection string"
    environment: "production"
    category: "database"
    owner: "backend-team"
    contact: "backend-team@company.com"
    rotation_schedule: "quarterly"
    sensitive_level: "critical"
    compliance_requirements: ["SOX", "PCI-DSS"]
    database_host: "db.prod.company.com"
    database_port: "5432"

  - name: "team-stripe-live"
    title: "Stripe Live API Key"
    description: "Production Stripe API key for payment processing"
    environment: "production"
    category: "api_key"
    owner: "payments-team"
    contact: "payments@company.com"
    rotation_schedule: "quarterly"
    sensitive_level: "critical"
    compliance_requirements: ["PCI-DSS"]
    vendor: "Stripe"
    api_version: "2024-06-20"
    permissions: ["charges", "customers", "subscriptions"]
```

## Describe Command Examples

### Secret with full configuration
```bash
$ gsecutil describe team-db-prod
Name: projects/team-secrets-prod/secrets/team-db-prod
Created: 2025-01-15T10:30:00Z
ETag: "abc123def456"
Labels:
  managed_by: gsecutil
  environment: production
  team: backend

Replication Strategy: Automatic (multi-region)

Default Version:
  Version: 3
  State: ENABLED
  Created: 2025-01-20T14:22:15Z
  ETag: "def456ghi789"

Config Attributes:
  Title: Production Database
  Description: PostgreSQL master database connection string
  Environment: production
  Category: database
  Owner: backend-team
  Contact: backend-team@company.com
  Rotation Schedule: quarterly
  Sensitive Level: critical
  Compliance Requirements: ["SOX", "PCI-DSS"]
  Database Host: db.prod.company.com
  Database Port: 5432
```

### Secret with many custom attributes
```bash
$ gsecutil describe team-stripe-live
Name: projects/team-secrets-prod/secrets/team-stripe-live
Created: 2025-01-10T09:15:30Z
ETag: "ghi789jkl012"
Labels:
  managed_by: gsecutil
  environment: production
  vendor: stripe

Replication Strategy: Automatic (multi-region)

Default Version:
  Version: 1
  State: ENABLED
  Created: 2025-01-10T09:15:30Z
  ETag: "jkl012mno345"

Config Attributes:
  Title: Stripe Live API Key
  Description: Production Stripe API key for payment processing
  Environment: production
  Category: api_key
  Owner: payments-team
  Contact: payments@company.com
  Rotation Schedule: quarterly
  Sensitive Level: critical
  Compliance Requirements: ["PCI-DSS"]
  Vendor: Stripe
  API Version: 2024-06-20
  Permissions: ["charges", "customers", "subscriptions"]
```

### Secret not in configuration file
```bash
$ gsecutil describe team-legacy-secret
Name: projects/team-secrets-prod/secrets/team-legacy-secret
Created: 2024-12-01T16:45:20Z
ETag: "old123abc456"
Labels:
  (none)

Replication Strategy: Automatic (multi-region)

Default Version:
  Version: 5
  State: ENABLED
  Created: 2024-12-15T11:30:45Z
  ETag: "old456def789"

Config Attributes:
  (No configuration found for this secret)
```

### Secret outside prefix scope
```bash
$ gsecutil describe other-teams-secret
Name: projects/team-secrets-prod/secrets/other-teams-secret
Created: 2025-01-05T12:20:10Z
ETag: "xyz789abc123"
Labels:
  team: other-team
  managed_by: terraform

Replication Strategy: User-managed (us-central1, us-east1)

Default Version:
  Version: 2
  State: ENABLED
  Created: 2025-01-18T08:45:30Z
  ETag: "xyz123def456"

Config Attributes:
  (Secret not managed by this configuration)
```

## Describe with Version History

```bash
$ gsecutil describe team-db-prod --show-versions
Name: projects/team-secrets-prod/secrets/team-db-prod
Created: 2025-01-15T10:30:00Z
ETag: "abc123def456"
Labels:
  managed_by: gsecutil
  environment: production
  team: backend

Replication Strategy: Automatic (multi-region)

Config Attributes:
  Title: Production Database
  Description: PostgreSQL master database connection string
  Environment: production
  Category: database
  Owner: backend-team
  Contact: backend-team@company.com
  Rotation Schedule: quarterly
  Sensitive Level: critical
  Compliance Requirements: ["SOX", "PCI-DSS"]
  Database Host: db.prod.company.com
  Database Port: 5432

--- Versions ---

Version: 3 (DEFAULT)
  State: ENABLED
  Created: 2025-01-20T14:22:15Z
  ETag: "def456ghi789"

Version: 2
  State: DISABLED
  Created: 2025-01-18T09:10:22Z
  ETag: "ghi789jkl012"

Version: 1
  State: DESTROYED
  Created: 2025-01-15T10:30:00Z
  Destroy Time: 2025-01-19T13:45:18Z
  ETag: "jkl012mno345"
```

## Benefits for Teams

### 1. Complete Documentation
Each secret shows both technical metadata and business context:
- **Secret Manager data**: Creation time, versions, replication strategy
- **Team context**: Purpose, owner, compliance requirements, rotation schedule

### 2. Centralized Information
All secret documentation is in one place:
- No need to check multiple systems for secret details
- Configuration file serves as team documentation
- Easy to see relationships and dependencies

### 3. Compliance Support
Built-in compliance tracking:
- See which secrets require specific compliance standards
- Track rotation schedules and ownership
- Document sensitivity levels and access requirements

### 4. Operational Context
Rich operational metadata:
- Contact information for troubleshooting
- Environment classification
- Service dependencies and configurations

## Usage Patterns

### Security Review
```bash
# Review all production secrets
gsecutil list --filter-attributes environment=production
for secret in $(gsecutil list --filter-attributes environment=production --format json | jq -r '.[].name'); do
  echo "=== $secret ==="
  gsecutil describe "$secret"
  echo ""
done
```

### Compliance Audit
```bash
# Find all PCI-DSS secrets
gsecutil list --filter-attributes compliance_requirements=PCI-DSS
gsecutil describe team-stripe-live | grep -A20 "Config Attributes:"
```

### Rotation Planning
```bash
# Find quarterly rotation secrets
gsecutil list --filter-attributes rotation_schedule=quarterly --show-attributes title,owner,contact
```

### Team Handover
```bash
# Document all backend team secrets
gsecutil list --filter-attributes owner=backend-team
for secret in $(gsecutil list --filter-attributes owner=backend-team --format json | jq -r '.[].name'); do
  gsecutil describe "$secret" > "docs/secrets/${secret}.txt"
done
```
