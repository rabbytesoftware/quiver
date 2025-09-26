package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/packages"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/repository"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/system"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		packages.SetupRoutes(v1.Group("/packages"))
		
		repository.SetupRoutes(v1.Group("/repository"))

		system.SetupRoutes(v1.Group("/system"))
	}
}
