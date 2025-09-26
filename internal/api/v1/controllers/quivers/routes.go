package quivers

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/quivers/health"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/quivers"
)

func SetupRoutes(
	router *gin.RouterGroup, 
	usecases *usecase.QuiversUsecase,
) {
	healthHandler := health.NewHealthHandler(usecases)
	healthHandler.SetupRoutes(router)
}
