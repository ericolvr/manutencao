package routes

import (
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/gin-gonic/gin"
)

func CostRoutes(
	router *gin.Engine,
	costHandler *handlers.CostHandler,
) {
	routes := router.Group("/api/v1/costs")
	{
		routes.GET("", costHandler.List)
		routes.GET("/:id", costHandler.FindByID)
		routes.PUT("/:id", costHandler.Update)
		routes.DELETE("/:id", costHandler.Delete)
	}
}
