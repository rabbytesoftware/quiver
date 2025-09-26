package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/repositories/health"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/repositories"
)

func SetupRoutes(
	router *gin.RouterGroup, 
	usecases *usecase.RepositoriesUsecase,
) {
	healthHandler := health.NewHealthHandler(usecases)
	healthHandler.SetupRoutes(router)
}
