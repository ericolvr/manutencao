package routes

import (
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/gin-gonic/gin"
)

func ProviderRoutes(
	router *gin.Engine,
	providerHandler *handlers.ProviderHandler,
) {
	routes := router.Group("/api/v1/providers")
	{
		routes.POST("", providerHandler.Create)
		routes.GET("", providerHandler.List)
		routes.GET("/:id", providerHandler.FindByID)
		routes.GET("/name/:name", providerHandler.FindByName)
		routes.PUT("/:id", providerHandler.Update)
		routes.DELETE("/:id", providerHandler.Delete)
	}
}
