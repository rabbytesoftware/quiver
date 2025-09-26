package repository

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
)

type RepositoryUsecase struct {
	infrastructure *infrastructure.Infrastructure
}

func NewRepositoryUsecase(
	infrastructure *infrastructure.Infrastructure,
) *RepositoryUsecase {
	return &RepositoryUsecase{
		infrastructure: infrastructure,
	}
}