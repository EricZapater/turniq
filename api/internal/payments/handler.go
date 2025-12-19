package payments

import (
	"api/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx := c.Request.Context()
	var request PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment created successfully", "data": response})
}

func (h *Handler) FindAll(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx := c.Request.Context()

	// Parse Filter Params
	var filter PaymentFilter

	if cid := c.Query("customer_id"); cid != "" {
		if uid, err := uuid.Parse(cid); err == nil {
			filter.CustomerID = &uid
		}
	}

	if from := c.Query("from"); from != "" {
		// Try parsing Layout "2006-01-02"
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.StartDate = &t
		} else {
             // Fallback or error? ignore for now
             if t, err := time.Parse(time.RFC3339, from); err == nil {
                 filter.StartDate = &t
             }
        }
	}

	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
            // Set to end of day? 23:59:59?
            // Actually usually 'to' implies inclusive up to that date.
            // Let's add 24 hours to cover the day if it's just date.
            t = t.Add(24 * time.Hour).Add(-1 * time.Second) 
			filter.EndDate = &t
		} else {
            if t, err := time.Parse(time.RFC3339, to); err == nil {
                filter.EndDate = &t
            }
        }
	}

	response, err := h.service.Search(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payments found successfully", "data": response})
}

func (h *Handler) FindById(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx := c.Request.Context()
	id := c.Param("id")
	response, err := h.service.FindById(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment found successfully", "data": response})
}

func (h *Handler) FindByCustomerId(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx := c.Request.Context()
	customerId := c.Param("customer_id")
	response, err := h.service.FindByCustomerId(ctx, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payments found successfully", "data": response})
}

func (h *Handler) Update(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx := c.Request.Context()
	id := c.Param("id")
	var request PaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Update(ctx, id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment updated successfully", "data": response})
}

func (h *Handler) Delete(c *gin.Context) {
	if !middleware.IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
