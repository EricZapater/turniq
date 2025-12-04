package tenants

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/tenants", handler.Create)
	router.GET("/tenants", handler.GetAll)
	router.GET("/tenants/:id", handler.GetByID)
	router.PUT("/tenants/:id", handler.Update)
	router.DELETE("/tenants/:id", handler.Delete)

}