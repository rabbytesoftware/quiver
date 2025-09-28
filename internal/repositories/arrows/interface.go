package arrows

import (
	domain "github.com/rabbytesoftware/quiver/internal/models/arrow"
	"github.com/rabbytesoftware/quiver/internal/repositories/common"
)

type ArrowsInterface interface {
	common.CRUD[domain.Arrow]
}
