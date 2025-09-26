package usecases

import (
	"github.com/rabbytesoftware/quiver/internal/respositories"
	"github.com/rabbytesoftware/quiver/internal/usecases/packages"
	repositories "github.com/rabbytesoftware/quiver/internal/usecases/repositories"
	repository "github.com/rabbytesoftware/quiver/internal/usecases/repositories"
	"github.com/rabbytesoftware/quiver/internal/usecases/system"
)

type Usecases struct {
	Packages *packages.PackagesUsecase
	Repositories *repositories.RepositoriesUsecase
	System *system.SystemUsecase
}

func NewUsecases(
	repositories *respositories.Repositories,
	system *respositories.System,
	packages *respositories.Packages,
) *Usecases {
	return &Usecases{
		Packages: packages.NewPackagesUsecase(repositories),
		Repository: repository.NewRepositoriesUsecase(repositories),
		System: system.NewSystemUsecase(repositories),
	}
}
