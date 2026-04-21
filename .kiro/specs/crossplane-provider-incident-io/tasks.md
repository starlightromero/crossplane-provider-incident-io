# Implementation Plan: Crossplane Provider for Incident.io

## Overview

This plan implements a Crossplane provider generated via Upjet from the `incident-io/incident` Terraform provider. The implementation proceeds incrementally: project scaffolding → Upjet configuration → code generation → provider runtime → packaging → validation → publishing. Each task builds on the previous, and all generated code is wired into the build pipeline before moving to the next phase.

## Tasks

- [ ] 1. Scaffold project structure and initialize Go module
  - [x] 1.1 Initialize Go module and project root
    - Create `go.mod` with module path `github.com/avodah-inc/crossplane-provider-incident-io`
    - Add core dependencies: `github.com/crossplane/upjet`, `github.com/crossplane/crossplane-runtime`, `github.com/crossplane/crossplane-tools`, `sigs.k8s.io/controller-runtime`
    - Add Terraform provider dependency: `github.com/incident-io/terraform-provider-incident`
    - Create directory structure: `apis/`, `config/`, `internal/`, `cmd/provider/`, `cmd/generator/`, `package/crds/`
    - Create `README.md` with provider overview, installation, and usage instructions
    - _Requirements: 13.1, 13.2, 13.3, 13.4, 13.5_

  - [x] 1.2 Create Makefile with build targets
    - Implement `make generate` target to run Upjet code generation via `cmd/generator/main.go`
    - Implement `make build` target to compile the provider binary from `cmd/provider/main.go`
    - Implement `make docker-build` target to build multi-arch container image (`linux/amd64`, `linux/arm64`) with distroless base
    - Implement `make docker-push` target to push image to `ghcr.io/avodah-inc/crossplane-provider-incident-io`
    - Implement `make xpkg-build` target to build Crossplane provider package
    - Implement `make xpkg-push` target to push xpkg to `ghcr.io/avodah-inc/crossplane-provider-incident-io`
    - Implement `make lint` and `make test` targets
    - _Requirements: 1.1, 9.1, 9.4, 10.3, 10.4_

  - [x] 1.3 Create Dockerfile for multi-arch distroless image
    - Write multi-stage Dockerfile: Go builder stage + distroless runtime stage
    - Use `gcr.io/distroless/static:nonroot` as the final base image
    - Copy only the provider binary into the final image
    - Support `linux/amd64` and `linux/arm64` via `--platform` build arg
    - _Requirements: 9.1, 9.2, 9.3_

- [ ] 2. Implement ProviderConfig type and Terraform setup client
  - [x] 2.1 Create ProviderConfig CRD type
    - Create `apis/v1alpha1/providerconfig_types.go` with `ProviderConfig` and `ProviderConfigUsage` types
    - Define `spec.credentials` with `source: Secret` and `secretRef` (namespace, name, key)
    - Embed `xpv1.ProviderConfigSpec` and `xpv1.ProviderConfigStatus` from crossplane-runtime
    - Generate deepcopy methods and register types with the scheme
    - _Requirements: 2.1, 2.2, 2.5_

  - [x] 2.2 Implement Terraform setup client
    - Create `internal/clients/incident.go` with `TerraformSetupBuilder` function
    - Read ProviderConfig CR and extract `secretRef`
    - Read the referenced Kubernetes Secret and extract the API key
    - Build and return `terraform.Setup` with provider configuration map (`api_key`, `endpoint`)
    - Handle error cases: missing secret (set `Ready=False`), missing key in secret, invalid API key
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 3. Checkpoint - Verify project compiles
  - Ensure `go build ./...` succeeds with no errors, ask the user if questions arise.

