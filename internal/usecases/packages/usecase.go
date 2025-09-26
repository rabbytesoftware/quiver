package packages

import "github.com/rabbytesoftware/quiver/internal/infrastructure"

type PackagesUsecase struct {
	infrastructure *infrastructure.Infrastructure
}

func NewPackagesUsecase(
	infrastructure *infrastructure.Infrastructure,
) *PackagesUsecase {
	return &PackagesUsecase{
		infrastructure: infrastructure,
	}
}
