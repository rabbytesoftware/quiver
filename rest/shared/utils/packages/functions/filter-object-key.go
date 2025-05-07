package packages_functions

import (
	packages "rounds.com.ar/watcher/packages"
	logger "rounds.com.ar/watcher/view/logger"
	pc "rounds.com.ar/watcher/sdk/base/package-config"
)

func filterRuntimeKey(pkg *packages.Package) pc.PackageConfig{
	return pc.PackageConfig{
		URL: pkg.URL,
		Name: pkg.Name,
		Version: pkg.Version,
		Maintainers: pkg.Maintainers,
		Icon: pkg.Icon,
		NetBridge: pkg.NetBridge,
		BuildNumber: pkg.BuildNumber,
	}
}

func FilterPackagesRuntimeKey(
  pkgs map[string]*packages.Package,
  pkgName string,
) (pc.PackageConfig, map[string]pc.PackageConfig) {

	if pkgName == "" {
		filteredPkgs := make(map[string]pc.PackageConfig)

		for name, pkg := range pkgs {
			filteredPkgs[name] = filterRuntimeKey(pkg)
		}

		return pc.PackageConfig{}, filteredPkgs
	}

	pkg, err := GetPackage(pkgName)
	if err != nil {
		logger.It.Error("Error: %v", err)
		return pc.PackageConfig{}, nil
	}

	filtered := filterRuntimeKey(pkg)
	return filtered, nil
}
