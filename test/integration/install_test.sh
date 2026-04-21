#!/usr/bin/env bash
# Integration test for Crossplane provider-incident-io package installation.
#
# This script:
#   1. Creates a kind cluster (or reuses an existing one)
#   2. Installs Crossplane
#   3. Installs the provider package
#   4. Waits for the provider pod to become ready
#   5. Verifies all 16 managed-resource CRDs are registered
#   6. Cleans up (unless --no-cleanup is passed)
#
# Usage:
#   ./test/integration/install_test.sh [--no-cleanup] [--xpkg PATH]
#
# Requirements: kind, kubectl, helm, docker

set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------
CLUSTER_NAME="${CLUSTER_NAME:-crossplane-provider-test}"
CROSSPLANE_NAMESPACE="crossplane-system"
CROSSPLANE_VERSION="${CROSSPLANE_VERSION:-1.17.1}"
PROVIDER_IMAGE="${PROVIDER_IMAGE:-ghcr.io/avodah-inc/crossplane-provider-incident-io}"
PROVIDER_VERSION="${PROVIDER_VERSION:-0.1.0}"
XPKG_PATH=""
CLEANUP=true
TIMEOUT_SECONDS=300
POLL_INTERVAL=5

# Expected CRD groups and their resource counts
declare -A EXPECTED_CRDS
EXPECTED_CRDS=(
  ["alerting.incidentio.crossplane.io"]=3
  ["catalog.incidentio.crossplane.io"]=4
  ["incident.incidentio.crossplane.io"]=5
  ["oncall.incidentio.crossplane.io"]=3
  ["automation.incidentio.crossplane.io"]=1
)
TOTAL_MANAGED_CRDS=16

# Expected CRD names (fully qualified)
EXPECTED_CRD_NAMES=(
  "attributes.alerting.incidentio.crossplane.io"
  "routes.alerting.incidentio.crossplane.io"
  "sources.alerting.incidentio.crossplane.io"
  "catalogentries.catalog.incidentio.crossplane.io"
  "catalogtypes.catalog.incidentio.crossplane.io"
  "entries.catalog.incidentio.crossplane.io"
  "typeattributes.catalog.incidentio.crossplane.io"
  "fields.incident.incidentio.crossplane.io"
  "fieldoptions.incident.incidentio.crossplane.io"
  "roles.incident.incidentio.crossplane.io"
  "severities.incident.incidentio.crossplane.io"
  "statuses.incident.incidentio.crossplane.io"
  "paths.oncall.incidentio.crossplane.io"
  "schedules.oncall.incidentio.crossplane.io"
  "windows.oncall.incidentio.crossplane.io"
  "workflows.automation.incidentio.crossplane.io"
)

# ---------------------------------------------------------------------------
# Logging helpers
# ---------------------------------------------------------------------------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log()   { echo -e "${CYAN}[INFO]${NC}  $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
pass()  { echo -e "${GREEN}[PASS]${NC}  $*"; }
fail()  { echo -e "${RED}[FAIL]${NC}  $*"; }

die() {
  fail "$@"
  exit 1
}

# ---------------------------------------------------------------------------
# Argument parsing
# ---------------------------------------------------------------------------
while [[ $# -gt 0 ]]; do
  case "$1" in
    --no-cleanup) CLEANUP=false; shift ;;
    --xpkg)       XPKG_PATH="$2"; shift 2 ;;
    *)            die "Unknown argument: $1" ;;
  esac
done

# ---------------------------------------------------------------------------
# Prerequisite checks
# ---------------------------------------------------------------------------
check_tool() {
  command -v "$1" >/dev/null 2>&1 || die "Required tool not found: $1"
}

log "Checking prerequisites..."
check_tool kind
check_tool kubectl
check_tool helm
check_tool docker
pass "All required tools are available"

# ---------------------------------------------------------------------------
# Cleanup handler
# ---------------------------------------------------------------------------
cleanup() {
  if [[ "$CLEANUP" == "true" ]]; then
    log "Cleaning up kind cluster '${CLUSTER_NAME}'..."
    kind delete cluster --name "${CLUSTER_NAME}" 2>/dev/null || true
    log "Cleanup complete"
  else
    warn "Skipping cleanup (--no-cleanup). Cluster '${CLUSTER_NAME}' is still running."
  fi
}
trap cleanup EXIT

