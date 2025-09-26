package packages

import "github.com/rabbytesoftware/quiver/internal/infrastructure"

type Packages struct {
	infrastructure *infrastructure.Infrastructure
}

func NewPackages(
	infrastructure *infrastructure.Infrastructure,
) *Packages {
	return &Packages{
		infrastructure: infrastructure,
	}
}
