package usecases

import (
	"github.com/rabbytesoftware/quiver/internal/repositories"
	"github.com/rabbytesoftware/quiver/internal/usecases/arrows"
	"github.com/rabbytesoftware/quiver/internal/usecases/quivers"
	"github.com/rabbytesoftware/quiver/internal/usecases/system"
)

type Usecases struct {
	Arrows *arrows.ArrowsUsecase
	Quivers *quivers.QuiversUsecase
	System *system.SystemUsecase
}

func NewUsecases(
	repositories *repositories.Repositories,
) *Usecases {
	return &Usecases{
		Arrows: arrows.NewArrowsUsecase(repositories),
		Quivers: quivers.NewQuiversUsecase(repositories),
		System: system.NewSystemUsecase(repositories),
	}
}