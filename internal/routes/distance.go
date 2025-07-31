package routes

import (
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/gin-gonic/gin"
)

func DistanceRoutes(
	router *gin.Engine,
	distanceHandler *handlers.DistanceHandler,
) {
	routes := router.Group("/api/v1/distances")
	{
		routes.POST("", distanceHandler.Create)
		routes.GET("", distanceHandler.List)
		routes.GET("/:id", distanceHandler.FindByID)
		routes.GET("/number/:number", distanceHandler.FindByNumber)
		routes.PUT("/:id", distanceHandler.Update)
		routes.DELETE("/:id", distanceHandler.Delete)
	}
}
