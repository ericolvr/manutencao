package routes

import (
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SlaRoutes(
	router *gin.Engine,
	slaHandler *handlers.SlaHandler,
) {
	routes := router.Group("/api/v1/slas")
	{
		routes.POST("", slaHandler.Create)
		routes.GET("", slaHandler.List)
		routes.GET("/:id", slaHandler.FindByID)
		routes.GET("/client/:client/priority/:priority", slaHandler.FindByParams)
		routes.PUT("/:id", slaHandler.Update)
		routes.DELETE("/:id", slaHandler.Delete)
	}
}
