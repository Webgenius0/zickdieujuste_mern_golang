package schedule

import (
	"errors"
	"fmt"
	"time"

	"gotickets/internal/domain/schedule/dto"

	"github.com/google/uuid"
)

// Service defines the business logic for the schedule domain.
type Service interface {
	GetOrCreateSchedule(userID uuid.UUID) (*dto.ScheduleResponse, error)
	UpdateSchedule(userID uuid.UUID, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new schedule Service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetOrCreateSchedule returns the user's schedule, creating a sensible default on first access
// rather than returning a 404. Default: 05:00 morning, 21:00 night, UTC timezone, push enabled.
func (s *service) GetOrCreateSchedule(userID uuid.UUID) (*dto.ScheduleResponse, error) {
	existing, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}
	if existing != nil {
		return toResponse(existing), nil
	}

	// Create default schedule on first access
	defaultSchedule := &UserSchedule{
		UserID:            userID,
		MorningPrayerTime: "05:00:00",
		NightPrayerTime:   "21:00:00",
		Timezone:          "UTC",
		PushEnabled:       true,
	}
	if err := s.repo.Upsert(defaultSchedule); err != nil {
		return nil, fmt.Errorf("failed to create default schedule: %w", err)
	}
	return toResponse(defaultSchedule), nil
}

func (s *service) UpdateSchedule(userID uuid.UUID, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	// Validate IANA timezone before saving — reject invalid values with a clear error.
	if _, err := time.LoadLocation(req.Timezone); err != nil {
		return nil, errors.New("invalid timezone: must be a valid IANA timezone string (e.g. Asia/Dhaka, America/New_York)")
	}

	existing, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}

	var sched *UserSchedule
	if existing == nil {
		sched = &UserSchedule{UserID: userID}
	} else {
		sched = existing
	}

	sched.MorningPrayerTime = req.MorningPrayerTime
	sched.NightPrayerTime = req.NightPrayerTime
	sched.Timezone = req.Timezone
	if req.PushEnabled != nil {
		sched.PushEnabled = *req.PushEnabled
	}

	if err := s.repo.Upsert(sched); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}
	return toResponse(sched), nil
}

func toResponse(s *UserSchedule) *dto.ScheduleResponse {
	return &dto.ScheduleResponse{
		ID:                s.ID,
		UserID:            s.UserID,
		MorningPrayerTime: s.MorningPrayerTime,
		NightPrayerTime:   s.NightPrayerTime,
		Timezone:          s.Timezone,
		PushEnabled:       s.PushEnabled,
		UpdatedAt:         s.UpdatedAt,
	}
}
