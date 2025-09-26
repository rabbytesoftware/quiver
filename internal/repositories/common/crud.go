package common

type CRUD[t any] interface {
	Get() []t
	GetById(id string) *t
	Create(t *t) *t
	Update(t *t) *t
	DeleteById(id string) error
}
