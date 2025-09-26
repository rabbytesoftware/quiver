package system

import "github.com/rabbytesoftware/quiver/internal/infrastructure"

type SystemUsecase struct {
	infrastructure *infrastructure.Infrastructure
}

func NewSystemUsecase(
	infrastructure *infrastructure.Infrastructure,
) *SystemUsecase {
	return &SystemUsecase{
		infrastructure: infrastructure,
	}
}
