package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/users", handler.Create)
	router.GET("/users", handler.GetAll)
	router.GET("/users/:id", handler.GetByID)
	router.GET("/users/customer/:customer_id", handler.GetByCustomerID)
	router.PUT("/users/:id", handler.Update)
	router.DELETE("/users/:id", handler.Delete)
}