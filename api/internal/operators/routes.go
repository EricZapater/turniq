package operators

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/operators", handler.Create)
	router.GET("/operators", handler.FindAll)
	router.GET("/operators/:id", handler.FindByID)
	router.PUT("/operators/:id", handler.Update)
	router.DELETE("/operators/:id", handler.Delete)
}	