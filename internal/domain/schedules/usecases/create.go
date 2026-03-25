package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules/dto"
)

type CreateSchedule struct {
	scheduleStorage schedules.ScheduleStorage
}

func NewCreateSchedule(scheduleStorage schedules.ScheduleStorage) *CreateSchedule {
	return &CreateSchedule{
		scheduleStorage: scheduleStorage,
	}
}

func (uc *CreateSchedule) Execute(ctx context.Context, input dto.CreateScheduleInput) (*schedules.Schedule, error) {
	// Check days of week
	seen := make(map[int]bool)
	for _, d := range input.DaysOfWeek {
		if d < 1 || d > 7 {
			return nil, common.ErrInvalidScheduleDay
		}
		if seen[d] {
			return nil, common.ErrDuplicateScheduleDay
		}
		seen[d] = true
	}

	// Parse time and validate
	start, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}
	end, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
	}
	if !start.Before(end) {
		return nil, common.ErrInvalidScheduleTime
	}

	// Check if schedule already
	exists, err := uc.scheduleStorage.IsScheduleExistsByRoomID(ctx, input.RoomID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing schedule: %w", err)
	}
	if exists {
		return nil, common.ErrScheduleExists
	}

	// Create schedule
	sched, err := uc.scheduleStorage.CreateSchedule(ctx, input.RoomID, start, end, input.DaysOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return sched, nil
}
