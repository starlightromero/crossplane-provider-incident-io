# Requirements Document

## Introduction

This document defines the requirements for a Crossplane Provider that enables Kubernetes-native management of Incident.io resources. The provider is generated using Upjet from the official Terraform provider (`incident-io/incident`). It allows platform teams to declare Incident.io resources (severities, custom fields, escalation paths, schedules, catalog entries, workflows, etc.) as Kubernetes Custom Resources, managed through the standard Crossplane reconciliation lifecycle.

## Glossary

- **Provider**: A Crossplane provider binary that reconciles managed resources against an external API (Incident.io)
- **ProviderConfig**: A Crossplane custom resource that holds authentication credentials for the Provider
- **Managed_Resource**: A Kubernetes custom resource representing an external Incident.io resource, reconciled by the Provider
- **CRD**: A Kubernetes Custom Resource Definition that extends the Kubernetes API with a new resource type
- **Upjet**: A Crossplane framework that generates Crossplane providers from Terraform provider schemas
- **Terraform_Provider**: The `incident-io/incident` Terraform provider whose schemas are used to generate CRDs and controllers
- **API_Group**: The Kubernetes API group (`incidentio.crossplane.io`) under which all Managed Resources are registered
- **xpkg**: A Crossplane provider package format used for distribution and installation
- **Reconciliation**: The Crossplane controller loop that observes, creates, updates, or deletes external resources to match desired state
- **kubeconform**: A Kubernetes manifest validation tool used to verify CRD structural correctness
- **Trivy**: A security scanner used to detect vulnerabilities in container images and filesystems
- **Distroless_Image**: A minimal container image containing only the application binary and runtime dependencies, with no shell or package manager

## Requirements

### Requirement 1: Upjet Code Generation

**User Story:** As a platform engineer, I want the provider's CRDs and controllers generated from the Terraform provider schemas, so that the Crossplane resources stay consistent with the upstream Incident.io API.

#### Acceptance Criteria

1. WHEN the `make generate` command is executed, THE Provider SHALL produce Go type definitions and controllers for all 16 Terraform resources defined in the `incident-io/incident` Terraform_Provider.
2. THE Provider SHALL generate CRDs organized into five domain groups: `alerting`, `catalog`, `incident`, `oncall`, and `automation`.
3. THE Provider SHALL generate all CRDs under the API_Group `incidentio.crossplane.io` with version `v1alpha1`.
4. WHEN a new version of the Terraform_Provider is released, THE Provider SHALL support regeneration of types and controllers by re-running the Upjet code generation pipeline.

### Requirement 2: Provider Authentication

**User Story:** As a platform engineer, I want the provider to authenticate with Incident.io using an API key stored in a Kubernetes Secret, so that credentials are managed securely.

#### Acceptance Criteria

1. THE ProviderConfig SHALL reference an Incident.io API key via a `secretRef` pointing to a Kubernetes Secret.
2. WHEN a ProviderConfig resource is created with a valid `secretRef`, THE Provider SHALL read the API key from the referenced Kubernetes Secret.
3. IF the referenced Kubernetes Secret does not exist, THEN THE Provider SHALL set the ProviderConfig status condition to `Ready=False` with a reason indicating the missing secret.
4. IF the API key in the referenced Secret is invalid or lacks sufficient permissions, THEN THE Provider SHALL set the ProviderConfig status condition to `Ready=False` with a reason indicating an authentication failure.
5. THE ProviderConfig SHALL accept credentials only from Kubernetes Secrets (source: `Secret`).

### Requirement 3: Alerting Managed Resources

**User Story:** As a platform engineer, I want to manage Incident.io alerting resources (alert attributes, alert routes, alert sources) as Kubernetes custom resources, so that alerting configuration is version-controlled and declarative.

#### Acceptance Criteria

1. THE Provider SHALL generate a Managed_Resource `AlertAttribute` (kind: `AlertAttribute`) that maps to the `incident_alert_attribute` Terraform resource.
2. THE Provider SHALL generate a Managed_Resource `AlertRoute` (kind: `AlertRoute`) that maps to the `incident_alert_route` Terraform resource.
3. THE Provider SHALL generate a Managed_Resource `AlertSource` (kind: `AlertSource`) that maps to the `incident_alert_source` Terraform resource.
4. WHEN an AlertAttribute, AlertRoute, or AlertSource resource is created in Kubernetes, THE Provider SHALL create the corresponding resource in Incident.io via the API.
5. WHEN an AlertAttribute, AlertRoute, or AlertSource resource is updated in Kubernetes, THE Provider SHALL update the corresponding resource in Incident.io via the API.
6. WHEN an AlertAttribute, AlertRoute, or AlertSource resource is deleted from Kubernetes, THE Provider SHALL delete the corresponding resource from Incident.io via the API.

