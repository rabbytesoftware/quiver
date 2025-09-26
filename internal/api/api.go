package api

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/core/config"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"

	"github.com/gin-gonic/gin"
	v1 "github.com/rabbytesoftware/quiver/internal/api/v1"
	"github.com/rabbytesoftware/quiver/internal/usecases"
)

type API struct {
	router *gin.Engine
	watcher *watcher.Watcher
	usecases *usecases.Usecases
}

func NewAPI(
	watcher *watcher.Watcher,
	usecases *usecases.Usecases,
) *API {
	return &API{
		router: gin.Default(),
		watcher: watcher,
		usecases: usecases,
	}
}

func (a *API) Run() {
	a.SetupMiddleware()
	a.SetupRoutes()

	a.watcher.Info(fmt.Sprintf(
		"Initializing API on %s:%d",
		config.GetAPI().Host,
		config.GetAPI().Port,
	))

	a.router.Run(
		fmt.Sprintf("%s:%d", config.GetAPI().Host, config.GetAPI().Port),
	)
}

func (a *API) SetupMiddleware() {
	a.router.Use(gin.Recovery())
}

func (a *API) SetupRoutes() {
	v1.SetupRoutes(a.router, a.usecases)
}