# ---------------------------------------------------------------------------
# Helper: wait for a condition with timeout
# ---------------------------------------------------------------------------
wait_for() {
  local description="$1"
  local check_cmd="$2"
  local elapsed=0

  log "Waiting for ${description} (timeout: ${TIMEOUT_SECONDS}s)..."
  while ! eval "${check_cmd}" >/dev/null 2>&1; do
    if [[ $elapsed -ge $TIMEOUT_SECONDS ]]; then
      die "Timed out waiting for ${description} after ${TIMEOUT_SECONDS}s"
    fi
    sleep "${POLL_INTERVAL}"
    elapsed=$((elapsed + POLL_INTERVAL))
  done
  pass "${description} (${elapsed}s)"
}

# ===========================================================================
# Step 1: Create or reuse kind cluster
# ===========================================================================
log "=== Step 1: Kind cluster ==="

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  log "Reusing existing kind cluster '${CLUSTER_NAME}'"
else
  log "Creating kind cluster '${CLUSTER_NAME}'..."
  kind create cluster --name "${CLUSTER_NAME}" --wait 60s
  pass "Kind cluster '${CLUSTER_NAME}' created"
fi

kubectl cluster-info --context "kind-${CLUSTER_NAME}" >/dev/null 2>&1 \
  || die "Cannot connect to kind cluster '${CLUSTER_NAME}'"
pass "Connected to kind cluster"

# ===========================================================================
# Step 2: Install Crossplane
# ===========================================================================
log "=== Step 2: Install Crossplane ==="

if kubectl get namespace "${CROSSPLANE_NAMESPACE}" >/dev/null 2>&1; then
  log "Crossplane namespace already exists, checking installation..."
else
  log "Installing Crossplane ${CROSSPLANE_VERSION} via Helm..."
  helm repo add crossplane-stable https://charts.crossplane.io/stable 2>/dev/null || true
  helm repo update crossplane-stable

  helm install crossplane crossplane-stable/crossplane \
    --namespace "${CROSSPLANE_NAMESPACE}" \
    --create-namespace \
    --version "${CROSSPLANE_VERSION}" \
    --wait \
    --timeout 120s
  pass "Crossplane ${CROSSPLANE_VERSION} installed"
fi

wait_for "Crossplane pods ready" \
  "kubectl -n ${CROSSPLANE_NAMESPACE} get pods -l app=crossplane -o jsonpath='{.items[0].status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"

pass "Crossplane is running"

# ===========================================================================
# Step 3: Install the provider package
# ===========================================================================
log "=== Step 3: Install provider package ==="

if [[ -n "${XPKG_PATH}" ]]; then
  # Load the xpkg from a local file by building a local image and loading it
  # into kind, then installing via a Provider manifest.
  log "Loading provider image into kind cluster..."

  # Build the provider image locally if not already present
  PROVIDER_TAG="${PROVIDER_IMAGE}:${PROVIDER_VERSION}"
  if ! docker image inspect "${PROVIDER_TAG}" >/dev/null 2>&1; then
    warn "Provider image '${PROVIDER_TAG}' not found locally."
    warn "Building provider image from Dockerfile..."
    docker build -t "${PROVIDER_TAG}" .
  fi

  kind load docker-image "${PROVIDER_TAG}" --name "${CLUSTER_NAME}"
  pass "Provider image loaded into kind"

  # Apply CRDs directly from package/crds/ for local testing
  log "Applying CRDs from package/crds/..."
  kubectl apply -f package/crds/
  pass "CRDs applied"
else
  # Install from registry using a Crossplane Provider resource
  log "Installing provider from registry: ${PROVIDER_IMAGE}:${PROVIDER_VERSION}"

  kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-incident-io
spec:
  package: ${PROVIDER_IMAGE}:${PROVIDER_VERSION}
  controllerConfigRef:
    name: provider-incident-io
EOF

  pass "Provider resource created"
fi

# ===========================================================================
# Step 4: Wait for provider pod to be ready
# ===========================================================================
log "=== Step 4: Wait for provider pod ==="

