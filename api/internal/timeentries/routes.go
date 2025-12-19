package timeentries

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/time-entries", handler.Create)
	router.GET("/time-entries", handler.FindAll)
	router.GET("/time-entries/:id", handler.FindByID)
	router.GET("/time-entries/customer/:customer_id", handler.FindByCustomerID)
	router.GET("/time-entries/operator/:operator_id", handler.FindByOperatorID)
	router.GET("/time-entries/current/:operator_id", handler.FindCurrent)
	router.PUT("/time-entries/:id", handler.Update)
	router.DELETE("/time-entries/:id", handler.Delete)
}
