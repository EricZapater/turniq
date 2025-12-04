package customers

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/customers", handler.Create)
	router.GET("/customers", handler.GetAll)
	router.GET("/customers/:id", handler.GetByID)
	router.PUT("/customers/:id", handler.Update)
	router.DELETE("/customers/:id", handler.Delete)
}