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

func (a *ArrowsRepository) Get() []domain.Arrow {
	return []domain.Arrow{}
}

func (a *ArrowsRepository) GetById(id string) *domain.Arrow {
	return nil
}

func (a *ArrowsRepository) Create(arrow *domain.Arrow) *domain.Arrow {
	return nil
}

func (a *ArrowsRepository) Update(arrow *domain.Arrow) *domain.Arrow {
	return nil
}

func (a *ArrowsRepository) DeleteById(id string) error {
	return nil
}
