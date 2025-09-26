package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/packages"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/repositories"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/system"
	"github.com/rabbytesoftware/quiver/internal/usecases"
)

func SetupRoutes(router *gin.Engine, usecases *usecases.Usecases) {
	v1 := router.Group("/api/v1")
	{
		packages.SetupRoutes(v1.Group("/package"), usecases.Packages)
		repositories.SetupRoutes(v1.Group("/repository"), usecases.Repository)
		system.SetupRoutes(v1.Group("/system"), usecases.System)
	}
}
