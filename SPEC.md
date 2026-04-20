# Crossplane Provider for Incident.io — Spec

## Overview

This provider enables Kubernetes-native management of
[Incident.io](https://incident.io) resources via Crossplane. It is
generated using [Upjet](https://github.com/crossplane/upjet) from the
official [Terraform provider](https://registry.terraform.io/providers/incident-io/incident/latest/docs)
(`incident-io/incident`).

## Provider Configuration

### ProviderConfig

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

The provider authenticates via an Incident.io API key stored in a
Kubernetes Secret. The API key requires sufficient permissions for the
resources being managed.

### Credential Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: incident-io-credentials
  namespace: crossplane-system
type: Opaque
stringData:
  api-key: <INCIDENT_IO_API_KEY>
```

## API Group

All managed resources use the API group `incidentio.crossplane.io`.
Initial version: `v1alpha1`.

## Managed Resources

Derived from the Terraform provider's resource schemas. Each Terraform
resource maps to a Crossplane Managed Resource (MR).

### Alerting

| Terraform Resource | Crossplane Kind | Description |
|---|---|---|
| `incident_alert_attribute` | `AlertAttribute` | Alert attribute definitions |
| `incident_alert_route` | `AlertRoute` | Alert routing rules |
| `incident_alert_source` | `AlertSource` | Alert source configurations |

### Catalog

| Terraform Resource | Crossplane Kind | Description |
|---|---|---|
| `incident_catalog_type` | `CatalogType` | Catalog type definitions |
| `incident_catalog_type_attribute` | `CatalogTypeAttribute` | Attributes on catalog types |
| `incident_catalog_entry` | `CatalogEntry` | Individual catalog entries |
| `incident_catalog_entries` | `CatalogEntries` | Bulk catalog entry management |

### Incident Configuration

| Terraform Resource | Crossplane Kind | Description |
|---|---|---|
| `incident_custom_field` | `CustomField` | Custom field definitions |
| `incident_custom_field_option` | `CustomFieldOption` | Options for single/multi-select custom fields |
| `incident_role` | `Role` | Incident role definitions |
| `incident_severity` | `Severity` | Severity level definitions |
| `incident_status` | `Status` | Incident status definitions |

### On-Call

| Terraform Resource | Crossplane Kind | Description |
|---|---|---|
| `incident_escalation_path` | `EscalationPath` | Escalation path configurations |
| `incident_schedule` | `Schedule` | On-call schedule definitions |
| `incident_maintenance_window` | `MaintenanceWindow` | Maintenance window definitions |

### Automation

| Terraform Resource | Crossplane Kind | Description |
|---|---|---|
| `incident_workflow` | `Workflow` | Incident workflow definitions |

## Example Managed Resource

```yaml
apiVersion: incidentio.crossplane.io/v1alpha1
kind: Severity
metadata:
  name: critical
spec:
  forProvider:
    name: Critical
    description: Total system outage or data loss
    rank: 1
  providerConfigRef:
    name: default
```

```yaml
apiVersion: incidentio.crossplane.io/v1alpha1
kind: CustomField
metadata:
  name: affected-service
spec:
  forProvider:
    name: Affected Service
    description: The service impacted by this incident
    fieldType: single_select
  providerConfigRef:
    name: default
```

```yaml
apiVersion: incidentio.crossplane.io/v1alpha1
kind: EscalationPath
metadata:
  name: platform-engineering
spec:
  forProvider:
    name: Platform Engineering
    path:
      - id: level-1
        type: level
        targets:
          - type: schedule
            id: platform-primary
            urgency: high
        timeToAckSeconds: 300
      - id: level-2
        type: level
        targets:
          - type: user
            id: <USER_ID>
            urgency: high
  providerConfigRef:
    name: default
```

## Project Structure

Standard Upjet provider layout:

```
crossplane-provider-incident-io/
├── apis/                          # Generated CRD types
│   ├── alerting/
│   │   └── v1alpha1/              # AlertAttribute, AlertRoute, AlertSource
│   ├── catalog/
│   │   └── v1alpha1/              # CatalogType, CatalogTypeAttribute, CatalogEntry, CatalogEntries
│   ├── incident/
│   │   └── v1alpha1/              # CustomField, CustomFieldOption, Role, Severity, Status
│   ├── oncall/
│   │   └── v1alpha1/              # EscalationPath, Schedule, MaintenanceWindow
│   ├── automation/
│   │   └── v1alpha1/              # Workflow
│   └── v1alpha1/
│       └── providerconfig.go      # ProviderConfig type
├── config/
│   ├── provider.go                # Upjet provider configuration
│   ├── alerting/                  # Per-resource Upjet config overrides
│   ├── catalog/
│   ├── incident/
│   ├── oncall/
│   └── automation/
├── internal/
│   └── controller/                # Generated controllers
├── package/
│   ├── crds/                      # Generated CRD YAML
│   └── crossplane.yaml            # Provider package metadata
├── cmd/
│   ├── provider/
│   │   └── main.go                # Provider binary entrypoint
│   └── generator/
│       └── main.go                # Upjet code generator entrypoint
├── Makefile
├── go.mod
└── README.md
```

## Build and Publish

1. Generate types and controllers from Terraform provider schemas
   using Upjet
2. Build provider binary as a multi-arch container image
3. Package as a Crossplane provider package (xpkg)
4. Publish to a container registry (GitHub Container Registry)

```bash
make generate        # Run Upjet code generation
make build           # Build provider binary
make docker-build    # Build container image
make docker-push     # Push to registry
make xpkg-build      # Build Crossplane package
make xpkg-push       # Push package to registry
```

## Integration with aws-eks-modules

The provider will be installed via the `crossplane` Flux module in
`aws-eks-modules`:

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-incident-io
spec:
  package: ghcr.io/avodah-inc/crossplane-provider-incident-io:v0.1.0
  controllerConfigRef:
    name: provider-incident-io
```

## Security

- API key stored in Kubernetes Secret, referenced via ProviderConfig
- No credentials in CRDs or provider package
- Trivy scan must pass with no critical/high/medium issues
- Provider container image is distroless

## Validation

All generated CRDs validated with kubeconform:

```bash
kubeconform \
  -schema-location default \
  -schema-location "https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json" \
  -strict \
  -summary package/crds/
```

## References

- [Terraform Provider](https://registry.terraform.io/providers/incident-io/incident/latest/docs)
- [Terraform Provider Source](https://github.com/incident-io/terraform-provider-incident)
- [Upjet Documentation](https://github.com/crossplane/upjet/blob/main/docs/generating-a-provider.md)
- [Crossplane Provider Development](https://docs.crossplane.io/latest/guides/write-a-composition-function-in-go/)
