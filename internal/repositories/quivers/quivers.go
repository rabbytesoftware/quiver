package quivers

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure"
	domain "github.com/rabbytesoftware/quiver/internal/models/quiver"
)

type QuiversRepository struct {
	infrastructure *infrastructure.Infrastructure
}

func NewQuiversRepository(
	infrastructure *infrastructure.Infrastructure,
) QuiversInterface {
	return &QuiversRepository{
		infrastructure: infrastructure,
	}
}

func (q *QuiversRepository) Get() []domain.Quiver {
	return []domain.Quiver{}
}

func (q *QuiversRepository) GetById(id string) *domain.Quiver {
	return nil
}

func (q *QuiversRepository) Create(quiver *domain.Quiver) *domain.Quiver {
	return nil
}

func (q *QuiversRepository) Update(quiver *domain.Quiver) *domain.Quiver {
	return nil
}

func (q *QuiversRepository) DeleteById(id string) error {
	return nil
}
