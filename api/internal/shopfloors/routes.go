package shopfloors

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/shopfloors", handler.Create)
	router.GET("/shopfloors", handler.FindAll)
	router.GET("/shopfloors/:id", handler.FindByID)
	router.GET("/shopfloors/tenant/:tenant_id", handler.FindByTenantID)
	router.PUT("/shopfloors/:id", handler.Update)
	router.DELETE("/shopfloors/:id", handler.Delete)
}