package arrows

import (
	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/api/v1/controllers/arrows/health"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/arrows"
)

func SetupRoutes(router *gin.RouterGroup, usecases *usecase.ArrowsUsecase) {
	healthHandler := health.NewHealthHandler(usecases)
	healthHandler.SetupRoutes(router)
}
