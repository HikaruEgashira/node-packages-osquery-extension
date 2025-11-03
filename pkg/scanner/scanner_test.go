package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageJSONParsing(t *testing.T) {
	// Test that package.json files are correctly parsed in scanDirectory
	tmpDir := t.TempDir()

	// Create a test package.json
	testPkg := PackageJSON{
		Name:    "test-package",
		Version: "1.0.0",
	}

	pkgPath := filepath.Join(tmpDir, "package.json")
	data, err := json.Marshal(testPkg)
	require.NoError(t, err)

	err = os.WriteFile(pkgPath, data, 0644)
	require.NoError(t, err)

	// Scan the directory to parse the package.json
	results, err := scanDirectory(tmpDir, "npm")
	require.NoError(t, err)

	assert.Len(t, results, 1)
	assert.Equal(t, "test-package", results[0].Name)
	assert.Equal(t, "1.0.0", results[0].Version)
	assert.Equal(t, "npm", results[0].Manager)
	assert.Equal(t, pkgPath, results[0].CachePath)
}

func TestScanDirectory(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create nested packages
	packages := []struct {
		path    string
		name    string
		version string
	}{
		{"pkg1/package.json", "package-one", "1.0.0"},
		{"pkg2/package.json", "package-two", "2.0.0"},
		{"nested/pkg3/package.json", "package-three", "3.0.0"},
	}

	for _, pkg := range packages {
		pkgDir := filepath.Dir(filepath.Join(tmpDir, pkg.path))
		err := os.MkdirAll(pkgDir, 0755)
		require.NoError(t, err)

		testPkg := PackageJSON{
			Name:    pkg.name,
			Version: pkg.version,
		}
		data, err := json.Marshal(testPkg)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(tmpDir, pkg.path), data, 0644)
		require.NoError(t, err)
	}

	// Scan the directory
	results, err := scanDirectory(tmpDir, "test")
	require.NoError(t, err)

	// Should find all 3 packages
	assert.Equal(t, 3, len(results))

	// Check that all packages were found
	names := make(map[string]bool)
	for _, pkg := range results {
		names[pkg.Name] = true
		assert.Equal(t, "test", pkg.Manager)
	}

	assert.True(t, names["package-one"])
	assert.True(t, names["package-two"])
	assert.True(t, names["package-three"])
}

func TestScanNonExistentDirectory(t *testing.T) {
	results, err := scanDirectory("/nonexistent/path", "npm")
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestInvalidPackageJSON(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "package.json")

	// Write invalid JSON
	err := os.WriteFile(pkgPath, []byte("invalid json"), 0644)
	require.NoError(t, err)

	// scanDirectory should skip invalid JSON gracefully
	results, err := scanDirectory(tmpDir, "npm")
	require.NoError(t, err)

	// Should not find any packages (invalid JSON is skipped)
	assert.Len(t, results, 0)
}

func TestScanAllManagers(t *testing.T) {
	// This test will scan actual system directories if they exist
	// It should not fail even if no packages are found
	packages, err := ScanAllManagers()
	require.NoError(t, err)

	// Just verify it returns a slice (may be empty)
	assert.NotNil(t, packages)

	t.Logf("Found %d packages across all managers", len(packages))

	// Log first few packages if any found
	for i, pkg := range packages {
		if i >= 5 {
			break
		}
		t.Logf("  %s@%s (%s)", pkg.Name, pkg.Version, pkg.Manager)
	}
}


func TestPackageManagerScans(t *testing.T) {
	tests := []struct {
		name    string
		scanner func() ([]Package, error)
	}{
		{"npm", ScanNpm},
		{"pnpm", ScanPnpm},
		{"yarn", ScanYarn},
		{"bun", ScanBun},
		{"deno", ScanDeno},
		{"jsr", ScanJsr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packages, err := tt.scanner()
			require.NoError(t, err)

			if len(packages) > 0 {
				t.Logf("%s: Found %d packages", tt.name, len(packages))
				// Verify package structure
				for i, pkg := range packages {
					if i >= 3 {
						break
					}
					assert.NotEmpty(t, pkg.Name)
					assert.NotEmpty(t, pkg.Version)
					assert.Equal(t, tt.name, pkg.Manager)
					assert.NotEmpty(t, pkg.CachePath)
					t.Logf("  %s@%s", pkg.Name, pkg.Version)
				}
			} else {
				t.Logf("%s: No packages found (cache may not exist)", tt.name)
			}
		})
	}
}