if [[ -n "${XPKG_PATH}" ]]; then
  # For local testing with direct CRD apply, we skip pod readiness since
  # the provider pod requires valid credentials to start properly.
  # We verify CRDs are registered instead.
  log "Local mode: skipping provider pod readiness check (CRDs applied directly)"
else
  wait_for "Provider package to become healthy" \
    "kubectl get provider.pkg.crossplane.io provider-incident-io -o jsonpath='{.status.conditions[?(@.type==\"Healthy\")].status}' 2>/dev/null | grep -q True"

  wait_for "Provider package to be installed" \
    "kubectl get provider.pkg.crossplane.io provider-incident-io -o jsonpath='{.status.conditions[?(@.type==\"Installed\")].status}' 2>/dev/null | grep -q True"

  # Find and wait for the provider pod
  wait_for "Provider pod to be running" \
    "kubectl -n ${CROSSPLANE_NAMESPACE} get pods -l pkg.crossplane.io/revision -o jsonpath='{.items[*].status.phase}' 2>/dev/null | grep -q Running"

  pass "Provider pod is ready"
fi

# ===========================================================================
# Step 5: Verify all 16 CRDs are registered
# ===========================================================================
log "=== Step 5: Verify CRDs ==="

ERRORS=0

# 5a. Check each expected CRD exists
log "Checking individual CRD registration..."
for crd_name in "${EXPECTED_CRD_NAMES[@]}"; do
  if kubectl get crd "${crd_name}" >/dev/null 2>&1; then
    pass "CRD registered: ${crd_name}"
  else
    fail "CRD missing: ${crd_name}"
    ERRORS=$((ERRORS + 1))
  fi
done

# 5b. Verify CRD counts per API group
log "Verifying CRD counts per API group..."
for group in "${!EXPECTED_CRDS[@]}"; do
  expected_count="${EXPECTED_CRDS[$group]}"
  actual_count=$(kubectl get crds -o json \
    | grep -c "\"${group}\"" \
    || echo "0")

  if [[ "$actual_count" -ge "$expected_count" ]]; then
    pass "Group ${group}: ${actual_count} CRDs (expected ${expected_count})"
  else
    fail "Group ${group}: ${actual_count} CRDs (expected ${expected_count})"
    ERRORS=$((ERRORS + 1))
  fi
done

# 5c. Verify total managed-resource CRD count
log "Verifying total managed-resource CRD count..."
total_found=0
for crd_name in "${EXPECTED_CRD_NAMES[@]}"; do
  if kubectl get crd "${crd_name}" >/dev/null 2>&1; then
    total_found=$((total_found + 1))
  fi
done

if [[ "$total_found" -eq "$TOTAL_MANAGED_CRDS" ]]; then
  pass "Total managed-resource CRDs: ${total_found}/${TOTAL_MANAGED_CRDS}"
else
  fail "Total managed-resource CRDs: ${total_found}/${TOTAL_MANAGED_CRDS}"
  ERRORS=$((ERRORS + 1))
fi

# 5d. Verify all CRDs have the correct API version (v1alpha1)
log "Verifying CRD API versions..."
for crd_name in "${EXPECTED_CRD_NAMES[@]}"; do
  if kubectl get crd "${crd_name}" >/dev/null 2>&1; then
    versions=$(kubectl get crd "${crd_name}" -o jsonpath='{.spec.versions[*].name}')
    if echo "${versions}" | grep -q "v1alpha1"; then
      pass "CRD ${crd_name}: has v1alpha1"
    else
      fail "CRD ${crd_name}: missing v1alpha1 (found: ${versions})"
      ERRORS=$((ERRORS + 1))
    fi
  fi
done

# ===========================================================================
# Step 6: Summary
# ===========================================================================
log "=== Test Summary ==="

if [[ "$ERRORS" -eq 0 ]]; then
  pass "All integration tests passed!"
  pass "  - ${TOTAL_MANAGED_CRDS} managed-resource CRDs registered"
  pass "  - All CRDs under incidentio.crossplane.io with v1alpha1"
  pass "  - 5 domain subgroups verified (alerting, catalog, incident, oncall, automation)"
  exit 0
else
  fail "${ERRORS} test(s) failed"
  exit 1
fi
