package jobs

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/jobs", handler.Create)
	router.GET("/jobs", handler.FindAll)
	router.GET("/jobs/:id", handler.FindByID)
	router.GET("/jobs/workcenter/:workcenterID", handler.FindByWorkcenterID)
	router.GET("/jobs/shopfloor/:shopfloorID", handler.FindByShopFloorID)
	router.GET("/jobs/customer/:customerID", handler.FindByCustomerID)
	router.PUT("/jobs/:id", handler.Update)
	router.DELETE("/jobs/:id", handler.Delete)
}