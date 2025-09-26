package repository

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
)

type RepositoriesUsecase struct {
	infrastructure *infrastructure.Infrastructure
}

func NewRepositoriesUsecase(
	infrastructure *infrastructure.Infrastructure,
) *RepositoriesUsecase {
	return &RepositoriesUsecase{
		infrastructure: infrastructure,
	}
}