package schedules

import (
	"context"
	"time"
)

type Schedule struct {
	ID         string    `json:"id"`
	RoomID     string    `json:"room_id"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	DaysOfWeek []int     `json:"days_of_week"`
	CreatedAt  time.Time `json:"created_at"`
}

type ScheduleStorage interface {
	CreateSchedule(
		ctx context.Context,
		roomID string,
		startTime, endTime time.Time,
		days []int,
	) (*Schedule, error)
	GetScheduleByID(ctx context.Context, scheduleID string) (*Schedule, error)
	IsScheduleExistsByRoomID(ctx context.Context, roomID string) (bool, error)
	GetScheduleByRoomID(ctx context.Context, roomID string) (*Schedule, error)
}
