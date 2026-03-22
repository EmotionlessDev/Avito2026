package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/schedules"
)

type pgSchedule struct {
	ID        string    `db:"id"`
	RoomID    string    `db:"room_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	CreatedAt time.Time `db:"created_at"`
}

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

const createScheduleSQL = `
INSERT INTO schedules (room_id, start_time, end_time)
VALUES ($1, $2, $3)
RETURNING id, created_at
`

const insertDaySQL = `
INSERT INTO schedule_days (schedule_id, day_of_week)
VALUES ($1, $2)
`

func (s *Storage) CreateSchedule(
	ctx context.Context,
	tx *sql.Tx,
	roomID string,
	startTime, endTime time.Time,
	days []int,
) (*schedules.Schedule, error) {

	if tx == nil {
		return nil, common.ErrNilTx
	}

	if !startTime.Before(endTime) {
		return nil, common.ErrInvalidScheduleTime
	}

	var sched pgSchedule
	err := tx.QueryRowContext(ctx, createScheduleSQL, roomID, startTime, endTime).
		Scan(&sched.ID, &sched.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	sched.RoomID = roomID
	sched.StartTime = startTime
	sched.EndTime = endTime

	seenDays := make(map[int]bool)
	for _, d := range days {
		if d < 1 || d > 7 {
			return nil, common.ErrInvalidScheduleDay
		}
		if seenDays[d] {
			return nil, common.ErrDuplicateScheduleDay
		}
		seenDays[d] = true

		if _, err := tx.ExecContext(ctx, insertDaySQL, sched.ID, d); err != nil {
			return nil, fmt.Errorf("failed to insert day: %w", err)
		}
	}

	result := pgScheduleToDomain(&sched)
	result.DaysOfWeek = days

	return result, nil
}

const getScheduleByIDSQL = `
SELECT id, room_id, start_time, end_time, created_at
FROM schedules
WHERE id = $1
`

func (s *Storage) GetScheduleByID(ctx context.Context, tx *sql.Tx, scheduleID string) (*schedules.Schedule, error) {
	if tx == nil {
		return nil, common.ErrNilTx
	}

	var sched pgSchedule
	err := tx.QueryRowContext(ctx, getScheduleByIDSQL, scheduleID).
		Scan(&sched.ID, &sched.RoomID, &sched.StartTime, &sched.EndTime, &sched.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	result := pgScheduleToDomain(&sched)

	days, err := s.getScheduleDays(ctx, tx, sched.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule days: %w", err)
	}
	result.DaysOfWeek = days

	return result, nil
}

const getScheduleDaysSQL = `
SELECT day_of_week
FROM schedule_days
WHERE schedule_id = $1
`

func (s *Storage) getScheduleDays(ctx context.Context, tx *sql.Tx, scheduleID string) ([]int, error) {
	rows, err := tx.QueryContext(ctx, getScheduleDaysSQL, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query schedule days: %w", err)
	}
	defer rows.Close()

	var days []int
	for rows.Next() {
		var d int
		if err := rows.Scan(&d); err != nil {
			return nil, fmt.Errorf("failed to scan day: %w", err)
		}
		days = append(days, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating schedule days: %w", err)
	}

	return days, nil
}

const isScheduleExistsByRoomIDSQL = `
SELECT id
FROM schedules
WHERE room_id = $1
`

func (s *Storage) IsScheduleExistsByRoomID(ctx context.Context, tx *sql.Tx, roomID string) (bool, error) {
	var id string
	err := tx.QueryRowContext(ctx, isScheduleExistsByRoomIDSQL, roomID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check schedule: %w", err)
	}
	return true, nil
}

func pgScheduleToDomain(s *pgSchedule) *schedules.Schedule {
	return &schedules.Schedule{
		ID:        s.ID,
		RoomID:    s.RoomID,
		StartTime: s.StartTime.Format("15:04"),
		EndTime:   s.EndTime.Format("15:04"),
		CreatedAt: s.CreatedAt,
	}
}
