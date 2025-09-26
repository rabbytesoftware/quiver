package usecases

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	"github.com/rabbytesoftware/quiver/internal/usecases/packages"
	"github.com/rabbytesoftware/quiver/internal/usecases/repository"
	"github.com/rabbytesoftware/quiver/internal/usecases/system"
)

type Usecases struct {
	infrastructure *infrastructure.Infrastructure
	Packages *packages.PackagesUsecase
	Repository *repository.RepositoryUsecase
	System *system.SystemUsecase
}

func NewUsecases(
	infrastructure *infrastructure.Infrastructure,
) *Usecases {
	return &Usecases{
		infrastructure: infrastructure,
		Packages: packages.NewPackagesUsecase(infrastructure),
		Repository: repository.NewRepositoryUsecase(infrastructure),
		System: system.NewSystemUsecase(infrastructure),
	}
}
