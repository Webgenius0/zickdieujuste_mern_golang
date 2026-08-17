package dto

import (
	"time"

	"github.com/google/uuid"
)

// ScheduleResponse is the shape returned for GET and PUT /api/v1/schedules/me.
type ScheduleResponse struct {
	ID                uuid.UUID `json:"id" example:"a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"`
	UserID            uuid.UUID `json:"user_id" example:"b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e"`
	MorningPrayerTime string    `json:"morning_prayer_time" example:"08:00"`
	NightPrayerTime   string    `json:"night_prayer_time" example:"22:00"`
	Timezone          string    `json:"timezone" example:"America/New_York"`
	PushEnabled       bool      `json:"push_enabled" example:"true"`
	UpdatedAt         time.Time `json:"updated_at" example:"2026-08-17T15:00:00Z"`
}
