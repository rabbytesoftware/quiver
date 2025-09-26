package quiver

import (
	"github.com/rabbytesoftware/quiver/internal/models/arrow"
	"github.com/rabbytesoftware/quiver/internal/models/system"
)

type Quiver struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Banner system.URL `json:"banner"`
	URL system.URL `json:"url"`
	Security system.Security `json:"security"`
	Maintainers []string `json:"maintainers"`
	Version string `json:"version"`
	InstalledArrows []arrow.Arrow `json:"installed_arrows"`
	ListedArrows []arrow.ArrowNamespace `json:"listed_arrows"`
}
