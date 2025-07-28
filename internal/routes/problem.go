package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ericolvr/maintenance-v2/internal/handlers"
)

func ProblemRoutes(router *gin.Engine, handler *handlers.ProblemHandler) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/problems", handler.Create)
		v1.GET("/problems", handler.List)
		v1.GET("/problems/:id", handler.GetByID)
		v1.PUT("/problems/:id", handler.Update)
		v1.DELETE("/problems/:id", handler.Delete)
	}
}
