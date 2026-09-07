package dto

type CreateWorshipReq struct {
	Title        string `json:"title" validate:"required,max=255"`
	Artist       string `json:"artist" validate:"required,max=255"`
	TimeOfDay    string `json:"time_of_day" validate:"required,oneof=Day Night"`
	Duration     string `json:"duration" validate:"required"`
	ThumbnailURL string `json:"thumbnail_url" validate:"required,url"`
	AudioURL     string `json:"audio_url" validate:"required,url"`
	PrayerText   string `json:"prayer_text"`
}

type UpdateWorshipReq struct {
	Title        string `json:"title" validate:"required,max=255"`
	Artist       string `json:"artist" validate:"required,max=255"`
	TimeOfDay    string `json:"time_of_day" validate:"required,oneof=Day Night"`
	Duration     string `json:"duration" validate:"required"`
	ThumbnailURL string `json:"thumbnail_url" validate:"required,url"`
	AudioURL     string `json:"audio_url" validate:"required,url"`
	PrayerText   string `json:"prayer_text"`
}
