package scanner

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Package struct {
	Name      string
	Version   string
	Manager   string
	CachePath string
}

type PackageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ScanAllManagers() ([]Package, error) {
	var allPackages []Package
	var mu sync.Mutex
	var wg sync.WaitGroup

	managers := []struct {
		name    string
		scanner func() ([]Package, error)
	}{
		{"npm", ScanNpm},
		{"pnpm", ScanPnpm},
		{"yarn", ScanYarn},
		{"bun", ScanBun},
		{"deno", ScanDeno},
	}

	// Scan all managers concurrently
	for _, mgr := range managers {
		wg.Add(1)
		go func(name string, scanFunc func() ([]Package, error)) {
			defer wg.Done()
			packages, err := scanFunc()
			if err != nil {
				log.Printf("Error scanning %s: %v", name, err)
				return
			}
			mu.Lock()
			allPackages = append(allPackages, packages...)
			mu.Unlock()
		}(mgr.name, mgr.scanner)
	}

	wg.Wait()
	return allPackages, nil
}

func scanDirectory(basePath string, manager string) ([]Package, error) {
	packages := []Package{}

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return packages, nil
	}

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		if d.Name() == "package.json" {
			data, err := os.ReadFile(path)
			if err == nil {
				var pkgJSON PackageJSON
				if json.Unmarshal(data, &pkgJSON) == nil && pkgJSON.Name != "" && pkgJSON.Version != "" {
					packages = append(packages, Package{
						Name:      pkgJSON.Name,
						Version:   pkgJSON.Version,
						Manager:   manager,
						CachePath: path,
					})
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return packages, nil
}

func scanPaths(manager string, paths []string) ([]Package, error) {
	var allPackages []Package
	for _, path := range paths {
		if packages, err := scanDirectory(path, manager); err == nil {
			allPackages = append(allPackages, packages...)
		}
	}
	return allPackages, nil
}

func ScanNpm() ([]Package, error) {
	home, _ := os.UserHomeDir()
	if home == "" {
		return []Package{}, nil
	}
	return scanPaths("npm", []string{
		filepath.Join(home, ".npm"),
		"/opt/node22/lib/node_modules",
		"/usr/local/lib/node_modules",
		"/usr/lib/node_modules",
	})
}

func ScanPnpm() ([]Package, error) {
	home, _ := os.UserHomeDir()
	if home == "" {
		return []Package{}, nil
	}
	paths := []string{
		filepath.Join(home, ".pnpm-store"),
		filepath.Join(home, ".local", "share", "pnpm", "store"),
		filepath.Join(home, "Library", "pnpm", "store"),
	}
	if pnpmHome := os.Getenv("PNPM_HOME"); pnpmHome != "" {
		paths = append(paths, filepath.Join(pnpmHome, "store"))
	}
	return scanPaths("pnpm", paths)
}

func ScanYarn() ([]Package, error) {
	home, _ := os.UserHomeDir()
	if home == "" {
		return []Package{}, nil
	}
	paths := []string{
		filepath.Join(home, ".yarn-cache"),
		filepath.Join(home, ".cache", "yarn"),
		filepath.Join(home, "Library", "Caches", "Yarn"),
	}
	if yarnCache := os.Getenv("YARN_CACHE_FOLDER"); yarnCache != "" {
		paths = append(paths, yarnCache)
	}
	return scanPaths("yarn", paths)
}

func ScanBun() ([]Package, error) {
	home, _ := os.UserHomeDir()
	if home == "" {
		return []Package{}, nil
	}
	return scanPaths("bun", []string{
		filepath.Join(home, ".bun", "install", "cache"),
		filepath.Join(home, ".bun", "install", "global"),
		filepath.Join(home, ".cache", ".bun", "install", "cache"),
	})
}

func ScanDeno() ([]Package, error) {
	home, _ := os.UserHomeDir()
	if home == "" {
		return []Package{}, nil
	}
	denoDir := os.Getenv("DENO_DIR")
	if denoDir == "" {
		denoDir = filepath.Join(home, ".cache", "deno")
	}
	return scanPaths("deno", []string{
		filepath.Join(denoDir, "npm"),
		filepath.Join(denoDir, "deps", "https"),
	})
}

