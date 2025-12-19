package shifts

import (
	"api/middleware"
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
	if customerID := middleware.GetCustomerID(c); customerID != "" {
		request.CustomerID = customerID
	}
	shift, err := h.service.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shift})
}

func (h *Handler) FindAll(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := middleware.GetCustomerID(c)
	if customerID == "" {
		customerID = c.Query("customer_id")
		if customerID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "customer_id required"})
			return
		}
	}
	
	shifts, err := h.service.FindByCustomerID(ctx, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shifts})
}

func (h *Handler) FindByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	shift, err := h.service.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shift})
}

func (h *Handler) FindByShopfloorID(c *gin.Context){
	ctx := c.Request.Context()
	shopfloorID := c.Param("shopfloorID")
	shifts, err := h.service.FindByShopfloorID(ctx, shopfloorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shifts})
}

func (h *Handler) Update(c *gin.Context){
	ctx := c.Request.Context()
	id := c.Param("id")
	var request ShiftRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if customerID := middleware.GetCustomerID(c); customerID != "" {
		request.CustomerID = customerID
	}
	shift, err := h.service.Update(ctx, id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": shift})
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

