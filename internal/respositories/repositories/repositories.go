package repositories

import "github.com/rabbytesoftware/quiver/internal/infrastructure"

type Repositories struct {
	infrastructure *infrastructure.Infrastructure
}

func NewRepositories(
	infrastructure *infrastructure.Infrastructure,
) *Repositories {
	return &Repositories{
		infrastructure: infrastructure,
	}
}
