package dto

// UpdateScheduleRequest is the body for PUT /api/v1/schedules/me.
type UpdateScheduleRequest struct {
	MorningPrayerTime string `json:"morning_prayer_time" validate:"required"`      // "HH:MM" or "HH:MM:SS"
	NightPrayerTime   string `json:"night_prayer_time"   validate:"required"`      // "HH:MM" or "HH:MM:SS"
	Timezone          string `json:"timezone"            validate:"required"`      // IANA tz, validated via time.LoadLocation
	PushEnabled       *bool  `json:"push_enabled"        validate:"required"`
}
