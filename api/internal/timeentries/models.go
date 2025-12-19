package timeentries

import (
	"time"

	"github.com/google/uuid"
)

type TimeEntry struct {
    ID           uuid.UUID  `json:"id"`    
    OperatorID   uuid.UUID  `json:"operator_id"`
    WorkcenterID *uuid.UUID `json:"workcenter_id,omitempty"`
    CheckIn      time.Time  `json:"check_in"`
    CheckOut     *time.Time `json:"check_out,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

type TimeEntryRequest struct {
    OperatorID   string     `json:"operator_id"`
    WorkcenterID *string    `json:"workcenter_id,omitempty"`
    CheckIn      time.Time  `json:"check_in"`
    CheckOut     *time.Time `json:"check_out,omitempty"`
}
