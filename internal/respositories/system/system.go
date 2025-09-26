package system

import "github.com/rabbytesoftware/quiver/internal/infrastructure"

type System struct {
	infrastructure *infrastructure.Infrastructure
}

func NewSystem(
	infrastructure *infrastructure.Infrastructure,
) *System {
	return &System{
		infrastructure: infrastructure,
	}
}
