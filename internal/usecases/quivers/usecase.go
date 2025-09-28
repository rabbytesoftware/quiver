package quivers

import "github.com/rabbytesoftware/quiver/internal/repositories"

type QuiversUsecase struct {
	repositories *repositories.Repositories
}

func NewQuiversUsecase(
	repositories *repositories.Repositories,
) *QuiversUsecase {
	return &QuiversUsecase{
		repositories: repositories,
	}
}
