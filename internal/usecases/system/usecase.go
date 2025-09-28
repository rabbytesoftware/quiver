package system

import "github.com/rabbytesoftware/quiver/internal/repositories"

type SystemUsecase struct {
	repositories *repositories.Repositories
}

func NewSystemUsecase(
	repositories *repositories.Repositories,
) *SystemUsecase {
	return &SystemUsecase{
		repositories: repositories,
	}
}
