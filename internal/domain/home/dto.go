package home

type UserHome struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
	AgeGroup  string `json:"ageGroup"`
}

type ScheduledPrayer struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PrayerType   string `json:"prayerType"` // "Morning Prayer" or "Night Prayer"
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
}

type DailyQuote struct {
	Text      string `json:"text"`
	Reference string `json:"reference"`
}

type HomeResponse struct {
	User     UserHome          `json:"user"`
	Schedule []ScheduledPrayer `json:"schedule"` // Returns [] instead of null
	Quote    DailyQuote        `json:"quote"`
}
