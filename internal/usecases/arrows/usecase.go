package arrows

import "github.com/rabbytesoftware/quiver/internal/repositories"

type ArrowsUsecase struct {
	repositories *repositories.Repositories
}

func NewArrowsUsecase(
	repositories *repositories.Repositories,
) *ArrowsUsecase {
	return &ArrowsUsecase{
		repositories: repositories,
	}
}
