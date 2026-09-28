package home

import (
	"fmt"
	"time"

	"gotickets/internal/domain/prayer"
	"gotickets/internal/domain/quote"
	"gotickets/internal/domain/schedule"
	"gotickets/internal/domain/user"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	GetHomeData(userID uuid.UUID) (HomeResponse, error)
}

type service struct {
	userRepo     user.Repository
	scheduleRepo schedule.Repository
	prayerRepo   prayer.Repository
	quoteRepo    quote.Repository
}

func NewService(userRepo user.Repository, scheduleRepo schedule.Repository, prayerRepo prayer.Repository, quoteRepo quote.Repository) Service {
	return &service{
		userRepo:     userRepo,
		scheduleRepo: scheduleRepo,
		prayerRepo:   prayerRepo,
		quoteRepo:    quoteRepo,
	}
}

func (s *service) GetHomeData(userID uuid.UUID) (HomeResponse, error) {
	var resp HomeResponse
	var ageGroup string

	// 1. Fetch User Profile (Needs to happen first to get AgeGroup)
	usr, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return resp, fmt.Errorf("failed to get user: %w", err)
	}

	avatar := ""
	if usr.AvatarURL != nil {
		avatar = *usr.AvatarURL
	}
	
	// Determine age group
	if usr.Age < 13 {
		ageGroup = "Kids"
	} else if usr.Age <= 18 {
		ageGroup = "Teens"
	} else {
		ageGroup = "Adults"
	}
	
	resp.User = UserHome{
		Name:      usr.Name,
		AvatarURL: avatar,
		AgeGroup:  fmt.Sprintf("Age %s", formatAgeBadge(usr.Age)),
	}

	// Fetch concurrently
	g := new(errgroup.Group)

	var morningTime, nightTime string
	// 2. Fetch User Schedule Config
	g.Go(func() error {
		sch, err := s.scheduleRepo.GetByUserID(userID)
		if err != nil {
			return err
		}
		if sch == nil {
			morningTime = "05:30 AM"
			nightTime = "09:00 PM"
		} else {
			morningTime = formatTime(sch.MorningPrayerTime, "05:30 AM")
			nightTime = formatTime(sch.NightPrayerTime, "09:00 PM")
		}
		return nil
	})

	var morningPrayer, nightPrayer *prayer.Prayer
	// 3. Fetch Featured Prayers
	g.Go(func() error {
		mornings, _, err := s.prayerRepo.FindAllPrayers(1, 1, map[string]interface{}{
			"prayerType": string(prayer.PrayerTypeMorning),
			"ageGroup":   ageGroup,
			"isAdmin":    false,
		})
		if err != nil {
			return err
		}
		if len(mornings) > 0 {
			morningPrayer = &mornings[0]
		}

		nights, _, err := s.prayerRepo.FindAllPrayers(1, 1, map[string]interface{}{
			"prayerType": string(prayer.PrayerTypeNight),
			"ageGroup":   ageGroup,
			"isAdmin":    false,
		})
		if err != nil {
			return err
		}
		if len(nights) > 0 {
			nightPrayer = &nights[0]
		}
		return nil
	})

	// 4. Fetch Daily Quote
	g.Go(func() error {
		quotes, err := s.quoteRepo.FindPublished()
		if err != nil {
			return err
		}
		if len(quotes) > 0 {
			resp.Quote = DailyQuote{
				Text:      quotes[0].QuoteText,
				Reference: quotes[0].Reference,
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return resp, err
	}

	// 5. Combine Schedule Time with Prayers
	scheduledPrayers := make([]ScheduledPrayer, 0)
	
	if morningPrayer != nil {
		endTime := addDuration(morningTime, morningPrayer.Duration)
		scheduledPrayers = append(scheduledPrayers, ScheduledPrayer{
			ID:           morningPrayer.ID.String(),
			Title:        morningPrayer.Title,
			ThumbnailURL: morningPrayer.ThumbnailURL,
			PrayerType:   string(morningPrayer.PrayerType),
			StartTime:    morningTime,
			EndTime:      endTime,
		})
	}

	if nightPrayer != nil {
		endTime := addDuration(nightTime, nightPrayer.Duration)
		scheduledPrayers = append(scheduledPrayers, ScheduledPrayer{
			ID:           nightPrayer.ID.String(),
			Title:        nightPrayer.Title,
			ThumbnailURL: nightPrayer.ThumbnailURL,
			PrayerType:   string(nightPrayer.PrayerType),
			StartTime:    nightTime,
			EndTime:      endTime,
		})
	}

	resp.Schedule = scheduledPrayers

	return resp, nil
}

func formatAgeBadge(age int) string {
	if age < 13 {
		return "Under 13"
	}
	if age <= 18 {
		return "13-18"
	}
	return "19+"
}

func formatTime(timeStr, defaultVal string) string {
	// Simple formatting assuming HH:MM:SS or HH:MM
	t, err := time.Parse("15:04:05", timeStr)
	if err == nil {
		return t.Format("03:04 PM")
	}
	t, err = time.Parse("15:04", timeStr)
	if err == nil {
		return t.Format("03:04 PM")
	}
	return defaultVal
}

func addDuration(startTime string, durationStr string) string {
	if durationStr == "" {
		return startTime
	}
	t, err := time.Parse("03:04 PM", startTime)
	if err != nil {
		return startTime
	}
	
	// Assuming durationStr is something like "08:00" (MM:SS) or "8" (minutes)
	// For simplicity, we'll assume it's just minutes as a string for now, or parse it properly
	var minutes int
	fmt.Sscanf(durationStr, "%d", &minutes)
	if minutes == 0 {
		// try MM:SS format
		var m, s int
		fmt.Sscanf(durationStr, "%d:%d", &m, &s)
		minutes = m
	}
	
	if minutes == 0 {
		return startTime
	}
	
	t = t.Add(time.Duration(minutes) * time.Minute)
	return t.Format("03:04 PM")
}
