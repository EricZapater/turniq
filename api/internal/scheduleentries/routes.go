package scheduleentries

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	router.POST("/schedule-entries/sync", handler.Sync)
	router.POST("/schedule-entries", handler.Create)
	router.GET("/schedule-entries", handler.FindAll)
	router.GET("/schedule-entries/:id", handler.FindByID)
	router.GET("/schedule-entries/filtered", handler.FindFiltered)
	router.PUT("/schedule-entries/:id", handler.Update)
	router.DELETE("/schedule-entries/:id", handler.Delete)
}
