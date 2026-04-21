package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// expectedDomains lists the 5 domain groups that must exist under apis/ and internal/controller/.
var expectedDomains = []string{
	"alerting",
	"catalog",
	"incident",
	"oncall",
	"automation",
}

// expectedControllerResources maps each domain to its expected resource controller subdirectories.
// These 16 entries correspond to the 16 Terraform resources.
var expectedControllerResources = map[string][]string{
	"alerting":   {"attribute", "route", "source"},
	"catalog":    {"catalogentries", "catalogtype", "entry", "typeattribute"},
	"incident":   {"field", "fieldoption", "role", "severity", "status"},
	"oncall":     {"path", "schedule", "window"},
	"automation": {"workflow"},
}

const (
	crdDir        = "package/crds"
	apisDir       = "apis"
	controllerDir = "internal/controller"
	expectedCRDs  = 16
)

func projectRoot(t *testing.T) string {
	t.Helper()
	// The test runs from the test/ directory, so go up one level.
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("failed to resolve project root: %v", err)
	}
	return root
}

// TestCRDFilesExist verifies that exactly 16 managed-resource CRD files exist under package/crds/.
// ProviderConfig and ProviderConfigUsage CRDs are excluded from the count since they are
// infrastructure CRDs, not managed resource CRDs.
func TestCRDFilesExist(t *testing.T) {
	root := projectRoot(t)
	dir := filepath.Join(root, crdDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read CRD directory %s: %v", dir, err)
	}

	var managedCRDs []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".yaml") {
			continue
		}
		// Exclude ProviderConfig/ProviderConfigUsage — they are not managed resource CRDs.
		if strings.HasPrefix(name, "incidentio.crossplane.io_provider") {
			continue
		}
		managedCRDs = append(managedCRDs, name)
	}

	if len(managedCRDs) != expectedCRDs {
		t.Errorf("expected %d managed-resource CRD files, got %d", expectedCRDs, len(managedCRDs))
		for _, f := range managedCRDs {
			t.Logf("  found: %s", f)
		}
	}
}

// TestCRDGroupAndVersion reads every managed-resource CRD YAML and verifies:
//   - The group field contains ".incidentio.crossplane.io"
//   - The versions list includes "v1alpha1"
//
// Uses simple string matching on the YAML content to avoid external dependencies.
func TestCRDGroupAndVersion(t *testing.T) {
	root := projectRoot(t)
	dir := filepath.Join(root, crdDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read CRD directory %s: %v", dir, err)
	}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".yaml") {
			continue
		}
		if strings.HasPrefix(name, "incidentio.crossplane.io_provider") {
			continue
		}

		t.Run(name, func(t *testing.T) {
			path := filepath.Join(dir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read CRD file %s: %v", path, err)
			}

			content := string(data)

			// Verify group contains .incidentio.crossplane.io
			// CRD YAML has a line like: "  group: alerting.incidentio.crossplane.io"
			if !strings.Contains(content, ".incidentio.crossplane.io") {
				t.Errorf("CRD %s: does not contain group *.incidentio.crossplane.io", name)
			}

			// Verify the CRD has a line with "group:" that ends with ".incidentio.crossplane.io"
			foundGroup := false
			for _, line := range strings.Split(content, "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "group:") {
					groupValue := strings.TrimSpace(strings.TrimPrefix(trimmed, "group:"))
					if strings.HasSuffix(groupValue, ".incidentio.crossplane.io") {
						foundGroup = true
						break
					}
				}
			}
			if !foundGroup {
				t.Errorf("CRD %s: no 'group:' field ending with .incidentio.crossplane.io", name)
			}

			// Verify version v1alpha1 is present
			if !strings.Contains(content, "v1alpha1") {
				t.Errorf("CRD %s: does not contain version v1alpha1", name)
			}
		})
	}
}

// TestAPITypePackagesExist verifies that Go type packages exist under apis/ for each domain group.
// Each domain must have a v1alpha1/ subdirectory with at least one Go file.
func TestAPITypePackagesExist(t *testing.T) {
	root := projectRoot(t)

	for _, domain := range expectedDomains {
		t.Run(domain, func(t *testing.T) {
			domainDir := filepath.Join(root, apisDir, domain, "v1alpha1")
			info, err := os.Stat(domainDir)
			if err != nil {
				t.Fatalf("apis/%s/v1alpha1 directory does not exist: %v", domain, err)
			}
			if !info.IsDir() {
				t.Fatalf("apis/%s/v1alpha1 is not a directory", domain)
			}

			// Verify at least one Go file exists in the package.
			entries, err := os.ReadDir(domainDir)
			if err != nil {
				t.Fatalf("failed to read apis/%s/v1alpha1: %v", domain, err)
			}
			hasGoFile := false
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".go") {
					hasGoFile = true
					break
				}
			}
			if !hasGoFile {
				t.Errorf("apis/%s/v1alpha1 contains no Go files", domain)
			}
		})
	}
}

// TestControllerPackagesExist verifies that controller packages exist under internal/controller/
// for all 16 resources across the 5 domain groups.
func TestControllerPackagesExist(t *testing.T) {
	root := projectRoot(t)

	totalControllers := 0
	for domain, resources := range expectedControllerResources {
		for _, resource := range resources {
			totalControllers++
			name := domain + "/" + resource
			t.Run(name, func(t *testing.T) {
				ctrlDir := filepath.Join(root, controllerDir, domain, resource)
				info, err := os.Stat(ctrlDir)
				if err != nil {
					t.Fatalf("controller directory %s does not exist: %v", name, err)
				}
				if !info.IsDir() {
					t.Fatalf("controller path %s is not a directory", name)
				}

				// Verify at least one Go file (controller) exists.
				entries, err := os.ReadDir(ctrlDir)
				if err != nil {
					t.Fatalf("failed to read controller directory %s: %v", name, err)
				}
				hasGoFile := false
				for _, e := range entries {
					if strings.HasSuffix(e.Name(), ".go") {
						hasGoFile = true
						break
					}
				}
				if !hasGoFile {
					t.Errorf("controller directory %s contains no Go files", name)
				}
			})
		}
	}

	if totalControllers != expectedCRDs {
		t.Errorf("expected %d controller resource directories, mapped %d", expectedCRDs, totalControllers)
	}
}
