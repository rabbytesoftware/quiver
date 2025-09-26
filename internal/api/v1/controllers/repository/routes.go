package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/repository/health"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/repository"
)

func SetupRoutes(
	router *gin.RouterGroup, 
	usecases *usecase.RepositoryUsecase,
) {
	healthHandler := health.NewHealthHandler(usecases)
	healthHandler.SetupRoutes(router)
}
