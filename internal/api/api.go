package api

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/core/config"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"

	"github.com/gin-gonic/gin"
	v1 "github.com/rabbytesoftware/quiver/internal/api/v1"
)

type API struct {
	router *gin.Engine
	watcher *watcher.Watcher
}

func NewAPI(
	watcher *watcher.Watcher,
) *API {
	return &API{
		router: gin.Default(),
		watcher: watcher,
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
	v1.SetupRoutes(a.router)
}
