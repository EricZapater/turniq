package shifts

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/shifts", handler.Create)
	router.GET("/shifts", handler.FindAll)
	router.GET("/shifts/:id", handler.FindByID)
	router.GET("/shifts/shopfloor/:shopfloorID", handler.FindByShopfloorID)
	router.PUT("/shifts/:id", handler.Update)
	router.DELETE("/shifts/:id", handler.Delete)
}