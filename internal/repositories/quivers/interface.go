package quivers

import (
	domain "github.com/rabbytesoftware/quiver/internal/models/quiver"
	"github.com/rabbytesoftware/quiver/internal/repositories/common"
)

type QuiversInterface interface {
	common.CRUD[domain.Quiver]
}
