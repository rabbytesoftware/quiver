package repository

import (
	"github.com/rabbytesoftware/quiver/internal/models/packages"
	"github.com/rabbytesoftware/quiver/internal/models/system"
)

type Repository struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Banner system.URL `json:"banner"`
	URL system.URL `json:"url"`
	Security system.Security `json:"security"`
	Maintainers []string `json:"maintainers"`
	Version string `json:"version"`
	InstalledPackages []packages.Package `json:"installed_packages"`
	ListedPackages []packages.PackageNamespace `json:"listed_packages"`
}
