package usecases

import (
	"github.com/rabbytesoftware/quiver/internal/usecases/packages"
	"github.com/rabbytesoftware/quiver/internal/usecases/repository"
	"github.com/rabbytesoftware/quiver/internal/usecases/system"
)

type Usecases struct {
	Packages *packages.PackagesUsecase
	Repository *repository.RepositoryUsecase
	System *system.SystemUsecase
}

func NewUsecases() *Usecases {
	return &Usecases{
		Packages: packages.NewPackagesUsecase(),
		Repository: repository.NewRepositoryUsecase(),
		System: system.NewSystemUsecase(),
	}
}
