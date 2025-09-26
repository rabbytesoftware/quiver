package arrows

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	domain "github.com/rabbytesoftware/quiver/internal/models/arrow"
)

type ArrowsRepository struct {
	infrastructure *infrastructure.Infrastructure
}

func NewArrowsRepository(
	infrastructure *infrastructure.Infrastructure,
) ArrowsInterface {
	return &ArrowsRepository{
		infrastructure: infrastructure,
	}
}

func (a *ArrowsRepository) GetArrows() []domain.Arrow {
	return []domain.Arrow{}
}

func (a *ArrowsRepository) GetArrow(id string) *domain.Arrow {
	return nil
}

func (a *ArrowsRepository) CreateArrow(arrow *domain.Arrow) *domain.Arrow {
	return nil
}

func (a *ArrowsRepository) UpdateArrow(arrow *domain.Arrow) *domain.Arrow {
	return nil
}

func (a *ArrowsRepository) DeleteArrow(id string) error {
	return nil
}
