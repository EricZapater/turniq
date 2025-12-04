package payments

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/payments", handler.Create)
	router.GET("/payments", handler.FindAll)
	router.GET("/payments/:id", handler.FindById)
	router.GET("/payments/customer/:customer_id", handler.FindByCustomerId)
	router.PUT("/payments/:id", handler.Update)
	router.DELETE("/payments/:id", handler.Delete)
}