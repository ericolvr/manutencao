package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ericolvr/maintenance-v2/internal/handlers"
)

func SolutionRoutes(router *gin.Engine, handler *handlers.SolutionHandler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/solutions", handler.Create)
		v1.GET("/solutions", handler.List)
		v1.GET("/solutions/:id", handler.GetByID)
		v1.PUT("/solutions/:id", handler.Update)
		v1.DELETE("/solutions/:id", handler.Delete)
	}
}
