package scheduleentries

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleEntry struct {
	ID          uuid.UUID `json:"id"`	
	CustomerID  uuid.UUID `json:"customer_id"`
	ShopfloorID uuid.UUID `json:"shopfloor_id"`
	ShiftID     uuid.UUID `json:"shift_id"`
	WorkcenterID uuid.NullUUID `json:"workcenter_id"`
	JobID       uuid.NullUUID `json:"job_id"`
	OperatorID  uuid.NullUUID `json:"operator_id"`
	Date time.Time `json:"date"`
	Order int `json:"order"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	IsCompleted bool   `json:"is_completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ScheduleEntryRequest struct {
	ID          string `json:"id"`
	CustomerID  string `json:"customer_id"`
	ShopfloorID string `json:"shopfloor_id"`
	ShiftID     string `json:"shift_id"`
	WorkcenterID string `json:"workcenter_id"`
	JobID       string `json:"job_id"`
	OperatorID  string `json:"operator_id"`
	Date string `json:"date"`
	Order int `json:"order"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	IsCompleted bool   `json:"is_completed"`
}