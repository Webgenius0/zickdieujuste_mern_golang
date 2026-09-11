package dto

import "time"

type CreateProverbReq struct {
	Title              string    `json:"title" validate:"required"`
	Category           string    `json:"category" validate:"required"`
	Duration           string    `json:"duration" validate:"required"`
	ThumbnailURL       string    `json:"thumbnail_url" validate:"required,url"`
	AudioURL           string    `json:"audio_url" validate:"required,url"`
	ScriptureReference string    `json:"scripture_reference" validate:"required"`
	MainText           string    `json:"main_text" validate:"required"`
	Explanation        string    `json:"explanation" validate:"required"`
	PublishDate        time.Time `json:"publish_date" validate:"required"`
}

type UpdateProverbReq struct {
	Title              string    `json:"title" validate:"required"`
	Category           string    `json:"category" validate:"required"`
	Duration           string    `json:"duration" validate:"required"`
	ThumbnailURL       string    `json:"thumbnail_url" validate:"required,url"`
	AudioURL           string    `json:"audio_url" validate:"required,url"`
	ScriptureReference string    `json:"scripture_reference" validate:"required"`
	MainText           string    `json:"main_text" validate:"required"`
	Explanation        string    `json:"explanation" validate:"required"`
	PublishDate        time.Time `json:"publish_date" validate:"required"`
}