### Requirement 4: Catalog Managed Resources

**User Story:** As a platform engineer, I want to manage Incident.io catalog resources (catalog types, type attributes, entries) as Kubernetes custom resources, so that the service catalog is declaratively managed.

#### Acceptance Criteria

1. THE Provider SHALL generate Managed Resources for `CatalogType`, `CatalogTypeAttribute`, `CatalogEntry`, and `CatalogEntries` that map to their respective Terraform resources.
2. WHEN a CatalogType resource is created in Kubernetes, THE Provider SHALL create the corresponding catalog type in Incident.io via the API.
3. WHEN a CatalogEntry resource is created in Kubernetes, THE Provider SHALL create the corresponding catalog entry in Incident.io via the API.
4. WHEN a CatalogEntries resource is created in Kubernetes, THE Provider SHALL perform bulk creation of catalog entries in Incident.io via the API.
5. WHEN any Catalog Managed_Resource is updated in Kubernetes, THE Provider SHALL update the corresponding resource in Incident.io via the API.
6. WHEN any Catalog Managed_Resource is deleted from Kubernetes, THE Provider SHALL delete the corresponding resource from Incident.io via the API.

### Requirement 5: Incident Configuration Managed Resources

**User Story:** As a platform engineer, I want to manage Incident.io incident configuration resources (custom fields, field options, roles, severities, statuses) as Kubernetes custom resources, so that incident response configuration is codified.

#### Acceptance Criteria

1. THE Provider SHALL generate Managed Resources for `CustomField`, `CustomFieldOption`, `Role`, `Severity`, and `Status` that map to their respective Terraform resources.
2. WHEN a CustomField resource is created in Kubernetes, THE Provider SHALL create the corresponding custom field in Incident.io via the API.
3. WHEN a Severity resource is created in Kubernetes, THE Provider SHALL create the corresponding severity level in Incident.io via the API.
4. WHEN a Status resource is created in Kubernetes, THE Provider SHALL create the corresponding incident status in Incident.io via the API.
5. WHEN any Incident Configuration Managed_Resource is updated in Kubernetes, THE Provider SHALL update the corresponding resource in Incident.io via the API.
6. WHEN any Incident Configuration Managed_Resource is deleted from Kubernetes, THE Provider SHALL delete the corresponding resource from Incident.io via the API.

### Requirement 6: On-Call Managed Resources

**User Story:** As a platform engineer, I want to manage Incident.io on-call resources (escalation paths, schedules, maintenance windows) as Kubernetes custom resources, so that on-call configuration is declarative and auditable.

#### Acceptance Criteria

1. THE Provider SHALL generate Managed Resources for `EscalationPath`, `Schedule`, and `MaintenanceWindow` that map to their respective Terraform resources.
2. WHEN an EscalationPath resource is created in Kubernetes, THE Provider SHALL create the corresponding escalation path in Incident.io via the API.
3. WHEN a Schedule resource is created in Kubernetes, THE Provider SHALL create the corresponding on-call schedule in Incident.io via the API.
4. WHEN a MaintenanceWindow resource is created in Kubernetes, THE Provider SHALL create the corresponding maintenance window in Incident.io via the API.
5. WHEN any On-Call Managed_Resource is updated in Kubernetes, THE Provider SHALL update the corresponding resource in Incident.io via the API.
6. WHEN any On-Call Managed_Resource is deleted from Kubernetes, THE Provider SHALL delete the corresponding resource from Incident.io via the API.

### Requirement 7: Automation Managed Resources

**User Story:** As a platform engineer, I want to manage Incident.io workflows as Kubernetes custom resources, so that incident automation is version-controlled.

#### Acceptance Criteria

1. THE Provider SHALL generate a Managed_Resource `Workflow` (kind: `Workflow`) that maps to the `incident_workflow` Terraform resource.
2. WHEN a Workflow resource is created in Kubernetes, THE Provider SHALL create the corresponding workflow in Incident.io via the API.
3. WHEN a Workflow resource is updated in Kubernetes, THE Provider SHALL update the corresponding workflow in Incident.io via the API.
4. WHEN a Workflow resource is deleted from Kubernetes, THE Provider SHALL delete the corresponding workflow from Incident.io via the API.

### Requirement 8: Reconciliation and Status Reporting

**User Story:** As a platform engineer, I want the provider to report accurate status conditions on all managed resources, so that I can observe convergence and diagnose issues.

#### Acceptance Criteria

