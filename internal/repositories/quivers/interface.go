package quivers

import domain "github.com/rabbytesoftware/quiver/internal/models/quiver"

type QuiversInterface interface {
	GetQuivers() []domain.Quiver
	GetQuiver(id string) *domain.Quiver
	CreateQuiver(quiver *domain.Quiver) *domain.Quiver
	UpdateQuiver(quiver *domain.Quiver) *domain.Quiver
	DeleteQuiver(id string) error
}