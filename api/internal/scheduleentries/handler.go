package scheduleentries

import (
	"net/http"

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
	var request ScheduleEntryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule entry created successfully", "data": response})
}

func (h *Handler) FindByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	response, err := h.service.FindByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule entry found successfully", "data": response})
}

func (h *Handler) FindAll(c *gin.Context) {
	ctx := c.Request.Context()
	
	var filter ScheduleFilter

	if cid := c.Query("customer_id"); cid != "" {
		if id, err := uuid.Parse(cid); err == nil {
			filter.CustomerID = &id
		}
	}
	if sid := c.Query("shopfloor_id"); sid != "" {
		if id, err := uuid.Parse(sid); err == nil {
			filter.ShopfloorID = &id
		}
	}
	if wcid := c.Query("workcenter_id"); wcid != "" {
		if id, err := uuid.Parse(wcid); err == nil {
			filter.WorkcenterID = &id
		}
	}
	if oid := c.Query("operator_id"); oid != "" {
		if id, err := uuid.Parse(oid); err == nil {
			filter.OperatorID = &id
		}
	}
	if jid := c.Query("job_id"); jid != "" {
		if id, err := uuid.Parse(jid); err == nil {
			filter.JobID = &id
		}
	}

	// Dates
	if date := c.Query("date"); date != "" {
		filter.StartDate = &date
		filter.EndDate = &date
	}
	if from := c.Query("from"); from != "" {
		filter.StartDate = &from
	}
	if to := c.Query("to"); to != "" {
		filter.EndDate = &to
	}

	response, err := h.service.Search(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) FindFiltered(c *gin.Context) {
	ctx := c.Request.Context()
	
	shopfloorID := c.Query("shopfloor_id")
	date := c.Query("date")

	var response []ScheduleEntry
	var err error

	if shopfloorID != "" && date != "" {
		response, err = h.service.GetPlanning(ctx, shopfloorID, date)
	} else {
		response, err = h.service.FindAll(ctx)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	var request ScheduleEntryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.service.Update(ctx, id, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule entry updated successfully", "data": response})
}

func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.service.Delete(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule entry deleted successfully"})
}

func (h *Handler) Sync(c *gin.Context) {
	ctx := c.Request.Context()
	
	type SyncRequest struct {
		ShopfloorID string                 `json:"shopfloor_id" binding:"required"`
		Date        string                 `json:"date" binding:"required"`
		Entries     []ScheduleEntryRequest `json:"entries" binding:"required"`
	}

	var req SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Sync(ctx, req.ShopfloorID, req.Date, req.Entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Planning synced successfully"})
}
