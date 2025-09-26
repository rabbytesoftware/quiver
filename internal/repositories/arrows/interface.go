package arrows

import domain "github.com/rabbytesoftware/quiver/internal/models/arrow"

type ArrowsInterface interface {
	GetArrows() []domain.Arrow
	GetArrow(id string) *domain.Arrow
	CreateArrow(arrow *domain.Arrow) *domain.Arrow
	UpdateArrow(arrow *domain.Arrow) *domain.Arrow
	DeleteArrow(id string) error
}
