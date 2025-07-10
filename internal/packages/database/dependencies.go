package database

import "fmt"

// Dependency management operations

// AddDependency adds a dependency relationship
func (m *Manager) AddDependency(packageName, dependencyName string) error {
	// Add to package's dependencies
	if pkg, exists := m.database.Installed[packageName]; exists {
		for _, dep := range pkg.Dependencies {
			if dep == dependencyName {
				return nil // Already exists
			}
		}
		pkg.Dependencies = append(pkg.Dependencies, dependencyName)
	}

	// Add to dependency's dependent_by
	if dep, exists := m.database.Installed[dependencyName]; exists {
		for _, dependent := range dep.DependentBy {
			if dependent == packageName {
				return m.Save() // Already exists
			}
		}
		dep.DependentBy = append(dep.DependentBy, packageName)
	}

	return m.Save()
}

// RemoveDependency removes a dependency relationship
func (m *Manager) RemoveDependency(packageName, dependencyName string) error {
	// Remove from package's dependencies
	if pkg, exists := m.database.Installed[packageName]; exists {
		for i, dep := range pkg.Dependencies {
			if dep == dependencyName {
				pkg.Dependencies = append(pkg.Dependencies[:i], pkg.Dependencies[i+1:]...)
				break
			}
		}
	}

	// Remove from dependency's dependent_by
	if dep, exists := m.database.Installed[dependencyName]; exists {
		for i, dependent := range dep.DependentBy {
			if dependent == packageName {
				dep.DependentBy = append(dep.DependentBy[:i], dep.DependentBy[i+1:]...)
				break
			}
		}
	}

	return m.Save()
}

// GetDependents returns packages that depend on the given package
func (m *Manager) GetDependents(packageName string) []string {
	if pkg, exists := m.database.Installed[packageName]; exists {
		return pkg.DependentBy
	}
	return []string{}
}

// GetDependencies returns the dependencies of a package
func (m *Manager) GetDependencies(packageName string) []string {
	if pkg, exists := m.database.Installed[packageName]; exists {
		return pkg.Dependencies
	}
	return []string{}
}

// HasDependents checks if a package has any dependents
func (m *Manager) HasDependents(packageName string) bool {
	dependents := m.GetDependents(packageName)
	return len(dependents) > 0
}

// GetDependencyTree returns the full dependency tree for a package
func (m *Manager) GetDependencyTree(packageName string) (map[string][]string, error) {
	if !m.IsInstalled(packageName) {
		return nil, fmt.Errorf("package %s is not installed", packageName)
	}

	tree := make(map[string][]string)
	visited := make(map[string]bool)
	
	m.buildDependencyTree(packageName, tree, visited)
	return tree, nil
}

// buildDependencyTree recursively builds the dependency tree
func (m *Manager) buildDependencyTree(packageName string, tree map[string][]string, visited map[string]bool) {
	if visited[packageName] {
		return // Avoid circular dependencies
	}
	
	visited[packageName] = true
	dependencies := m.GetDependencies(packageName)
	tree[packageName] = dependencies
	
	for _, dep := range dependencies {
		if m.IsInstalled(dep) {
			m.buildDependencyTree(dep, tree, visited)
		}
	}
}

// ValidateDependencies checks if all dependencies are satisfied
func (m *Manager) ValidateDependencies(packageName string) error {
	dependencies := m.GetDependencies(packageName)
	
	for _, dep := range dependencies {
		if !m.IsInstalled(dep) {
			return fmt.Errorf("dependency %s is not installed", dep)
		}
	}
	
	return nil
}

// GetUnusedDependencies returns dependencies that are not used by any package
func (m *Manager) GetUnusedDependencies() []string {
	var unused []string
	
	for packageName := range m.database.Installed {
		if len(m.GetDependents(packageName)) == 0 {
			// Check if this package is a root package (not a dependency of anything)
			isRootPackage := true
			for _, otherPkg := range m.database.Installed {
				for _, dep := range otherPkg.Dependencies {
					if dep == packageName {
						isRootPackage = false
						break
					}
				}
				if !isRootPackage {
					break
				}
			}
			
			if !isRootPackage {
				unused = append(unused, packageName)
			}
		}
	}
	
	return unused
} 