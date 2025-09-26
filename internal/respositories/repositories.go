package respositories

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	"github.com/rabbytesoftware/quiver/internal/respositories/packages"
	"github.com/rabbytesoftware/quiver/internal/respositories/repositories"
	"github.com/rabbytesoftware/quiver/internal/respositories/system"
)

type Repositories struct {
	infrastructure *infrastructure.Infrastructure

	packages *packages.Packages
	system *system.System
	repositories *repositories.Repositories
}

func NewRepositories(
	infrastructure *infrastructure.Infrastructure,
) *Repositories {
	return &Repositories{
		infrastructure: infrastructure,
		packages: packages.NewPackages(infrastructure),
		system: system.NewSystem(infrastructure),
		repositories: repositories.NewRepositories(infrastructure),
	}
}