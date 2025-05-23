package packages_functions

import (
	"fmt"
	"strings"

	packages "rounds.com.ar/watcher/packages"
	packages_global_variables "rounds.com.ar/watcher/rest/shared/utils/packages/global-variables"
)

func GetPackage(name string) (*packages.Package, error) {
	packagesList := packages_global_variables.Packages
	// Convert package name
	// from this 'package-name'
	// to this 'pkgs/package-name.watcher'
	pkgFullName := strings.Join([]string{"pkgs/", name, ".watcher"}, "")

	// Fetch package by fullName
	pkg, ok := packagesList[pkgFullName]

	if !ok {
    return nil, fmt.Errorf("package '%s' not found", name)
	}

  return pkg, nil
}