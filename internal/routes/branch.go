package routes

import (
	"github.com/ericolvr/maintenance-v2/internal/handlers"
	"github.com/gin-gonic/gin"
)

func BranchRoutes(
	router *gin.Engine,
	branchHandler *handlers.BranchHandler,
) {
	routes := router.Group("/api/v1/branchs")
	{
		routes.POST("", branchHandler.Create)
		routes.GET("", branchHandler.List)
		routes.GET("/:id", branchHandler.FindByID)
		routes.GET("/client/:client", branchHandler.GetByClient)
		routes.PUT("/:id", branchHandler.Update)
		routes.DELETE("/:id", branchHandler.Delete)
	}
}
