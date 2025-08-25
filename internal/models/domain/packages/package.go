package packages

import (
	"github.com/google/uuid"
	"github.com/rabbytesoftware/quiver/internal/models/domain/port"
	"github.com/rabbytesoftware/quiver/internal/models/domain/requirement"
	"github.com/rabbytesoftware/quiver/internal/models/domain/runtime"
	"github.com/rabbytesoftware/quiver/internal/models/domain/system"
	"github.com/rabbytesoftware/quiver/internal/models/domain/variable"
)

type Package struct {
	ID uuid.UUID `json:"id"`
	Namespace PackageNamespace `json:"namespace"`
	ArrowVersion []string `json:"arrow_version"`
	Name string `json:"name"`
	Description string `json:"description"`
	Version string `json:"version"`
	License string `json:"license"`
	Maintainers []string `json:"maintainers"`
	Credits []string `json:"credits"`
	URL system.URL `json:"url"`
	Documentation string `json:"documentation"`

	Requirements requirement.Requirement `json:"requirements"`
	Dependencies []PackageNamespace `json:"dependencies"`

	Netbridge []port.PortRule `json:"netbridge"`
	Variables []variable.Variable `json:"variables"`

	Methods []runtime.Method `json:"methods"`
}