1. WHEN a Managed_Resource is successfully reconciled with Incident.io, THE Provider SHALL set the `Synced` status condition to `True` on the Managed_Resource.
2. WHEN a Managed_Resource is successfully created or observed in Incident.io, THE Provider SHALL set the `Ready` status condition to `True` on the Managed_Resource.
3. IF the Incident.io API returns an error during Reconciliation, THEN THE Provider SHALL set the `Synced` status condition to `False` with the error message from the API.
4. IF the external Incident.io resource is not found during an observe operation, THEN THE Provider SHALL trigger a create operation to restore the resource.
5. THE Provider SHALL populate the `status.atProvider` field with the current state of the external resource after each successful Reconciliation.

### Requirement 9: Container Image Build

**User Story:** As a platform engineer, I want the provider packaged as a multi-arch distroless container image, so that it runs securely in production Kubernetes clusters.

#### Acceptance Criteria

1. THE Provider SHALL be built as a multi-architecture container image supporting `linux/amd64` and `linux/arm64`.
2. THE Provider SHALL use a Distroless_Image as the base image for the final container.
3. THE Provider container image SHALL contain only the provider binary and its runtime dependencies.
4. WHEN `make docker-build` is executed, THE Provider SHALL produce a container image tagged with the current version.

### Requirement 10: Crossplane Package (xpkg)

**User Story:** As a platform engineer, I want the provider distributed as a Crossplane package, so that it can be installed via the standard Crossplane package manager.

#### Acceptance Criteria

1. THE Provider SHALL include a `crossplane.yaml` metadata file that declares the provider package name, version, and Crossplane compatibility.
2. THE Provider SHALL include all generated CRD YAML files in the `package/crds/` directory.
3. WHEN `make xpkg-build` is executed, THE Provider SHALL produce a valid xpkg archive containing the `crossplane.yaml` and all CRDs.
4. WHEN `make xpkg-push` is executed, THE Provider SHALL publish the xpkg to `ghcr.io/avodah-inc/crossplane-provider-incident-io`.

### Requirement 11: CRD Validation

**User Story:** As a platform engineer, I want all generated CRDs validated against Kubernetes schemas, so that malformed CRDs are caught before deployment.

#### Acceptance Criteria

1. WHEN CRDs are generated, THE Provider build pipeline SHALL validate all CRD YAML files in `package/crds/` using kubeconform.
2. THE Provider build pipeline SHALL run kubeconform in strict mode with the default schema location and the CRDs-catalog schema location.
3. IF kubeconform reports any validation errors, THEN THE Provider build pipeline SHALL fail and report the specific validation errors.

### Requirement 12: Security Scanning

**User Story:** As a platform engineer, I want the provider scanned for security vulnerabilities, so that no critical, high, or medium issues ship to production.

#### Acceptance Criteria

1. THE Provider build pipeline SHALL run a Trivy filesystem scan on the provider source code including dev dependencies.
2. THE Provider build pipeline SHALL run a Trivy container image scan on the built provider image.
3. IF Trivy reports any critical, high, or medium severity vulnerabilities, THEN THE Provider build pipeline SHALL fail and report the specific findings.
4. THE Provider build pipeline SHALL run a SonarQube scan on the provider source code.

### Requirement 13: Project Structure

**User Story:** As a developer, I want the provider to follow the standard Upjet provider layout, so that the codebase is familiar to Crossplane contributors.

#### Acceptance Criteria

1. THE Provider SHALL organize generated CRD types under `apis/` with subdirectories for each domain group: `alerting`, `catalog`, `incident`, `oncall`, and `automation`.
2. THE Provider SHALL place Upjet provider configuration and per-resource overrides under `config/`.
3. THE Provider SHALL place generated controllers under `internal/controller/`.
4. THE Provider SHALL place the provider binary entrypoint at `cmd/provider/main.go` and the generator entrypoint at `cmd/generator/main.go`.
5. THE Provider SHALL place generated CRD YAML files under `package/crds/` and the package metadata at `package/crossplane.yaml`.

### Requirement 14: Integration with aws-eks-modules

**User Story:** As a platform engineer, I want the provider installable via the Flux module in aws-eks-modules, so that it integrates with the existing GitOps deployment pipeline.

#### Acceptance Criteria

1. THE Provider package SHALL be published to `ghcr.io/avodah-inc/crossplane-provider-incident-io` with semantic version tags.
2. THE Provider package SHALL be installable as a Crossplane `Provider` resource referencing the published package image.
3. WHEN the Provider is installed in a cluster, THE Provider SHALL register all 16 CRDs under the API_Group `incidentio.crossplane.io`.
