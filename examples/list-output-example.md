# gsecutil list Command with Configuration

This document shows how the `gsecutil list` command output would look when using configuration files with metadata.

## Configuration File Example

```yaml
# gsecutil.conf
project: "team-secrets-prod"
prefix: "team-"

# Configure list command display
list:
  attributes:
    - title
    - owner
    - environment
    - sensitive_level

credentials:
  - name: "team-db-prod"
    title: "Production Database"
    description: "PostgreSQL master database"
    environment: "production"
    owner: "backend-team"
    sensitive_level: "critical"

  - name: "team-db-staging"
    title: "Staging Database"
    description: "PostgreSQL staging database"
    environment: "staging"
    owner: "backend-team"
    sensitive_level: "medium"

  - name: "team-stripe-live"
    title: "Stripe Live API Key"
    description: "Production payment processing"
    environment: "production"
    owner: "payments-team"
    sensitive_level: "critical"
    vendor: "Stripe"

  - name: "team-sendgrid-api"
    title: "SendGrid API Key"
    description: "Email delivery service"
    environment: "production"
    owner: "platform-team"
    sensitive_level: "medium"
    vendor: "SendGrid"
```

## List Command Examples

### Basic list (default output from config)
```bash
$ gsecutil list
NAME               TITLE                    OWNER           ENVIRONMENT      SENSITIVE_LEVEL
team-db-prod       Production Database      backend-team    production       critical
team-db-staging    Staging Database         backend-team    staging          medium
team-stripe-live   Stripe Live API Key      payments-team   production       critical
team-sendgrid-api  SendGrid API Key         platform-team   production       medium
```

### List with custom attribute display
```bash
$ gsecutil list --show-attributes title,owner
NAME               TITLE                    OWNER
team-db-prod       Production Database      backend-team
team-db-staging    Staging Database         backend-team
team-stripe-live   Stripe Live API Key      payments-team
team-sendgrid-api  SendGrid API Key         platform-team

$ gsecutil list --show-attributes title,vendor
NAME               TITLE                    VENDOR
team-db-prod       Production Database      (none)
team-db-staging    Staging Database         (none)
team-stripe-live   Stripe Live API Key      Stripe
team-sendgrid-api  SendGrid API Key         SendGrid
```

### List with attribute filtering
```bash
$ gsecutil list --filter-attributes environment=production
NAME               TITLE                    OWNER           ENVIRONMENT      SENSITIVE_LEVEL
team-db-prod       Production Database      backend-team    production       critical
team-stripe-live   Stripe Live API Key      payments-team   production       critical
team-sendgrid-api  SendGrid API Key         platform-team   production       medium

$ gsecutil list --filter-attributes owner=backend-team
NAME               TITLE                    OWNER           ENVIRONMENT      SENSITIVE_LEVEL
team-db-prod       Production Database      backend-team    production       critical
team-db-staging    Staging Database         backend-team    staging          medium

$ gsecutil list --filter-attributes sensitive_level=critical
NAME               TITLE                    OWNER           ENVIRONMENT      SENSITIVE_LEVEL
team-db-prod       Production Database      backend-team    production       critical
team-stripe-live   Stripe Live API Key      payments-team   production       critical

$ gsecutil list --filter-attributes environment=production --show-attributes title,vendor
NAME               TITLE                    VENDOR
team-db-prod       Production Database      (none)
team-stripe-live   Stripe Live API Key      Stripe
team-sendgrid-api  SendGrid API Key         SendGrid
```

### List with detailed output
```bash
$ gsecutil list --detailed
NAME: team-db-prod
  Title: Production Database
  Description: PostgreSQL master database
  Environment: production
  Owner: backend-team
  Sensitive Level: critical
  Created: 2025-01-15T10:30:00Z
  Versions: 3

NAME: team-db-staging
  Title: Staging Database
  Description: PostgreSQL staging database
  Environment: staging
  Owner: backend-team
  Sensitive Level: medium
  Created: 2025-01-10T14:20:00Z
  Versions: 2

NAME: team-stripe-live
  Title: Stripe Live API Key
  Description: Production payment processing
  Environment: production
  Owner: payments-team
  Sensitive Level: critical
  Vendor: Stripe
  Created: 2025-01-20T09:15:00Z
  Versions: 1

NAME: team-sendgrid-api
  Title: SendGrid API Key
  Description: Email delivery service
  Environment: production
  Owner: platform-team
  Sensitive Level: medium
  Vendor: SendGrid
  Created: 2025-01-18T16:45:00Z
  Versions: 1
```

### JSON output for scripting
```bash
$ gsecutil list --format json
[
  {
    "name": "team-db-prod",
    "title": "Production Database",
    "description": "PostgreSQL master database",
    "environment": "production",
    "owner": "backend-team",
    "sensitive_level": "critical",
    "created": "2025-01-15T10:30:00Z",
    "versions": 3,
    "latest_version": "3"
  },
  {
    "name": "team-db-staging",
    "title": "Staging Database",
    "description": "PostgreSQL staging database",
    "environment": "staging",
    "owner": "backend-team",
    "sensitive_level": "medium",
    "created": "2025-01-10T14:20:00Z",
    "versions": 2,
    "latest_version": "2"
  }
]
```

## Benefits for Teams

1. **Clear Overview**: Everyone can quickly understand what each secret is for
2. **Ownership Clarity**: Easy to see who's responsible for each credential
3. **Environment Safety**: Filter to avoid accidentally using production secrets
4. **Security Classification**: Understand sensitivity levels for compliance
5. **Vendor Tracking**: Know which external services are involved
6. **Rotation Planning**: See what needs to be updated and when

## Behavior Without Configuration

### No config file or no credentials section
```bash
# Without config file - only shows secret names
$ gsecutil list
NAME
my-secret-1
my-secret-2
another-secret

# Force showing attributes when no config (shows empty values)
$ gsecutil list --show-attributes title,owner
NAME               TITLE                    OWNER
my-secret-1        (no title)              (unknown)
my-secret-2        (no title)              (unknown)
another-secret     (no title)              (unknown)
```

### Config file with project/prefix but no credentials section
```yaml
# gsecutil.conf
project: "my-project"
prefix: "team-"
# No credentials section
```

```bash
# Shows only secret names (no attributes since no credentials defined)
$ gsecutil list
NAME
team-secret-1
team-secret-2

# Can force attribute display with CLI parameter
$ gsecutil list --show-attributes title
NAME               TITLE
team-secret-1      (no title)
team-secret-2      (no title)
```

## Integration with Existing Secrets

If you have existing secrets in Secret Manager that don't have configuration metadata:

```bash
# List shows configured secrets with metadata, others with basic info
$ gsecutil list
NAME                    TITLE                    OWNER           ENVIRONMENT
team-db-prod           Production Database      backend-team    production
team-stripe-live       Stripe Live API Key      payments-team   production
team-legacy-secret     (no title)              (unknown)        (unknown)
team-undocumented      (no title)              (unknown)        (unknown)

# Use --all to see everything in the project
$ gsecutil list --all
# Shows all secrets, not just those with team prefix

# Use --no-prefix-filter to see all team secrets plus others
$ gsecutil list --no-prefix-filter
# Shows all secrets starting with team prefix, plus any documented in config
```
