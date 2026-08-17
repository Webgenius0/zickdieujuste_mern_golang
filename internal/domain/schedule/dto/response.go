package dto

import (
	"time"

	"github.com/google/uuid"
)

// ScheduleResponse is the shape returned for GET and PUT /api/v1/schedules/me.
type ScheduleResponse struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	MorningPrayerTime string    `json:"morning_prayer_time"`
	NightPrayerTime   string    `json:"night_prayer_time"`
	Timezone          string    `json:"timezone"`
	PushEnabled       bool      `json:"push_enabled"`
	UpdatedAt         time.Time `json:"updated_at"`
}
