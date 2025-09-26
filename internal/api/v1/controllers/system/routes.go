package system

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/system/health"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/system"
)

func SetupRoutes(router *gin.RouterGroup, usecases *usecase.SystemUsecase) {
	healthHandler := health.NewHealthHandler(usecases)
	healthHandler.SetupRoutes(router)
}
