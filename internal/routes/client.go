package routes

import (
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/gin-gonic/gin"
)

func ClientRoutes(
	router *gin.Engine,
	clientHandler *handlers.ClientHandler,

) {
	routes := router.Group("/api/v1/clients")
	{
		routes.POST("", clientHandler.Create)
		routes.GET("", clientHandler.List)
		routes.GET("/:id", clientHandler.FindByID)
		routes.PUT("/:id", clientHandler.Update)
		routes.DELETE("/:id", clientHandler.Delete)
	}
}
