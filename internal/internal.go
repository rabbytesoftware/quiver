package internal

import (
	"github.com/rabbytesoftware/quiver/internal/api"
	"github.com/rabbytesoftware/quiver/internal/core"
	"github.com/rabbytesoftware/quiver/internal/modules"
	"github.com/rabbytesoftware/quiver/internal/usecases"
)

// ? Internal DI Container
// ? This is the initiator for internal services (Dependency Injection)

// ? This is spectial, because here we only expose the core services to the outside world.
// ? All other services are internal and are not exposed to the outside world as they are
// ? essential for the internal workings of the application and not intended to be used directly.

type Internal struct {
	core *core.Core
	api *api.API
	modules *modules.Modules
	usecases *usecases.Usecases
}

func NewInternal() *Internal {
	core := core.Init()

	return &Internal{
		core: core,
		api: api.NewAPI(core.GetWatcher()),
		modules: modules.NewModules(),
		usecases: usecases.NewUsecases(),
	}
}

func (i *Internal) Run() {
	i.api.Run()
}

func (i *Internal) GetCore() *core.Core {
	return i.core
}
