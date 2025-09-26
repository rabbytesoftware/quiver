package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usecase "github.com/rabbytesoftware/quiver/internal/usecases/arrows"
)

type HealthHandler struct {
	usecases *usecase.ArrowsUsecase
}

func NewHealthHandler(
	usecases *usecase.ArrowsUsecase,
) *HealthHandler {
	return &HealthHandler{
		usecases: usecases,
	}
}

func (h *HealthHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.Handler())
}

func (h *HealthHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Arrow ready to fly",
		})	
	}
}
