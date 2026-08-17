package dto

// UpdateScheduleRequest is the body for PUT /api/v1/schedules/me.
type UpdateScheduleRequest struct {
	MorningPrayerTime string `json:"morning_prayer_time" validate:"required" example:"08:00"`      // "HH:MM" or "HH:MM:SS"
	NightPrayerTime   string `json:"night_prayer_time"   validate:"required" example:"22:00"`      // "HH:MM" or "HH:MM:SS"
	Timezone          string `json:"timezone"            validate:"required" example:"America/New_York"` // IANA tz
	PushEnabled       *bool  `json:"push_enabled"        validate:"required" example:"true"`
}