- [ ] 4. Configure Upjet provider and external names
  - [x] 4.1 Create central provider configuration
    - Create `config/provider.go` with Upjet `ProviderConfiguration`
    - Set provider name to `incident-io`, root API group to `incidentio.crossplane.io`
    - Define the resource inclusion list for all 16 Terraform resources
    - Configure default resource behavior (late initialization, external name)
    - Wire in all per-domain configuration functions
    - _Requirements: 1.1, 1.2, 1.3_

  - [x] 4.2 Create external name configuration
    - Create `config/external_name.go` mapping all 16 Terraform resources to `IdentifierFromProvider`
    - All Incident.io resources use server-assigned UUIDs
    - _Requirements: 1.1, 8.5_

  - [x] 4.3 Create alerting domain resource configuration
    - Create `config/alerting/config.go`
    - Register `incident_alert_attribute`, `incident_alert_route`, `incident_alert_source`
    - Set API subgroup to `alerting`, version `v1alpha1`
    - Apply any resource-specific overrides (references, sensitive fields)
    - _Requirements: 1.2, 3.1, 3.2, 3.3_

  - [x] 4.4 Create catalog domain resource configuration
    - Create `config/catalog/config.go`
    - Register `incident_catalog_type`, `incident_catalog_type_attribute`, `incident_catalog_entry`, `incident_catalog_entries`
    - Set API subgroup to `catalog`, version `v1alpha1`
    - Apply resource-specific overrides (cross-resource references for `catalogTypeId`)
    - _Requirements: 1.2, 4.1_

  - [x] 4.5 Create incident domain resource configuration
    - Create `config/incident/config.go`
    - Register `incident_custom_field`, `incident_custom_field_option`, `incident_role`, `incident_severity`, `incident_status`
    - Set API subgroup to `incident`, version `v1alpha1`
    - Apply resource-specific overrides (cross-resource references for `customFieldId`)
    - _Requirements: 1.2, 5.1_

  - [x] 4.6 Create on-call domain resource configuration
    - Create `config/oncall/config.go`
    - Register `incident_escalation_path`, `incident_schedule`, `incident_maintenance_window`
    - Set API subgroup to `oncall`, version `v1alpha1`
    - _Requirements: 1.2, 6.1_

  - [x] 4.7 Create automation domain resource configuration
    - Create `config/automation/config.go`
    - Register `incident_workflow`
    - Set API subgroup to `automation`, version `v1alpha1`
    - _Requirements: 1.2, 7.1_

- [ ] 5. Implement code generator and provider entrypoints
  - [x] 5.1 Create generator entrypoint
    - Create `cmd/generator/main.go`
    - Load the Terraform provider binary schema
    - Pass provider configuration from `config/provider.go` to the Upjet pipeline
    - Generate Go types under `apis/`, CRD YAML under `package/crds/`, controllers under `internal/controller/`
    - _Requirements: 1.1, 1.4, 13.1, 13.3, 13.5_

  - [x] 5.2 Create provider entrypoint
    - Create `cmd/provider/main.go`
    - Parse CLI flags: `--debug`, `--poll`, `--max-reconcile-rate`
    - Initialize controller-runtime manager with leader election
    - Register all 16 resource controllers via generated `Setup` functions
    - Wire in the `TerraformSetupBuilder` from `internal/clients/incident.go`
    - Start the controller manager
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 13.4_

- [ ] 6. Run Upjet code generation and verify output
  - [x] 6.1 Execute `make generate` and validate generated artifacts
    - Run `make generate` to produce all generated code
    - Verify all 16 Go type files exist under `apis/` in correct domain subdirectories (`alerting/v1alpha1/`, `catalog/v1alpha1/`, `incident/v1alpha1/`, `oncall/v1alpha1/`, `automation/v1alpha1/`)
    - Verify all 16 CRD YAML files exist under `package/crds/`
    - Verify all 16 controller files exist under `internal/controller/`
    - Verify all CRDs use API group `incidentio.crossplane.io` and version `v1alpha1`
    - _Requirements: 1.1, 1.2, 1.3, 3.1, 3.2, 3.3, 4.1, 5.1, 6.1, 7.1, 13.1, 13.3, 13.5_

  - [x] 6.2 Verify generated code compiles
    - Run `make build` to compile the provider binary
    - Ensure no compilation errors in generated types or controllers
    - _Requirements: 1.1_

- [x] 7. Checkpoint - Code generation and build verified
  - Ensure `make generate` and `make build` both succeed cleanly, ask the user if questions arise.

- [ ] 8. Create package metadata and example resources
  - [x] 8.1 Create Crossplane package metadata
    - Create `package/crossplane.yaml` with provider package declaration
    - Set `apiVersion: meta.pkg.crossplane.io/v1`, `kind: Provider`
    - Set `metadata.name: provider-incident-io`
    - Set `spec.crossplane.version: ">=v1.14.0"`
    - Set `spec.controller.image` to `ghcr.io/avodah-inc/crossplane-provider-incident-io:VERSION`
    - _Requirements: 10.1, 10.2, 14.1_

  - [x] 8.2 Create example ProviderConfig and managed resource manifests
    - Create `examples/providerconfig.yaml` with ProviderConfig referencing a Secret
    - Create `examples/secret.yaml` with placeholder API key Secret
    - Create example manifests for at least one resource per domain group (Severity, CatalogType, AlertAttribute, EscalationPath, Workflow)
    - _Requirements: 2.1, 3.4, 4.2, 5.2, 6.2, 7.2_

