package shifts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h *Handler) Create(c *gin.Context){
	ctx := c.Request.Context()
	var request ShiftRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shift, err := h.service.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shift)
}

func (h *Handler) FindByShopfloorID(c *gin.Context){
	ctx := c.Request.Context()
	shopfloorID := c.Param("shopfloorID")
	shifts, err := h.service.FindByShopfloorID(ctx, shopfloorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shifts)
}

func (h *Handler) Update(c *gin.Context){
	ctx := c.Request.Context()
	id := c.Param("id")
	var request ShiftRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shift, err := h.service.Update(ctx, id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shift)
}

func (h *Handler) Delete(c *gin.Context){
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shift deleted successfully"})
}

