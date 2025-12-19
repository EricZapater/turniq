package timeentries

import (
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
	ctx := c.Request.Context()
	var request TimeEntryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entry created successfully", "data": response})
}

func (h *Handler) FindByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	response, err := h.service.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entry found successfully", "data": response})
}

func (h *Handler) FindByCustomerID(c *gin.Context) {
	ctx := c.Request.Context()
	customerID := c.Param("customer_id")
	response, err := h.service.FindByCustomerID(ctx, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entries found successfully", "data": response})
}

func (h *Handler) FindAll(c *gin.Context) {
	ctx := c.Request.Context()
	
	var filter TimeEntryFilter

	if cid := c.Query("customer_id"); cid != "" {
		if id, err := uuid.Parse(cid); err == nil {
			filter.CustomerID = &id
		}
	}
	if oid := c.Query("operator_id"); oid != "" {
		if id, err := uuid.Parse(oid); err == nil {
			filter.OperatorID = &id
		}
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.StartDate = &t
		} else {
             if t, err := time.Parse(time.RFC3339, from); err == nil {
                 filter.StartDate = &t
             }
        }
	}

	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
             // Inclusive day
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
	c.JSON(http.StatusOK, gin.H{"message": "Time entries found successfully", "data": response})
}

func (h *Handler) FindByOperatorID(c *gin.Context) {
	ctx := c.Request.Context()
	operatorID := c.Param("operator_id")
	response, err := h.service.FindByOperatorID(ctx, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entries found successfully", "data": response})
}

func (h *Handler) FindCurrent(c *gin.Context) {
	ctx := c.Request.Context()
	operatorID := c.Param("operator_id")
	response, err := h.service.FindCurrent(ctx, operatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entry found successfully", "data": response})
}

func (h *Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var request TimeEntryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Update(ctx, id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entry updated successfully", "data": response})
}

func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Time entry deleted successfully"})
}
