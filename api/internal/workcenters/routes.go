package workcenters

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/workcenters", handler.Create)
	router.GET("/workcenters", handler.FindAll)
	router.GET("/workcenters/:id", handler.FindByID)
	router.GET("/workcenters/tenant/:tenantID", handler.FindByTenantID)
	router.PUT("/workcenters/:id", handler.Update)
	router.DELETE("/workcenters/:id", handler.Delete)
}