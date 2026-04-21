# Crossplane Provider for Incident.io

[![CI](https://github.com/avodah-inc/crossplane-provider-incident-io/actions/workflows/ci.yml/badge.svg)](https://github.com/avodah-inc/crossplane-provider-incident-io/actions/workflows/ci.yml)
[![Release](https://github.com/avodah-inc/crossplane-provider-incident-io/actions/workflows/release.yml/badge.svg)](https://github.com/avodah-inc/crossplane-provider-incident-io/actions/workflows/release.yml)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

`provider-incident-io` is a [Crossplane](https://crossplane.io/) provider generated
with [Upjet](https://github.com/crossplane/upjet) from the official
[incident-io/incident](https://registry.terraform.io/providers/incident-io/incident/latest/docs)
Terraform provider.

It brings Kubernetes-native, declarative management to
[Incident.io](https://incident.io/) resources — alerting configuration,
service catalog, incident settings, on-call schedules, and automation workflows
all reconciled through the standard Crossplane lifecycle.

## Supported Resources

16 managed resources across 5 domain groups, all under `incidentio.crossplane.io/v1alpha1`:

| Domain | API Group | Kind | Terraform Resource |
|--------|-----------|------|--------------------|
| Alerting | `alerting.incidentio.crossplane.io` | `Attribute` | `incident_alert_attribute` |
| Alerting | `alerting.incidentio.crossplane.io` | `Route` | `incident_alert_route` |
| Alerting | `alerting.incidentio.crossplane.io` | `Source` | `incident_alert_source` |
| Catalog | `catalog.incidentio.crossplane.io` | `CatalogType` | `incident_catalog_type` |
| Catalog | `catalog.incidentio.crossplane.io` | `TypeAttribute` | `incident_catalog_type_attribute` |
| Catalog | `catalog.incidentio.crossplane.io` | `Entry` | `incident_catalog_entry` |
| Catalog | `catalog.incidentio.crossplane.io` | `CatalogEntries` | `incident_catalog_entries` |
| Incident | `incident.incidentio.crossplane.io` | `Field` | `incident_custom_field` |
| Incident | `incident.incidentio.crossplane.io` | `FieldOption` | `incident_custom_field_option` |
| Incident | `incident.incidentio.crossplane.io` | `Role` | `incident_incident_role` |
| Incident | `incident.incidentio.crossplane.io` | `Severity` | `incident_severity` |
| Incident | `incident.incidentio.crossplane.io` | `Status` | `incident_status` |
| On-Call | `oncall.incidentio.crossplane.io` | `Path` | `incident_escalation_path` |
| On-Call | `oncall.incidentio.crossplane.io` | `Schedule` | `incident_schedule` |
| On-Call | `oncall.incidentio.crossplane.io` | `Window` | `incident_maintenance_window` |
| Automation | `automation.incidentio.crossplane.io` | `Workflow` | `incident_workflow` |

## Quick Start

### 1. Install the Provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-incident-io
spec:
  package: ghcr.io/avodah-inc/crossplane-provider-incident-io:v0.1.0
```

### 2. Create a Secret with Your API Key

Generate an API key at [app.incident.io/settings/api-keys](https://app.incident.io/settings/api-keys).
The key is org-scoped and remains valid even if the creating user is deactivated.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: incident-io-credentials
  namespace: crossplane-system
type: Opaque
stringData:
  api-key: "YOUR_INCIDENT_IO_API_KEY"
```

### 3. Create a ProviderConfig

```yaml
apiVersion: incidentio.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: incident-io-credentials
      key: api-key
```

### 4. Declare Resources

```yaml
apiVersion: incident.incidentio.crossplane.io/v1alpha1
kind: Severity
metadata:
  name: critical
spec:
  forProvider:
    name: "Critical"
    description: "Critical severity - major customer impact requiring immediate response"
    rank: 1
  providerConfigRef:
    name: default
```

```yaml
apiVersion: catalog.incidentio.crossplane.io/v1alpha1
kind: CatalogType
metadata:
  name: service
spec:
  forProvider:
    name: "Service"
    description: "A service in the product catalog"
    sourceRepoUrl: "https://github.com/avodah-inc/crossplane-provider-incident-io"
    typeName: 'Custom["Service"]'
    categories:
      - "service"
  providerConfigRef:
    name: default
```

```yaml
apiVersion: alerting.incidentio.crossplane.io/v1alpha1
kind: Attribute
metadata:
  name: service-tier
spec:
  forProvider:
    name: "service_tier"
    type: "single_select"
    array: false
  providerConfigRef:
    name: default
```

More examples in the [`examples/`](examples/) directory.

## How It Works

The provider follows the standard Crossplane reconciliation model:

1. You declare an Incident.io resource as a Kubernetes CR
2. The provider's controller observes the CR and calls the Incident.io API
3. On each reconciliation loop, the controller compares desired state (spec) with actual state (API)
4. Drift is corrected automatically — if someone changes a resource in the Incident.io UI, the provider restores it
5. Status conditions (`Ready`, `Synced`) report convergence state on each CR

Under the hood, Upjet delegates to the Terraform provider's CRUD logic, so API parity with the upstream Terraform provider is maintained automatically.

### Status Conditions

| Condition | Meaning |
|-----------|---------|
| `Ready=True` | Resource exists in Incident.io and matches desired state |
| `Synced=True` | Last reconciliation succeeded |
| `Synced=False` | Last reconciliation failed — check `message` for details |

### Authentication

The provider reads an API key from a Kubernetes Secret referenced by the ProviderConfig.
If the secret is missing or the key is invalid, the ProviderConfig reports `Ready=False`.

Rate limit: 1200 requests/minute per API key. The provider uses Terraform's built-in retry
with exponential backoff for 429 responses.

## Development

### Prerequisites

- Go 1.23+
- Python 3 (for code generation scripts)
- Docker with Buildx (for container builds)
- [controller-gen](https://book.kubebuilder.io/reference/controller-gen) v0.16.5+
- [Crossplane CLI](https://docs.crossplane.io/latest/cli/) (`crank`)

### Build

```bash
# Full pipeline: generate types/CRDs/controllers, then compile
make generate
make build

# Or both at once
make all
```

### Code Generation

`make generate` runs the Upjet pipeline end-to-end:

1. Reads the Terraform provider schema from `config/schema.json`
2. Applies provider configuration from `config/provider.go` and per-domain overrides
3. Generates Go types under `apis/<domain>/v1alpha1/`
4. Generates managed resource interface methods via `hack/generate-managed.py`
5. Generates deepcopy methods via `controller-gen`
6. Generates CRD YAML under `package/crds/`

To regenerate after upstream Terraform provider changes:

```bash
# Extract fresh schema from the Terraform provider
terraform providers schema -json > config/schema.json

# Convert protocol v6 nested_type format to block_types (Upjet v1.11 compat)
python3 hack/convert-schema.py config/schema.json config/schema.json

# Regenerate everything
make generate
```

### Test

```bash
# Unit tests
go test -race -coverprofile=coverage.out ./cmd/... ./config/... ./internal/clients/...

# Smoke tests (verify generated artifacts exist and are well-formed)
cd test && go test -v -count=1 ./...

# Integration tests (requires kind, helm, kubectl)
./test/integration/install_test.sh
```

### Container Image

```bash
# Build multi-arch image (linux/amd64 + linux/arm64, distroless base)
make docker-build

# Push to GHCR
make docker-push VERSION=0.2.0
```

### Crossplane Package

```bash
# Build xpkg archive
make xpkg-build

# Push to GHCR
make xpkg-push VERSION=0.2.0
```

### Validation and Security

```bash
# Validate CRDs against Kubernetes schemas
make validate-crds

# Trivy filesystem scan (fails on CRITICAL/HIGH/MEDIUM)
make trivy-fs

# Trivy container image scan
make trivy-image

# SonarQube analysis
make sonar

# Lint
make lint
```

## CI/CD

Two GitHub Actions workflows handle the full lifecycle:

### CI (`.github/workflows/ci.yml`)

Runs on every push to main and every PR. Six parallel jobs:

- **test-unit** — Go tests with race detection and coverage
- **test-integration** — CRD validation with kubeconform, smoke tests
- **trivy-fs** — Filesystem vulnerability scan (SARIF upload + severity gate)
- **trivy-image** — Container image vulnerability scan (SARIF upload + severity gate)
- **sbom** — CycloneDX SBOM generation (90-day artifact retention)
- **codeql** — CodeQL static analysis for Go

### Release (`.github/workflows/release.yml`)

Triggered automatically when CI passes on main. Uses [semantic-release](https://semantic-release.gitbook.io/)
with Angular commit conventions to determine version bumps:

- `feat(...)` → minor version
- `fix(...)` → patch version
- `BREAKING CHANGE` → major version

When a new version is published:

1. Multi-arch Docker image is built and pushed to GHCR
2. Crossplane xpkg is built with the embedded runtime image and pushed to GHCR
3. SBOM and xpkg are attached to the GitHub Release

## Project Structure

```
.
├── apis/                          # Generated Go types per domain
│   ├── alerting/v1alpha1/         #   Alert attributes, routes, sources
│   ├── automation/v1alpha1/       #   Workflows
│   ├── catalog/v1alpha1/          #   Catalog types, entries, attributes
│   ├── incident/v1alpha1/         #   Custom fields, roles, severities, statuses
│   ├── oncall/v1alpha1/           #   Escalation paths, schedules, maintenance windows
│   └── v1alpha1/                  #   ProviderConfig and ProviderConfigUsage
├── cmd/
│   ├── generator/main.go         # Upjet code generation entrypoint
│   └── provider/main.go          # Provider binary entrypoint
├── config/
│   ├── provider.go               # Central Upjet provider configuration
│   ├── external_name.go          # External name mappings (all IdentifierFromProvider)
│   ├── schema.json               # Terraform provider schema (embedded at build)
│   ├── alerting/config.go        # Alerting domain resource overrides
│   ├── automation/config.go      # Automation domain resource overrides
│   ├── catalog/config.go         # Catalog domain resource overrides
│   ├── incident/config.go        # Incident domain resource overrides
│   └── oncall/config.go          # On-call domain resource overrides
├── internal/
│   ├── clients/incident.go       # TerraformSetupBuilder (credential handling)
│   ├── controller/               # Generated + hand-written controllers
│   └── features/                 # Feature flags
├── package/
│   ├── crossplane.yaml           # Crossplane provider package metadata
│   └── crds/                     # Generated CRD YAML (18 files)
├── examples/                     # Example manifests for each domain
├── hack/
│   ├── convert-schema.py         # TF protocol v6 → v5 schema converter
│   └── generate-managed.py       # Managed resource interface generator
├── test/
│   ├── smoke_test.go             # Code generation output verification
│   └── integration/              # Kind cluster integration tests
├── .github/workflows/
│   ├── ci.yml                    # CI pipeline (tests, scans, SBOM, CodeQL)
│   └── release.yml               # Semantic release, image + xpkg publish
├── Dockerfile                    # Multi-stage distroless build
├── Makefile                      # All build/test/scan/publish targets
├── .releaserc.json               # Semantic release configuration
└── sonar-project.properties      # SonarQube configuration
```

## Compatibility

| Component | Version |
|-----------|---------|
| Crossplane | >= 1.14.0 |
| Kubernetes | >= 1.27 |
| Terraform Provider (incident-io/incident) | 5.35.0 |
| Go | 1.23 |
| Architectures | linux/amd64, linux/arm64 |

## Upgrading the Terraform Provider

When a new version of the `incident-io/incident` Terraform provider is released:

1. Update the dependency in `go.mod`
2. Re-extract the provider schema: `terraform providers schema -json > config/schema.json`
3. Convert the schema: `python3 hack/convert-schema.py config/schema.json config/schema.json`
4. Regenerate: `make generate`
5. Verify: `make build && make validate-crds`
6. Run tests: `go test ./... && cd test && go test -v ./...`

New resources added by the Terraform provider will need entries in `config/external_name.go`
and a domain config file under `config/<domain>/config.go`.

## License

This project is licensed under the GNU General Public License v3.0. See [LICENSE](LICENSE) for the full text.
