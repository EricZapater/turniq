package customers

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/customers", handler.Create)
	router.GET("/customers", handler.FindAll)
	router.GET("/customers/:id", handler.FindByID)
	router.PUT("/customers/:id", handler.Update)
	router.DELETE("/customers/:id", handler.Delete)
}