package packages

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/packages/health"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/packages"
)

func SetupRoutes(router *gin.RouterGroup, usecases *usecase.PackagesUsecase) {
	healthHandler := health.NewHealthHandler(usecases)
	healthHandler.SetupRoutes(router)
}
