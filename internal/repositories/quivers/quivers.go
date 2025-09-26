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

func (q *QuiversRepository) GetQuivers() []domain.Quiver {
	return []domain.Quiver{}
}

func (q *QuiversRepository) GetQuiver(id string) *domain.Quiver {
	return nil
}

func (q *QuiversRepository) CreateQuiver(quiver *domain.Quiver) *domain.Quiver {
	return nil
}

func (q *QuiversRepository) UpdateQuiver(quiver *domain.Quiver) *domain.Quiver {
	return nil
}

func (q *QuiversRepository) DeleteQuiver(id string) error {
	return nil
}