- [ ] 9. Implement CRD validation pipeline
  - [x] 9.1 Add kubeconform CRD validation to build pipeline
    - Add `make validate-crds` target to Makefile
    - Run kubeconform in strict mode with default and CRDs-catalog schema locations
    - Validate all CRD YAML files in `package/crds/`
    - Fail the pipeline if any validation errors are reported
    - _Requirements: 11.1, 11.2, 11.3_

  - [x] 9.2 Write smoke tests for code generation output
    - Write a test script or Go test that verifies:
      - All 16 CRD files exist under `package/crds/`
      - All CRDs contain `group: incidentio.crossplane.io` and `version: v1alpha1`
      - All 16 Go type packages exist under `apis/`
      - All 16 controller packages exist under `internal/controller/`
    - _Requirements: 1.1, 1.2, 1.3, 13.1, 13.3_

- [ ] 10. Implement security scanning pipeline
  - [x] 10.1 Add Trivy filesystem scan to build pipeline
    - Add `make trivy-fs` target to Makefile
    - Run `trivy fs . --include-dev-deps` on the provider source code
    - Fail the pipeline if critical, high, or medium severity vulnerabilities are found
    - _Requirements: 12.1, 12.3_

  - [x] 10.2 Add Trivy container image scan to build pipeline
    - Add `make trivy-image` target to Makefile
    - Run `trivy image` on the built provider container image
    - Fail the pipeline if critical, high, or medium severity vulnerabilities are found
    - _Requirements: 12.2, 12.3_

  - [x] 10.3 Add SonarQube scan configuration
    - Create `sonar-project.properties` with project key, source paths, and exclusions
    - Add `make sonar` target to Makefile to run `sonar-scanner`
    - _Requirements: 12.4_

- [x] 11. Checkpoint - Validation and security scanning verified
  - Ensure `make validate-crds`, `make trivy-fs`, and `make sonar` all pass, ask the user if questions arise.

- [ ] 12. Build container image and Crossplane package
  - [x] 12.1 Build and validate multi-arch container image
    - Run `make docker-build` to produce the container image
    - Verify multi-arch manifest includes `linux/amd64` and `linux/arm64`
    - Verify the image uses distroless base (no shell, no package manager)
    - Run `make trivy-image` on the built image to verify no vulnerabilities
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 12.2_

  - [x] 12.2 Build and validate Crossplane package
    - Run `make xpkg-build` to produce the xpkg archive
    - Verify the xpkg contains `crossplane.yaml` and all 16 CRD files
    - _Requirements: 10.1, 10.2, 10.3_

  - [x] 12.3 Write integration tests for package installation
    - Write a test that installs the xpkg in a test cluster (kind or similar)
    - Verify all 16 CRDs are registered under `incidentio.crossplane.io` after installation
    - Verify the provider pod starts and becomes ready
    - _Requirements: 10.3, 14.2, 14.3_

- [ ] 13. Wire publishing targets and finalize
  - [x] 13.1 Configure container registry publishing
    - Ensure `make docker-push` pushes to `ghcr.io/avodah-inc/crossplane-provider-incident-io` with semantic version tags
    - Ensure `make xpkg-push` pushes the Crossplane package to the same registry
    - _Requirements: 10.4, 14.1_

  - [x] 13.2 Create CI/CD pipeline configuration
    - Create GitHub Actions workflow (or equivalent) that runs the full build pipeline:
      - `make generate` → `make build` → `make validate-crds` → `make trivy-fs` → `make docker-build` → `make trivy-image` → `make xpkg-build`
    - On tag push: additionally run `make docker-push` and `make xpkg-push`
    - _Requirements: 1.1, 9.4, 10.3, 10.4, 11.1, 12.1, 12.2, 12.3, 12.4_

- [x] 14. Final checkpoint - Full pipeline validation
  - Ensure the complete build pipeline (`make generate` → `make build` → `make validate-crds` → `make trivy-fs` → `make docker-build` → `make trivy-image` → `make xpkg-build`) runs end-to-end without errors, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirement acceptance criteria for traceability
- Checkpoints ensure incremental validation at key milestones
- Property-based testing is not applicable for this IaC project (see design document Testing Strategy section)
- All code is Go, following the standard Upjet provider patterns
- The build pipeline order is: generate → build → validate → scan → package → publish
